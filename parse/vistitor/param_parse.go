package vistitor

import (
	"caseGenerator/parse/bo"
	"go/ast"
	"log"
	"strings"
)

// ParamParse 处理参数：expr待处理的节点， name：节点名， typeParams 泛型关系
func ParamParse(expr ast.Expr, name string) *bo.ParamParseResult {
	typeParamsMap := bo.GetTypeParamMap()
	var db bo.ParamParseResult

	db.ParamName = name

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			bo.AppendImportList(bo.GetImportPathFromAliasMap(firstField))
		}
		if expr == "context.Context" {
			db.ParamInitValue = "context.Background()"
			bo.AppendImportList("\"context\"")
		} else if expr == "time.Time" {
			db.ParamInitValue = "time.Now()"
			bo.AppendImportList("\"time\"")
		} else {
			db.ParamInitValue = "utils.Empty[" + expr + "]()"
			bo.AppendImportList("\"caseGenerator/utils\"")
		}
	case *ast.Ident:
		result, ok := typeParamsMap[dbType.Name]
		if ok {
			db.ParamType = result.ParamType
		} else {
			db.ParamType = dbType.Name
		}
		db.ParamInitValue = "utils.Empty[" + db.ParamType + "]()"
		bo.AppendImportList("\"caseGenerator/utils\"")
		// 指针类型
	case *ast.StarExpr:
		// todo typeParamMap
		param := ParamParse(dbType.X, name)
		db.ParamInitValue = "lo.ToPtr(" + param.ParamInitValue + ")"
		db.ParamType = "*" + param.ParamType
		bo.AppendImportList("\"github.com/samber/lo\"")
	case *ast.FuncType:
		paramType := parseFuncType(dbType)

		db.ParamType = paramType
		db.ParamInitValue = "nil"
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
		db.ParamInitValue = "nil"
	case *ast.ArrayType:
		requestType := parseParamArrayType(dbType)
		db.ParamType = requestType
		db.ParamInitValue = "make(" + requestType + ", 0, 10)"
	case *ast.MapType:
		requestType := parseParamMapType(dbType)
		db.ParamType = requestType
		db.ParamInitValue = "make(" + db.ParamType + ", 10)"
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name)
		db.ParamType = "[]" + param.ParamType
		db.ParamInitValue = "[]" + param.ParamType + "{utils.Empty[" + param.ParamType + "]()}"
		bo.AppendImportList("\"caseGenerator/utils\"")

		db.IsEllipsis = true
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, name)
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
	return &db
}

func parseParamMapType(mpType *ast.MapType) string {
	typeParamsMap := bo.GetTypeParamMap()
	var keyInfo, valueInfo string
	// 处理key
	switch eltType := mpType.Key.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		keyInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			bo.AppendImportList(bo.GetImportPathFromAliasMap(firstField))
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
			bo.AppendImportList(bo.GetImportPathFromAliasMap(firstField))
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
	case *ast.InterfaceType:
		valueInfo = "any"
	case *ast.MapType:
		valueInfo = parseParamMapType(eltType)
	case *ast.ArrayType:
		valueInfo = parseParamArrayType(eltType)
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo
}

func parseParamArrayType(dbType *ast.ArrayType) string {
	paramTypeMap := bo.GetTypeParamMap()
	var requestType string
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			bo.AppendImportList(bo.GetImportPathFromAliasMap(firstField))
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
		return "[]*" + param.ParamType
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:

		requestType = parseParamArrayType(eltType)
		requestType = "[]" + requestType
	case *ast.MapType:
		requestType = parseParamMapType(eltType)
		requestType = "[]" + requestType
	case *ast.FuncType:
		paramType := parseFuncType(eltType)
		requestType = "[]" + paramType
	default:
		log.Fatalf("未知类型...")
	}
	return requestType
}

// ParseParamWithoutInit 不设置初始化的值，也就没有相应的依赖
func ParseParamWithoutInit(expr ast.Expr, name string) *bo.ParamParseResult {
	paramTypeMap := bo.GetTypeParamMap()
	// "_" 这种不处理了
	var db bo.ParamParseResult

	db.ParamName = name

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			bo.AppendImportList(bo.GetImportPathFromAliasMap(firstField))
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
	case *ast.FuncType:
		paramType := parseFuncType(dbType)
		db.ParamType = paramType
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
	case *ast.ArrayType:
		requestType := parseParamArrayType(dbType)
		db.ParamType = requestType
	case *ast.MapType:
		requestType := parseParamMapType(dbType)
		db.ParamType = requestType
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := ParamParse(dbType.Elt, name)
		db.ParamType = "[]" + param.ParamType
		db.IsEllipsis = true
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, name)
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
	return &db
}

func parseFuncType(dbType *ast.FuncType) string {
	var paramType = "func("
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := ParseParamWithoutInit(v.Type, "")
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
		return paramType
	}
	fields := dbType.Results.List
	if len(fields) > 0 {
		paramType = paramType + " ("
		for _, v := range fields {
			if len(v.Names) > 0 {
				for _, name := range v.Names {
					param := ParseParamWithoutInit(v.Type, name.Name)
					paramType = paramType + param.ParamType + ", "
				}
			} else if v.Type != nil {
				param := ParseParamWithoutInit(v.Type, "param")
				paramType = paramType + param.ParamType + ", "
			}

		}
		// 去掉最后一个逗号
		lastCommaIndex := strings.LastIndex(paramType, ",")
		if lastCommaIndex != -1 {
			paramType = paramType[:lastCommaIndex]
		}
		paramType = paramType + ")"
	}
	return paramType
}
