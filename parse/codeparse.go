package parse

import (
	"caseGenerator/generate"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

		// 4. 定义所有receiver.
		var paramVisitor ParamVisitor
		var rec Receiver
		if decl.Recv != nil {
			var recvName string
			switch recvType := decl.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				switch astStartExpr := recvType.X.(type) {
				case *ast.Ident:
					recvName = astStartExpr.Name
				// 下标类型，实际上是泛型。泛型先不处理
				case *ast.IndexExpr:
					continue
				}
			case *ast.Ident:
				recvName = recvType.Name
			default:
				recvName = decl.Recv.List[0].Names[0].Name
			}
			rec = Receiver{
				ReceiverName: recvName,
			}
			switch typeType := decl.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				switch astStartExpr := typeType.X.(type) {
				case *ast.Ident:
					rec.ReceiverValue = InvocationUnary{InvocationName: astStartExpr.Name}
				// 下标类型，实际上是泛型。泛型先不处理
				case *ast.IndexExpr:
					continue
				}
			case *ast.Ident:
				rec.ReceiverValue = InvocationUnary{InvocationName: typeType.Name}
			}
			paramVisitor.AddDetail(recvName, &rec)
			paramVisitorMarshal, err := json.Marshal(&paramVisitor)
			if err != nil {
				continue
			}
			fmt.Printf("接受者信息:%v \n", string(paramVisitorMarshal))
		}
		// 5. 获取typeParam
		typeParamMap := TypeParamParse(decl, ipInfo)
		SetTypeParamMap(typeParamMap)
		// 6. 判断是否有类型断言
		var typeAssertVisitor TypeAssertionVisitor

		ast.Walk(&typeAssertVisitor, cg)

		// 7. 定义所有request
		dbs := RequestInfoParse(decl, ipInfo)
		dbsMarshal, err := json.Marshal(&dbs)
		if err != nil {
			continue
		}
		fmt.Printf("请求参数信息:%v \n", string(dbsMarshal))
		// 8. 获取所有赋值、变量
		ast.Walk(&paramVisitor, cg)
		paramVMarshal, err := json.Marshal(&paramVisitor)
		if err != nil {
			continue
		}
		fmt.Printf("参数信息:%v \n", string(paramVMarshal))
		// 9. 开始处理数据，request信息，名称和类型
		rdList := make([]generate.RequestDetail, 0, 10)
		if len(dbs) > 0 {
			for _, db := range dbs {
				rdList = append(rdList, db)
			}
		}
		// 10. 开始处理receiver
		uu, _ := uuid.NewUUID()
		cd := generate.CaseDetail{
			CaseName:    uu.String(),
			MethodName:  methodInfo.MethodName,
			RequestList: rdList,
		}
		if decl.Recv != nil {
			cd.ReceiverType = "utils.Empty[" + rec.ReceiverValue.InvocationName + "]()"
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

// TypeParamParse 解析typeParam
func TypeParamParse(decl *ast.FuncDecl, ipInfo Import) map[string]*ParamParseResult {
	// 1. 处理typeParam
	typeParams := decl.Type.TypeParams
	typeParamsMap := make(map[string]*ParamParseResult, 10)

	if typeParams != nil && len(typeParams.List) > 0 {
		// 1. 一般只有一个
		field := typeParams.List[0]
		for _, v := range field.Names {
			ident, ok := field.Type.(*ast.Ident)
			if ok && ident.Name == "comparable" {
				typeParamsMap[v.Name] = lo.ToPtr(ParamParseResult{
					ParamName:      ident.Name,
					ParamType:      "string",
					ParamInitValue: "\"\"",
				})
				continue
			}
			//v.Obj.Decl
			init := parseParamWithoutInit(field.Type, v.Name, ipInfo, typeParamsMap)
			typeParamsMap[v.Name] = init
		}
	}

	return typeParamsMap
}

// RequestInfoParse RequestName:  "ctx",
// RequestType:  "context.Context",
// RequestValue: "context.Background()",
func RequestInfoParse(decl *ast.FuncDecl, ipInfo Import) []generate.RequestDetail {
	dbs := make([]generate.RequestDetail, 0, 10)
	for i, requestParam := range decl.Type.Params.List {
		// "_" 这种不处理了
		var db generate.RequestDetail
		if requestParam.Names == nil && requestParam.Type != nil {
			db.RequestName = "param" + strconv.Itoa(i)
			// todo typeParam
			result := ParamParse(requestParam.Type, db.RequestName, ipInfo, GetTypeParamMap())
			if result == nil {
				fmt.Println(result)
			}
			db.RequestType = result.ParamType
			db.RequestValue = result.ParamInitValue
			db.ImportPkgPath = result.ImportPkgPath
			db.IsEllipsis = result.IsEllipsis
			dbs = append(dbs, db)
			continue
		}

		names := requestParam.Names
		for j, name := range names {
			if name.Name == "_" {
				db.RequestName = "param" + strconv.Itoa(i) + strconv.Itoa(j)
			} else {
				db.RequestName = name.Name
			}
			// todo typeParamMap
			result := ParamParse(requestParam.Type, name.Name, ipInfo, GetTypeParamMap())
			if result == nil {
				fmt.Println(result)
			}
			db.RequestType = result.ParamType
			db.RequestValue = result.ParamInitValue
			db.ImportPkgPath = result.ImportPkgPath
			db.IsEllipsis = result.IsEllipsis
			dbs = append(dbs, db)
		}
	}
	return dbs
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

// ParamParse 处理参数：expr待处理的节点， name：节点名， ipInfo依赖关系， typeParams 泛型关系
func ParamParse(expr ast.Expr, name string, ipInfo Import, typeParamsMap map[string]*ParamParseResult) *ParamParseResult {
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
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
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
		param := ParamParse(dbType.X, name, ipInfo, GetTypeParamMap())
		db.ParamInitValue = "lo.ToPtr(" + param.ParamInitValue + ")"
		db.ParamType = "*" + param.ParamType
		importPaths = append(importPaths, "\"github.com/samber/lo\"")

		if len(param.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.FuncType:
		paramType, importList := parseFuncType(dbType, ipInfo, typeParamsMap)
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
		requestType, importList := parseParamArrayType(dbType, ipInfo, typeParamsMap)
		for _, v := range importList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + requestType + ", 0, 10)"
	case *ast.MapType:
		requestType, portList := parseParamMapType(dbType, ipInfo, typeParamsMap)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + db.ParamType + ", 10)"
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name, ipInfo, nil)
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
		param := ParamParse(dbType.Value, name, ipInfo, nil)
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

func parseParamMapType(mpType *ast.MapType, ipInfo Import, typeParamsMap map[string]*ParamParseResult) (string, []string) {
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
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			keyInfo = result.ParamType
		} else {
			keyInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := parseParamWithoutInit(eltType.X, "", ipInfo, nil)
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
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			valueInfo = result.ParamType
		} else {
			valueInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := parseParamWithoutInit(eltType.X, "", ipInfo, nil)
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
		valueInfo, portList = parseParamMapType(eltType, ipInfo, typeParamsMap)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.ArrayType:
		var portList []string
		valueInfo, portList = parseParamArrayType(eltType, ipInfo, typeParamsMap)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo, importPaths
}

func parseParamArrayType(dbType *ast.ArrayType, ipInfo Import, paramTypeMap map[string]*ParamParseResult) (string, []string) {
	var requestType string
	importPaths := make([]string, 0, 10)
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		result, ok := paramTypeMap[eltType.Name]
		if ok {
			requestType = "[]" + result.ParamType
		} else {
			requestType = "[]" + eltType.Name
		}
	case *ast.StarExpr:
		param := parseParamWithoutInit(eltType.X, "", ipInfo, nil)
		return "[]*" + param.ParamType, param.ImportPkgPath
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:
		var portList []string

		requestType, portList = parseParamArrayType(eltType, ipInfo, paramTypeMap)
		requestType = "[]" + requestType
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.MapType:
		var portList []string
		requestType, portList = parseParamMapType(eltType, ipInfo, paramTypeMap)
		requestType = "[]" + requestType
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.FuncType:
		paramType, imports := parseFuncType(eltType, ipInfo, paramTypeMap)
		requestType = "[]" + paramType
		for _, v := range imports {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return requestType, importPaths
}

// parseParamWithoutInit 不设置初始化的值，也就没有相应的依赖
func parseParamWithoutInit(expr ast.Expr, name string, paramTypeMap map[string]*ParamParseResult) *ParamParseResult {
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
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
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
		param := parseParamWithoutInit(dbType.X, name, ipInfo, paramTypeMap)
		db.ParamType = "*" + param.ParamType
		if len(db.ImportPkgPath) > 0 {
			for _, v := range param.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.FuncType:
		paramType, imports := parseFuncType(dbType, ipInfo, paramTypeMap)
		db.ParamType = paramType
		for _, v := range imports {
			db.ImportPkgPath = append(db.ImportPkgPath, v)
		}
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
	case *ast.ArrayType:
		requestType, importList := parseParamArrayType(dbType, ipInfo, paramTypeMap)
		for _, v := range importList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
	case *ast.MapType:
		requestType, portList := parseParamMapType(dbType, ipInfo, paramTypeMap)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name, ipInfo, nil)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		db.ParamType = "[]" + param.ParamType
		db.IsEllipsis = true
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, name, ipInfo, nil)
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
		param := ParamParse(dbType.Y, name, ipInfo, nil)
		db.ParamType = param.ParamType
	default:
		log.Fatalf("未知类型...")
	}
	db.ImportPkgPath = importPaths
	return &db
}

func parseFuncType(dbType *ast.FuncType, ipInfo Import, paramTypeMap map[string]*ParamParseResult) (string, []string) {
	importPaths := make([]string, 0, 10)
	var paramType = "func("
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := parseParamWithoutInit(v.Type, "", ipInfo, paramTypeMap)
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
					param := parseParamWithoutInit(v.Type, name.Name, ipInfo, paramTypeMap)
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
