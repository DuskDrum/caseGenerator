package parse

import (
	"caseGenerator/generate"
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
	_ = Package{PackagePath: "", PackageName: f.Name.Name}
	// 2. 组装import信息
	InitImport()
	importParse(f)
	// 组装所有方法
	methods := make([]string, 0, 10)
	caseDetailList := make([]generate.CaseDetail, 0, 10)
	importPackageList := make([]string, 0, 10)
	// Loop in comment groups
OuterLoop:
	for _, cg := range f.Decls {
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
		methodInfo := Method{MethodName: decl.Name.Name}
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
		var requestVisitor vistitor.RequestVisitor
		ast.Walk(&requestVisitor, decl.Type.Params)

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
			RequestList: GetRequestDetailList(),
		}
		if GetReceiverInfo() == nil {
			cd.ReceiverType = "utils.Empty[" + GetReceiverInfo().ReceiverValue.InvocationName + "]()"
			importPackageList = append(importPackageList, "\"caseGenerator/utils\"")
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
		ImportPkgPaths: importPackageList,
		CaseDetailList: caseDetailList,
	}

	generate.GenFile(gm)

	return
}

func importParse(af *ast.File) {
	for _, importSpec := range af.Imports {
		if importSpec.Name == nil {
			AppendImportList(importSpec.Path.Value)
		} else {
			AppendAliasImport(importSpec.Name.Name, importSpec.Path.Value)
		}
	}
}

// containsMethod 判断是否存在于列表
func containsMethod(methodList []string, methodName string) bool {
	if len(methodList) == 0 {
		return false
	}
	for _, m := range methodList {
		if m == methodName {
			return true
		}
	}
	return false
}

type ParamParseResult struct {
	// 参数名
	ParamName string
	// 参数类型
	ParamType string
	// 参数初始化值
	ParamInitValue string
	// 参数校验值
	ParamCheckValue string
	// 参数引入的依赖包
	ImportPkgPath []string
	// 是否是省略号语法
	IsEllipsis bool
}

// ParamParse 处理参数：expr待处理的节点， name：节点名， typeParams 泛型关系
func ParamParse(expr ast.Expr, name string) *ParamParseResult {
	typeParamsMap := GetTypeParamMap()
	var db ParamParseResult

	db.ParamName = name
	importPaths := make([]string, 0, 10)

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, GetImportInfo().GetImportPath(firstField))
		}
		if expr == "context.Context" {
			db.ParamInitValue = "context.Background()"
			importPaths = append(importPaths, "\"context\"")
		} else if expr == "time.Time" {
			db.ParamInitValue = "time.Now()"
			importPaths = append(importPaths, "\"time\"")
		} else {
			db.ParamInitValue = "utils.Empty[" + expr + "]()"
			importPaths = append(importPaths, "\"caseGenerator/utils\"")
		}
	case *ast.Ident:
		result, ok := typeParamsMap[dbType.Name]
		if ok {
			db.ParamType = result.ParamType
		} else {
			db.ParamType = dbType.Name
		}
		db.ParamInitValue = "utils.Empty[" + db.ParamType + "]()"
		importPaths = append(importPaths, "\"caseGenerator/utils\"")
		// 指针类型
	case *ast.StarExpr:
		// todo typeParamMap
		param := ParamParse(dbType.X, name)
		db.ParamInitValue = "lo.ToPtr(" + param.ParamInitValue + ")"
		db.ParamType = "*" + param.ParamType
		importPaths = append(importPaths, "\"github.com/samber/lo\"")

		if len(param.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.FuncType:
		paramType, importList := parseFuncType(dbType)
		if len(importList) > 0 {
			for _, v := range importList {
				importPaths = append(importPaths, v)
			}
		}
		var resultVisit ResultVisitor
		ast.Walk(&resultVisit, dbType)

		db.ParamType = paramType
		db.ParamInitValue = "nil"
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
		db.ParamInitValue = "nil"
	case *ast.ArrayType:
		requestType, importList := parseParamArrayType(dbType)
		for _, v := range importList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + requestType + ", 0, 10)"
	case *ast.MapType:
		requestType, portList := parseParamMapType(dbType)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + db.ParamType + ", 10)"
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		db.ParamType = "[]" + param.ParamType
		db.ParamInitValue = "[]" + param.ParamType + "{utils.Empty[" + param.ParamType + "]()}"
		importPaths = append(importPaths, "\"caseGenerator/utils\"")
		db.IsEllipsis = true
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, name)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		if dbType.Dir == ast.RECV {
			db.ParamType = "<-chan " + param.ParamType
		} else {
			db.ParamType = "chan<- " + param.ParamType
		}
		db.ParamInitValue = "make(" + db.ParamType + ")"
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return nil
	default:
		log.Fatalf("未知类型...")
	}
	db.ImportPkgPath = importPaths
	return &db
}

