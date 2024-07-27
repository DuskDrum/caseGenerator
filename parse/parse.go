package parse

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"caseGenerator/parse/vistitor"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Extract(path string, excludedPaths ...string) error {
	// 解析出path下每个目录下的私有方法，落入map
	//privateFunctionLinkedMap, err := extractPrivateFile(path)
	//if err != nil {
	//	return err
	//}

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
		for _, v := range list {
			// bo.AppendImportList("bytedanceMockey \"github.com/bytedance/mockey\"")
			//// 组装mock信息
			//mp := MockParam{
			//	// 判断是哪个类的方法，是指针还是值
			//	Caller:           v.RightType.Code,
			//	CallFunctionName: v.RightType.Desc,
			//	ReturnList:       v.LeftName,
			//}

			bo.AppendMockInfoList(generate.MockInstruct{
				MockResponseParam:  v.LeftName,
				MockFunction:       v.RightFormula,
				MockFunctionParam:  nil,
				MockFunctionResult: nil,
			})
		}
		// 9. 解析出所有需要go-link的function
		goLinkedList := make([]string, 0, 10)

		// 10. 开始处理receiver
		cd := generate.CaseDetail{
			CaseName:    methodInfo.MethodName,
			MethodName:  methodInfo.MethodName,
			RequestList: bo.GetRequestDetailList(),
			MockList:    bo.GetMockInfoList(),
			GoLinkList:  goLinkedList,
		}
		if bo.GetReceiverInfo() != nil {
			cd.ReceiverType = "utils.Empty[" + bo.GetReceiverInfo().ReceiverValue.InvocationName + "]()"
			//bo.AppendImportList("\"caseGenerator/utils\"")
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
		panic("路径存在问题")
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

// PrivateFunctionLinked 将所有private方法的信息存储在一起。一部分是内置函数，一部分是包中所有的私有方法的goLinked写法
type PrivateFunctionLinked struct {
	PrivateFunctionLinkedMap map[string]string
}

// GenerateKey path是全路径，判断如果是.go结尾，那么就拿最后一个目录文件夹
func (pfl PrivateFunctionLinked) GenerateKey(path string, funcName string) string {
	// 不以go结尾，直接拼接
	if filepath.Ext(path) != ".go" {
		return path + funcName
	} else {
		return path[0:strings.LastIndex(path, "/")]
	}
}

func (pfl PrivateFunctionLinked) GetLinkedInfoByFunctionName(key string) (string, error) {
	s, ok := pfl.PrivateFunctionLinkedMap[key]
	if ok {
		return s, nil
	} else {
		return "", errors.New("don't exist")
	}
}

func (pfl PrivateFunctionLinked) AddLinkedInfo(key string, linkedInfo string) {
	pfl.PrivateFunctionLinkedMap[key] = linkedInfo
}

type FunctionParam struct {
	requestList   []generate.CaseRequest
	responseList  []generate.CaseResponse
	aliasFuncName string
	funcName      string
	moduleName    string
	filepath      string
	// 需要依赖的import
	ImportPkgPaths []string
}

// GenerateMockGoLinked 根据解析出来的方法信息，生成go-linked记录
// //go:linkname awxCommonConvertSettlementReportAlgorithm slp/reconcile/core/message/standard.commonConvertSettlementReportAlgorithm
// func awxCommonConvertSettlementReportAlgorithm(transactionType enums.TransactionType, createdAt time.Time, ctx context.Context, dataBody dto.AwxSettlementReportDataBody) (result []service.OrderAlgorithmResult, err error)
func (fp FunctionParam) GenerateMockGoLinked() string {
	// 1. 使用 stringBuilder 解析go linkname
	var stringBuilder strings.Builder
	stringBuilder.WriteString("//go:linkname ")
	stringBuilder.WriteString(fp.aliasFuncName)
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(fp.moduleName)
	stringBuilder.WriteString("/")
	stringBuilder.WriteString(fp.filepath)
	stringBuilder.WriteString(".")
	stringBuilder.WriteString(fp.funcName)
	stringBuilder.WriteString("\n")
	// 2. 组装内部方法对应的结构
	stringBuilder.WriteString("func ")
	stringBuilder.WriteString(fp.aliasFuncName)
	stringBuilder.WriteString("(")
	// 2.1 解析内部方法的请求
	requestStr := ""
	for _, v := range fp.requestList {
		content := v.GenerateRequestContent()
		requestStr = requestStr + content + ", "
	}
	if len(requestStr) > 0 {
		requestStr = requestStr[0 : len(requestStr)-2]
	}
	stringBuilder.WriteString(requestStr)
	stringBuilder.WriteString(")")
	// 2.2 解析内部方法的响应
	if len(fp.responseList) > 0 {
		stringBuilder.WriteString("(")
		responseStr := ""
		for _, v := range fp.responseList {
			content := v.GenerateResponseContent()
			responseStr = responseStr + content + ", "
		}
		responseStr = responseStr[0 : len(responseStr)-2]
		stringBuilder.WriteString(responseStr)
		stringBuilder.WriteString(")")
	}

	return stringBuilder.String()
}

// clearingRecords := append(records, &record)
// mockey.Mock((*repo.ClearingPipeConfigRepo).GetAllConfigs).Return(clearingRecords).Build()
// 解析出mockey的方法，首先是调用方，然后是返回值
type MockParam struct {
	Caller           string
	CallFunctionName string
	ReturnList       []MockReturnParam
}

// MockReturnParam mock的返回值
// 首先是返回值的类型，然后是返回值的值(要mock的值)
type MockReturnParam struct {
	Name  string
	Value string
}

func (mp MockParam) GenerateMockString() string {
	var stringBuilder strings.Builder
	// 1. 判断是不是需要返回参数
	if len(mp.ReturnList) > 0 {
		for _, v := range mp.ReturnList {
			stringBuilder.WriteString(v.Name)
			stringBuilder.WriteString(" := ")
			stringBuilder.WriteString(v.Value)
			stringBuilder.WriteString("\n")
		}
	}
	// 2. 使用 stringBuilder 解析 mockey
	stringBuilder.WriteString("mockey.Mock((")
	stringBuilder.WriteString(mp.Caller)
	stringBuilder.WriteString(").")
	stringBuilder.WriteString(mp.CallFunctionName)
	stringBuilder.WriteString(").Return(")
	// 3. 解析返回值
	if len(mp.ReturnList) > 0 {
		for i, v := range mp.ReturnList {
			stringBuilder.WriteString(v.Name)
			if i != len(mp.ReturnList)-1 {
				stringBuilder.WriteString(",")
			}
		}
	}
	stringBuilder.WriteString(").Build()\n")
	return stringBuilder.String()
}
