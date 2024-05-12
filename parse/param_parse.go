package parse

import (
	"caseGenerator/parse/bo"
	"caseGenerator/parse/vistitor"
	"go/ast"
	"log"
	"strings"
)

// ParamParse 处理参数：expr待处理的节点， name：节点名， typeParams 泛型关系
func ParamParse(expr ast.Expr, name string) *bo.ParamParseResult {
	typeParamsMap := bo.GetTypeParamMap()
	var db bo.ParamParseResult

	db.ParamName = name
	importPaths := make([]string, 0, 10)

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := vistitor.GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, bo.GetImportInfo().GetImportPath(firstField))
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
		var resultVisit vistitor.ResultVisitor
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
	typeParamsMap := bo.GetTypeParamMap()
	var keyInfo, valueInfo string
	importPaths := make([]string, 0, 10)
	// 处理key
	switch eltType := mpType.Key.(type) {
	case *ast.SelectorExpr:
		expr := vistitor.GetRelationFromSelectorExpr(eltType)
		keyInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, bo.GetImportInfo().GetImportPath(firstField))
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
		expr := vistitor.GetRelationFromSelectorExpr(eltType)
		valueInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, bo.GetImportInfo().GetImportPath(firstField))
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
	paramTypeMap := bo.GetTypeParamMap()
	var requestType string
	importPaths := make([]string, 0, 10)
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := vistitor.GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, bo.GetImportInfo().GetImportPath(firstField))
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
func ParseParamWithoutInit(expr ast.Expr, name string) *bo.ParamParseResult {
	paramTypeMap := bo.GetTypeParamMap()
	// "_" 这种不处理了
	var db bo.ParamParseResult

	db.ParamName = name
	importPaths := make([]string, 0, 10)

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := vistitor.GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, bo.GetImportInfo().GetImportPath(firstField))
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
