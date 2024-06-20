package parser

import (
	"caseGenerator/common/enum"
	"caseGenerator/parse/bo"
	"go/ast"
	"log"
	"strings"
)

// Param 参数信息，请求参数、返回参数、类型断言参数、变量、常量等
type Param struct {
	Name    string
	Type    string
	AstType enum.ParamAstType
}

// ParamParse 处理参数：expr待处理的节点， name：节点名， typeParams 泛型关系
func (s *SourceInfo) ParamParse(expr ast.Expr) *Param {
	genericsMap := s.GenericsMap
	var paramInfo Param

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		paramInfo.Type = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			s.AppendImportList(s.GetImportPathFromAliasMap(firstField))
		}
		paramInfo.AstType = enum.PARAM_AST_TYPE_SelectStmt
	case *ast.Ident:
		result, ok := genericsMap[dbType.Name]
		if ok {
			paramInfo.Type = result.Type
		} else {
			paramInfo.Type = dbType.Name
		}
		paramInfo.AstType = enum.PARAM_AST_TYPE_Ident
		// 指针类型
	case *ast.StarExpr:
		param := s.ParamParse(dbType.X)
		paramInfo.Type = "*" + param.Type
		paramInfo.AstType = enum.PARAM_AST_TYPE_StarExpr
	case *ast.FuncType:
		paramType := s.parseFuncType(dbType)
		paramInfo.Type = paramType
		paramInfo.AstType = enum.PARAM_AST_TYPE_FuncType
	case *ast.InterfaceType:
		// 啥也不做
		paramInfo.Type = "interface{}"
		paramInfo.AstType = enum.PARAM_AST_TYPE_InterfaceType
	case *ast.ArrayType:
		requestType := s.parseParamArrayType(dbType)
		paramInfo.Type = requestType
		paramInfo.AstType = enum.PARAM_AST_TYPE_ArrayType
	case *ast.MapType:
		requestType := s.parseParamMapType(dbType)
		paramInfo.Type = requestType
		paramInfo.AstType = enum.PARAM_AST_TYPE_MapType
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := s.ParamParse(dbType.Elt)
		paramInfo.Type = "[]" + param.Type
		paramInfo.AstType = enum.PARAM_AST_TYPE_Ellipsis
	case *ast.ChanType:
		// 处理value
		param := s.ParamParse(dbType.Value)
		if dbType.Dir == ast.RECV {
			paramInfo.Type = "<-chan " + param.Type
		} else {
			paramInfo.Type = "chan<- " + param.Type
		}
		paramInfo.AstType = enum.PARAM_AST_TYPE_ChanType
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		paramInfo.AstType = enum.PARAM_AST_TYPE_IndexExpr
		return nil
	default:
		log.Fatalf("未知类型...")
	}
	return &paramInfo
}

func (s *SourceInfo) parseParamMapType(mpType *ast.MapType) string {
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
			s.AppendImportList(s.GetImportPathFromAliasMap(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			keyInfo = result.ParamType
		} else {
			keyInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := s.ParamParse(eltType.X)
		keyInfo = "*" + param.Type
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
			s.AppendImportList(s.GetImportPathFromAliasMap(firstField))
		}
	case *ast.Ident:
		result, ok := typeParamsMap[eltType.Name]
		if ok {
			valueInfo = result.ParamType
		} else {
			valueInfo = eltType.Name
		}
	case *ast.StarExpr:
		param := s.ParamParse(eltType.X)
		valueInfo = "*" + param.Type
	case *ast.InterfaceType:
		valueInfo = "any"
	case *ast.MapType:
		valueInfo = s.parseParamMapType(eltType)
	case *ast.ArrayType:
		valueInfo = s.parseParamArrayType(eltType)
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo
}

func (s *SourceInfo) parseParamArrayType(dbType *ast.ArrayType) string {
	paramTypeMap := bo.GetTypeParamMap()
	var requestType string
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			s.AppendImportList(s.GetImportPathFromAliasMap(firstField))
		}
	case *ast.Ident:
		result, ok := paramTypeMap[eltType.Name]
		if ok {
			requestType = "[]" + result.ParamType
		} else {
			requestType = "[]" + eltType.Name
		}
	case *ast.StarExpr:
		param := s.ParamParse(eltType.X)
		return "[]*" + param.Type
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:
		requestType = s.parseParamArrayType(eltType)
		requestType = "[]" + requestType
	case *ast.MapType:
		requestType = s.parseParamMapType(eltType)
		requestType = "[]" + requestType
	case *ast.FuncType:
		paramType := s.parseFuncType(eltType)
		requestType = "[]" + paramType
	default:
		log.Fatalf("未知类型...")
	}
	return requestType
}

func (s *SourceInfo) parseFuncType(dbType *ast.FuncType) string {
	var paramType = "func("
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := s.ParamParse(v.Type)
		paramType = paramType + param.Type + ", "
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
				for _, _ = range v.Names {
					param := s.ParamParse(v.Type)
					paramType = paramType + param.Type + ", "
				}
			} else if v.Type != nil {
				param := s.ParamParse(v.Type)
				paramType = paramType + param.Type + ", "
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
func (s *SourceInfo) ParseFuncTypeParamParseResult(dbType *ast.FuncType) ([]*Param, []*Param) {
	// 请求参数列表
	requests := make([]*Param, 0, 10)
	// 响应参数列表
	results := make([]*Param, 0, 10)
	// 解析func的入参
	list := dbType.Params.List
	for _, v := range list {
		param := s.ParamParse(v.Type)
		requests = append(requests, param)
	}
	// 解析func的出参
	if dbType.Results == nil || dbType.Results.List == nil {
		return requests, results
	}
	fields := dbType.Results.List
	if len(fields) > 0 {
		for _, v := range fields {
			if len(v.Names) > 0 {
				for _, _ = range v.Names {
					param := s.ParamParse(v.Type)
					results = append(results, param)
				}
			} else if v.Type != nil {
				param := s.ParamParse(v.Type)
				results = append(results, param)
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
