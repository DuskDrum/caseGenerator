package parser

import (
	"caseGenerator/common/enum"
	"caseGenerator/parse/bo"
	"github.com/samber/lo"
	"go/ast"
	"log"
	"strings"
)

// Param 参数信息，请求参数、返回参数、类型断言参数、变量、常量等
type Param struct {
	Name    string            `json:"name,omitempty"`
	Type    string            `json:"type,omitempty"`
	AstType enum.ParamAstType `json:"astType,omitempty"`
}

// ParamValue 带有参数值的参数
type ParamValue struct {
	Param
	Value string
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
		// 下标类型，一般是泛型
		paramInfo.AstType = enum.PARAM_AST_TYPE_IndexExpr
		// 先解析主体
		param := s.ParamParse(dbType.X)
		// 再解析下标结构
		indexParam := s.ParamParseValue(dbType.Index)
		var index string
		if indexParam.Value == "" {
			index = indexParam.Type
		} else {
			index = indexParam.Value
		}
		paramInfo.Type = param.Type + "[" + index + "]"
	case *ast.CallExpr:
		// 一般无需处理CallExpr
	default:
		panic("未知类型...")
	}
	return &paramInfo
}

// ParamParseValue 同时处理几种可能存在值的类型，如BasicLit、FuncLit、CompositeLit、CallExpr
func (s *SourceInfo) ParamParseValue(expr ast.Expr) *ParamValue {
	genericsMap := s.GenericsMap
	var paramInfo ParamValue

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
	case *ast.BasicLit:
		paramInfo.Type = dbType.Kind.String()
		paramInfo.AstType = enum.PARAM_AST_TYPE_BasicLit
		paramInfo.Value = dbType.Value
		// FuncLit 等待解析出内容值
	case *ast.FuncLit:
		funcType := s.parseFuncType(dbType.Type)
		paramInfo.Type = funcType
		paramInfo.AstType = enum.PARAM_AST_TYPE_FuncLit
		// CompositeLit 等待解析出内容值
	case *ast.CompositeLit:
		var clv CompositeLitValue
		if dbType.Type != nil {
			init := s.ParamParse(dbType.Type)
			paramInfo.Type = init.Type
			clv.Param = init
		}
		paramInfo.AstType = enum.PARAM_AST_TYPE_CompositeLit
		// 如果Incomplete为true，那么这个类型是不完整的
		if dbType.Incomplete {
			log.Fatal("compositeLit is Incomplete")
		}
		// 解析composite的内容
		values := make([]ParamValue, 0, 10)
		for _, v := range dbType.Elts {
			paramValue := s.ParamParseValue(v)
			values = append(values, *paramValue)
		}
		clv.Values = values
		paramInfo.Value = clv.ToString()
		// CallExpr 等待解析出内容值
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		var paramType string
		param := s.ParamParse(dbType.Fun)
		paramType = param.Type + "("
		for _, v := range dbType.Args {
			reqInfos := s.ParseParamRequest(v)
			for _, reqInfo := range reqInfos {
				if reqInfo.Value != "" {
					paramType += reqInfo.Value + ", "
				} else if reqInfo.Type != "" {
					paramType += reqInfo.Type + ", "
				}
			}
		}
		if len(dbType.Args) > 0 {
			paramType = paramType[:len(paramType)-2]
		}
		paramType += ")"

		paramInfo.Type = paramType
		paramInfo.AstType = enum.PARAM_AST_TYPE_CallExpr
	case *ast.KeyValueExpr:
		key := s.ParamParse(dbType.Key)
		value := s.ParamParseValue(dbType.Value)
		paramInfo.Type = key.Type
		paramInfo.AstType = enum.PARAM_AST_TYPE_KeyValueExp
		if value.Value == "" {
			paramInfo.Value = value.Type
		} else {
			paramInfo.Value = value.Value
		}
		// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		result := s.ParamParseValue(dbType.X)
		paramInfo.Value = AssignmentUnaryNode{Param: result.Param, Op: dbType.Op}.ToString(result.Value)
		paramInfo.Type = result.Type
	default:
		panic("未知类型...")
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
		panic("未知类型...")
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
		panic("未知类型...")
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
	case *ast.StructType:
		requestType = "[]" + s.parseStructType(eltType)
	default:
		panic("未知类型...")
	}
	return requestType
}

func (s *SourceInfo) parseStructType(dbType *ast.StructType) string {
	var requestType = "struct{"
	// 解析struct的字段
	fields := dbType.Fields.List
	for _, v := range fields {
		typeParam := s.ParamParse(v.Type)
		for _, name := range v.Names {
			nameParam := s.ParamParse(name)
			requestType = requestType + " \n " + nameParam.Type + " " + typeParam.Type
		}
	}
	// 去掉最后一个逗号
	requestType = requestType + "}"
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

// ParseReceiver 解析receiver
func (s *SourceInfo) ParseReceiver(funcDecl *ast.FuncDecl) *Param {
	if funcDecl.Recv == nil {
		return nil
	}
	if len(funcDecl.Recv.List) == 0 {
		return nil
	}
	var recvName string
	switch recvType := funcDecl.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		switch astStartExpr := recvType.X.(type) {
		case *ast.Ident:
			recvName = astStartExpr.Name
		// 下标类型，实际上是泛型。泛型先不处理
		case *ast.IndexExpr:
			return nil
		}
	case *ast.Ident:
		recvName = recvType.Name
	default:
		recvName = funcDecl.Recv.List[0].Names[0].Name
	}
	rec := Param{
		Name: recvName,
	}
	switch typeType := funcDecl.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		switch astStartExpr := typeType.X.(type) {
		case *ast.Ident:
			rec.Type = astStartExpr.Name
			rec.AstType = enum.PARAM_AST_TYPE_StarExpr
		// 下标类型，实际上是泛型。泛型先不处理
		case *ast.IndexExpr:
			return nil
		}
	case *ast.Ident:
		rec.Type = typeType.Name
		rec.AstType = enum.PARAM_AST_TYPE_Ident
	}
	return &rec
}

