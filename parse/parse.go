package parse

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"caseGenerator/parse/enum"
	"caseGenerator/parse/vistitor"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Extract(path string, excludedPaths ...string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Skip excluded paths
		for _, p := range excludedPaths {
			if p == path {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip vendor and all directories beginning with a .
		if info.IsDir() && (info.Name() == "vendor" || (len(info.Name()) > 1 && info.Name()[0] == '.')) {
			return filepath.SkipDir
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}
			if err = extractFile(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func extractFile(filename string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("extractFile parse error: %s", err)
		}
	}()
	// Parse file and create the AST
	var fset = token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments); err != nil {
		return
	}

	// 1. 组装package、method信息
	_ = bo.Package{PackagePath: "", PackageName: f.Name.Name}
	// 2. 组装import信息
	bo.InitImport(f)
	// 组装所有方法
	methods := make([]string, 0, 10)
	caseDetailList := make([]generate.CaseDetail, 0, 10)
	// Loop in comment groups
OuterLoop:
	for _, cg := range f.Decls {
		// 1. 先清理bo
		bo.ClearBo()
		decl, ok := cg.(*ast.FuncDecl)
		if !ok {
			fmt.Print("不处理非 ast.FuncDecl 的内容")
			continue
		}
		// 跳过某些方法的设置
		if decl.Doc != nil {
			for _, com := range decl.Doc.List {
				// 跳过处理
				if strings.Contains(com.Text, "skipcodeautoparse") {
					fmt.Print("skip auto parse...")
					continue OuterLoop
				}
			}
		}

		methods = append(methods, decl.Name.Name)
		// 3. 组装method方法
		methodInfo := bo.Method{MethodName: decl.Name.Name}
		fmt.Print("\n开始处理:" + decl.Name.Name)
		// 4. walk所有receiver.
		var receiverVisitor vistitor.ReceiverVisitor
		ast.Walk(&receiverVisitor, cg)
		// 5. 获取typeParam
		var typeParamVisitor vistitor.TypeParamVisitor
		ast.Walk(&typeParamVisitor, cg)
		// 6. 定义所有request
		vistitor.ParseRequest(decl.Type.Params)
		// 7. 判断是否有类型断言
		var typeAssertVisitor vistitor.TypeAssertionVisitor
		ast.Walk(&typeAssertVisitor, cg)
		combinationSLice := typeAssertVisitor.CombinationTypeAssertionRequest(bo.GetRequestDetailList())
		for i, v := range combinationSLice {
			cd := generate.CaseDetail{
				CaseName:    methodInfo.MethodName + strconv.Itoa(i),
				MethodName:  methodInfo.MethodName,
				RequestList: v,
			}
			caseDetailList = append(caseDetailList, cd)
		}

		// 8. 获取所有赋值、变量
		var assignmentVisitor vistitor.AssignmentVisitor
		ast.Walk(&assignmentVisitor, cg)
		// 8.1 遍历所有赋值和变量
		list := bo.GetAssignmentDetailInfoList()
		marshal, err := json.Marshal(list)
		if err != nil {
			fmt.Printf("遍历所有赋值失败，详情为:%s", err.Error())

		}
		fmt.Printf("遍历所有赋值为:%+v", string(marshal))
		for _, v := range list {
			if v.RightType == enum.RIGHT_TYPE_CALL {
				// mockey.Mock((*repo.ClearingPipeConfigRepo).GetAllConfigs).Return(clearingPipeConfigs).Build()
				// 	"github.com/bytedance/mockey"
				bo.AppendImportList("\"github.com/bytedance/mockey\"")
				bo.AppendMockInfoList(bo.MockInstruct{
					MockResponseParam:  v.LeftName,
					MockFunction:       v.RightFormula,
					MockFunctionParam:  nil,
					MockFunctionResult: nil,
				})
			}
		}

		// 10. 开始处理receiver
		cd := generate.CaseDetail{
			CaseName:    methodInfo.MethodName,
			MethodName:  methodInfo.MethodName,
			RequestList: bo.GetRequestDetailList(),
		}
		if bo.GetReceiverInfo() != nil {
			cd.ReceiverType = "utils.Empty[" + bo.GetReceiverInfo().ReceiverValue.InvocationName + "]()"
			bo.AppendImportList("\"caseGenerator/utils\"")
		}
		// 11. 开始处理mock

		// 12. 处理import
		caseDetailList = append(caseDetailList, cd)

	}

	packageName := f.Name.Name
	// ./core/service/complex_conditon_service.go，去掉前缀和后缀
	newFilename := strings.Replace(filename, "./", "", 1)
	newFilename = strings.TrimSuffix(newFilename, ".go")
	// 获取最后一个斜杠的位置
	lastIndex := strings.LastIndex(newFilename, "/")
	// 如果找到了斜杠
	if lastIndex == -1 {
		log.Fatalf("路径存在问题")
	}
	// 分割字符串
	filePath := newFilename[:lastIndex+1]
	fileName := newFilename[lastIndex+1:]

	gm := generate.GenMeta{
		Package:        packageName,
		FileName:       fileName,
		FilePath:       filePath,
		ImportPkgPaths: bo.GetImportInfo(),
		CaseDetailList: caseDetailList,
	}

	generate.GenFile(gm)

	return
}