func parseParamMapType(mpType *ast.MapType) (string, []string) {
	typeParamsMap := GetTypeParamMap()
	var keyInfo, valueInfo string
	importPaths := make([]string, 0, 10)
	// 处理key
	switch eltType := mpType.Key.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		keyInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, GetImportInfo().GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			keyInfo = result.ParamType
		} else {
			keyInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := ParseParamWithoutInit(eltType.X, "")
		keyInfo = "*" + param.ParamType
		if len(param.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.InterfaceType:
		keyInfo = "any"
	default:
		log.Fatalf("未知类型...")
	}

	// 处理value
	switch eltType := mpType.Value.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		valueInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, GetImportInfo().GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			valueInfo = result.ParamType
		} else {
			valueInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := ParseParamWithoutInit(eltType.X, "")
		valueInfo = "*" + param.ParamType
		if len(param.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.InterfaceType:
		valueInfo = "any"
	case *ast.MapType:
		var portList []string
		valueInfo, portList = parseParamMapType(eltType)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.ArrayType:
		var portList []string
		valueInfo, portList = parseParamArrayType(eltType)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo, importPaths
}

func parseParamArrayType(dbType *ast.ArrayType) (string, []string) {
	paramTypeMap := GetTypeParamMap()
	var requestType string
	importPaths := make([]string, 0, 10)
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, GetImportInfo().GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := paramTypeMap[eltType.Name]
		if ok {
			requestType = "[]" + result.ParamType
		} else {
			requestType = "[]" + eltType.Name
		}
	case *ast.StarExpr:
		param := ParseParamWithoutInit(eltType.X, "")
		return "[]*" + param.ParamType, param.ImportPkgPath
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:
		var portList []string

		requestType, portList = parseParamArrayType(eltType)
		requestType = "[]" + requestType
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.MapType:
		var portList []string
		requestType, portList = parseParamMapType(eltType)
		requestType = "[]" + requestType
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.FuncType:
		paramType, imports := parseFuncType(eltType)
		requestType = "[]" + paramType
		for _, v := range imports {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return requestType, importPaths
}

// ParseParamWithoutInit 不设置初始化的值，也就没有相应的依赖
func ParseParamWithoutInit(expr ast.Expr, name string) *ParamParseResult {
	paramTypeMap := GetTypeParamMap()
	// "_" 这种不处理了
	var db ParamParseResult

	db.ParamName = name
	importPaths := make([]string, 0, 10)

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, GetImportInfo().GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := paramTypeMap[dbType.Name]
		if ok {
			db.ParamType = result.ParamType
		} else {
			db.ParamType = dbType.Name
		}
		// 指针类型
	case *ast.StarExpr:
		param := ParseParamWithoutInit(dbType.X, name)
		db.ParamType = "*" + param.ParamType
		if len(db.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.FuncType:
		paramType, imports := parseFuncType(dbType)
		db.ParamType = paramType
		for _, v := range imports {
			db.ImportPkgPath = append(db.ImportPkgPath, v)
		}
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
	case *ast.ArrayType:
		requestType, importList := parseParamArrayType(dbType)
		for _, v := range importList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
	case *ast.MapType:
		requestType, portList := parseParamMapType(dbType)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		db.ParamType = "[]" + param.ParamType
		db.IsEllipsis = true
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, name)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		if dbType.Dir == ast.RECV {
			db.ParamType = "<-chan " + param.ParamType
		} else {
			db.ParamType = "chan<- " + param.ParamType
		}
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return nil
	case *ast.BinaryExpr:
		// 先直接取Y
		param := ParamParse(dbType.Y, name)
		db.ParamType = param.ParamType
	default:
		log.Fatalf("未知类型...")
	}
	db.ImportPkgPath = importPaths
	return &db
}

func parseFuncType(dbType *ast.FuncType) (string, []string) {
	importPaths := make([]string, 0, 10)
	var paramType = "func("
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := ParseParamWithoutInit(v.Type, "")
		if len(param.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		paramType = paramType + param.ParamType + ", "
	}
	// 去掉最后一个逗号
	lastCommaIndex := strings.LastIndex(paramType, ",")
	if lastCommaIndex != -1 {
		paramType = paramType[:lastCommaIndex]
	}
	paramType = paramType + ")"
	// 解析func的出参
	if dbType.Results == nil || dbType.Results.List == nil {
		return paramType, importPaths
	}
	fields := dbType.Results.List
	if len(fields) > 0 {
		paramType = paramType + " ("
		for _, v := range fields {
			if len(v.Names) > 0 {
				for _, name := range v.Names {
					param := ParseParamWithoutInit(v.Type, name.Name)
					if len(param.ImportPkgPath) > 0 {
						for _, v := range param.ImportPkgPath {
							importPaths = append(importPaths, v)
						}
					}
					paramType = paramType + param.ParamType + ", "
				}
			}

		}
		// 去掉最后一个逗号
		lastCommaIndex := strings.LastIndex(paramType, ",")
		if lastCommaIndex != -1 {
			paramType = paramType[:lastCommaIndex]
		}
		paramType = paramType + ")"
	}
	return paramType, importPaths
}