func (s *SourceInfo) ParseGenericsMap(decl *ast.FuncDecl) map[string]*Param {
	if decl.Type == nil {
		return nil
	}
	if decl.Type.TypeParams == nil {
		return nil
	}
	if len(decl.Type.TypeParams.List) == 0 {
		return nil
	}
	// 1. 处理typeParam
	typeParams := decl.Type.TypeParams
	genericsMap := make(map[string]*Param, 10)

	if typeParams != nil && len(typeParams.List) > 0 {
		// 2. 一般只有一个
		field := typeParams.List[0]
		for _, v := range field.Names {
			ident, ok := field.Type.(*ast.Ident)
			if ok && ident.Name == "comparable" {
				genericsMap[v.Name] = lo.ToPtr(Param{
					Name:    ident.Name,
					Type:    "string",
					AstType: enum.PARAM_AST_TYPE_Ident,
				})
				continue
			}
			//v.Obj.Decl
			init := s.ParamParse(field.Type)
			genericsMap[v.Name] = init
		}
	}
	return genericsMap
}

// ParseParamRequest 解析请求
func (s *SourceInfo) ParseParamRequest(expr ast.Expr) []*ParamValue {
	genericsMapMap := s.GenericsMap
	// "_" 这种不处理了
	var db ParamValue

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.Type = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			s.AppendImportList(s.GetImportPathFromAliasMap(firstField))
		}
		db.AstType = enum.PARAM_AST_TYPE_SelectorExpr
	case *ast.Ident:
		result, ok := genericsMapMap[dbType.Name]
		if ok {
			db.Type = result.Type
		} else {
			db.Type = dbType.Name
		}
		db.AstType = enum.PARAM_AST_TYPE_Ident
		// 指针类型
	case *ast.StarExpr:
		param := s.ParamParse(dbType.X)
		db.Type = "*" + param.Type
		db.AstType = enum.PARAM_AST_TYPE_StarExpr
	case *ast.FuncType:
		paramType := s.parseFuncType(dbType)
		db.Type = paramType
		db.AstType = enum.PARAM_AST_TYPE_FuncType
	case *ast.InterfaceType:
		// 啥也不做
		db.Type = "interface{}"
		db.AstType = enum.PARAM_AST_TYPE_InterfaceType
	case *ast.ArrayType:
		requestType := s.parseParamArrayType(dbType)
		db.Type = requestType
		db.AstType = enum.PARAM_AST_TYPE_ArrayType
	case *ast.MapType:
		requestType := s.parseParamMapType(dbType)
		db.Type = requestType
		db.AstType = enum.PARAM_AST_TYPE_MapType
	// 可变长度，省略号表达式
	case *ast.Ellipsis:
		// 处理Elt
		param := s.ParamParse(dbType.Elt)
		db.Type = "[]" + param.Type
		db.AstType = enum.PARAM_AST_TYPE_Ellipsis
	case *ast.ChanType:
		// 处理value
		param := s.ParamParse(dbType.Value)
		if dbType.Dir == ast.RECV {
			db.Type = "<-chan " + param.Type
		} else {
			db.Type = "chan<- " + param.Type
		}
		db.AstType = enum.PARAM_AST_TYPE_ChanType
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		db.AstType = enum.PARAM_AST_TYPE_IndexExpr
		return nil
	case *ast.BinaryExpr:
		// 先直接取Y
		// todo 这里有个问题待解决， 	fmt.Print("convert str3 is: " + strconv.Itoa(str3))的解析会走到这里，这次先不处理后半部分的解析
		param := s.ParamParse(dbType.Y)
		db.Type = param.Type
		db.AstType = enum.PARAM_AST_TYPE_BinaryExpr
	case *ast.BasicLit:
		db.Type = dbType.Kind.String()
		db.AstType = enum.PARAM_AST_TYPE_BasicLit
		db.Value = dbType.Value
	case *ast.FuncLit:
		funcType := s.parseFuncType(dbType.Type)
		db.Type = funcType
		db.AstType = enum.PARAM_AST_TYPE_FuncLit
	case *ast.CompositeLit:
		init := s.ParamParse(dbType.Type)
		db.Type = init.Type
		db.AstType = enum.PARAM_AST_TYPE_CompositeLit
	case *ast.CallExpr:
		args := dbType.Args
		requests := make([]*ParamValue, 0, 10)
		for _, v := range args {
			request := s.ParseParamRequest(v)
			for _, s := range request {
				requests = append(requests, s)
			}
		}
		return requests
	default:
		panic("未知类型...")
	}
	if db.Name == "" && db.Type != "" {
		db.Name = db.Type
	}
	return []*ParamValue{&db}
}
