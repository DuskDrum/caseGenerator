package visitor_v2

import (
	"caseGenerator/common/enum"
	"caseGenerator/parse/bo"
	"go/ast"
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
		if expr == "context.Context" {
			db.ParamInitValue = "context.Background()"
		} else if expr == "time.Time" {
			db.ParamInitValue = "time.Now()"
		} else {
			db.ParamInitValue = "utils.Empty[" + expr + "]()"
		}
	case *ast.Ident:
		result, ok := typeParamsMap[dbType.Name]
		if ok {
			db.ParamType = result.ParamType
		} else {
			db.ParamType = dbType.Name
		}
		db.ParamInitValue = "utils.Empty[" + db.ParamType + "]()"
		// 指针类型
	case *ast.StarExpr:
		// todo typeParamMap
		param := ParamParse(dbType.X, name)
		db.ParamInitValue = "lo.ToPtr(" + param.ParamInitValue + ")"
		db.ParamType = "*" + param.ParamType
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
		panic("未知类型...")
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
		panic("未知类型...")
	}

	// 处理value
	switch eltType := mpType.Value.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		valueInfo = expr
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
		panic("未知类型...")
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
		panic("未知类型...")
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
		panic("未知类型...")
	}
	return &db
}

// ParseParamRequest 解析请求
func ParseParamRequest(expr ast.Expr) []*bo.ParamParseRequest {
	paramTypeMap := bo.GetTypeParamMap()
	// "_" 这种不处理了
	var db bo.ParamParseRequest

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
	case *ast.Ident:
		result, ok := paramTypeMap[dbType.Name]
		if ok {
			db.ParamType = result.ParamType
		} else {
			db.ParamType = dbType.Name
		}
		// 指针类型
	case *ast.StarExpr:
		param := ParseParamWithoutInit(dbType.X, "")
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
		param := ParamParse(dbType.Elt, "")
		db.ParamType = "[]" + param.ParamType
	case *ast.ChanType:
		// 处理value
		param := ParamParse(dbType.Value, "")
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
		param := ParamParse(dbType.Y, "")
		db.ParamType = param.ParamType
	case *ast.BasicLit:
		db.ParamVale = dbType.Value
		db.ParamType = dbType.Kind.String()
	case *ast.FuncLit:
		funcType := parseFuncType(dbType.Type)
		db.ParamVale = funcType
		db.ParamType = enum.ASSIGNMENT_TYPE_FUNCTION.Code
	case *ast.CompositeLit:
		db.ParamType = enum.ASSIGNMENT_TYPE_COMPOSITE.Code
		init := ParseParamWithoutInit(dbType.Type, "")
		db.ParamVale = init.ParamType
	case *ast.CallExpr:
		args := dbType.Args
		requests := make([]*bo.ParamParseRequest, 0, 10)
		for _, v := range args {
			request := ParseParamRequest(v)
			for _, s := range request {
				requests = append(requests, s)
			}
		}
		return requests
	default:
		panic("未知类型...")
	}
	return []*bo.ParamParseRequest{&db}
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

// ParseFuncTypeParamParseResult 解析结果，第一个响应是请求参数列表，第二个参数是响应参数列表
func ParseFuncTypeParamParseResult(dbType *ast.FuncType) ([]bo.ParamParseResult, []bo.ParamParseResult) {
	// 请求参数列表
	requests := make([]bo.ParamParseResult, 0, 10)
	// 响应参数列表
	results := make([]bo.ParamParseResult, 0, 10)
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := ParseParamWithoutInit(v.Type, "")
		requests = append(requests, *param)
	}
	// 解析func的出参
	if dbType.Results == nil || dbType.Results.List == nil {
		return requests, results
	}
	fields := dbType.Results.List
	if len(fields) > 0 {
		for _, v := range fields {
			if len(v.Names) > 0 {
				for _, name := range v.Names {
					param := ParseParamWithoutInit(v.Type, name.Name)
					results = append(results, *param)
				}
			} else if v.Type != nil {
				param := ParseParamWithoutInit(v.Type, "param")
				results = append(results, *param)
			}
		}
	}
	return requests, results
}

func GetRelationFromSelectorExpr(se *ast.SelectorExpr) string {
	if si, ok := se.X.(*ast.Ident); ok {
		return si.Name + "." + se.Sel.Name
	}
	if sse, ok := se.X.(*ast.SelectorExpr); ok {
		return GetRelationFromSelectorExpr(sse) + "." + se.Sel.Name
	}
	return se.Sel.Name
}
