package parse

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"caseGenerator/parse/vistitor"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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
	// Parse file and create the AST
	var fset = token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments); err != nil {
		return
	}

	// 1. 组装package、method信息
	_ = bo.Package{PackagePath: "", PackageName: f.Name.Name}
	// 2. 组装import信息
	bo.InitImport()
	importParse(f)
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
		// 6. 判断是否有类型断言
		var typeAssertVisitor vistitor.TypeAssertionVisitor
		ast.Walk(&typeAssertVisitor, cg)
		// 7. 定义所有request
		vistitor.ParseRequest(decl.Type.Params)

		// 8. 获取所有赋值、变量
		var paramVisitor vistitor.ParamVisitor
		ast.Walk(&paramVisitor, cg)
		paramVMarshal, err := json.Marshal(&paramVisitor)
		if err != nil {
			continue
		}
		fmt.Printf("参数信息:%v \n", string(paramVMarshal))

		// 10. 开始处理receiver
		uu, _ := uuid.NewUUID()
		cd := generate.CaseDetail{
			CaseName:    uu.String(),
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

func importParse(af *ast.File) {
	for _, importSpec := range af.Imports {
		if importSpec.Name == nil {
			bo.AppendImportList(importSpec.Path.Value)
		} else {
			bo.AppendAliasImport(importSpec.Name.Name, importSpec.Path.Value)
		}
	}
}
