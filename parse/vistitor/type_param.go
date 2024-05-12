package vistitor

import (
	"caseGenerator/parse"
	"github.com/samber/lo"
	"go/ast"
)

// TypeParamVisitor typeParam
type TypeParamVisitor struct {
}

func (v *TypeParamVisitor) Visit(n ast.Node) ast.Visitor {
	decl, ok := n.(*ast.FuncDecl)
	if !ok {
		return v
	}
	if decl.Type == nil {
		return v
	}
	if decl.Type.TypeParams == nil {
		return v
	}
	if len(decl.Type.TypeParams.List) == 0 {
		return v
	}
	// 1. 处理typeParam
	typeParams := decl.Type.TypeParams
	typeParamsMap := make(map[string]*parse.ParamParseResult, 10)

	if typeParams != nil && len(typeParams.List) > 0 {
		// 1. 一般只有一个
		field := typeParams.List[0]
		for _, v := range field.Names {
			ident, ok := field.Type.(*ast.Ident)
			if ok && ident.Name == "comparable" {
				typeParamsMap[v.Name] = lo.ToPtr(parse.ParamParseResult{
					ParamName:      ident.Name,
					ParamType:      "string",
					ParamInitValue: "\"\"",
				})
				continue
			}
			//v.Obj.Decl
			init := parse.ParseParamWithoutInit(field.Type, v.Name)
			typeParamsMap[v.Name] = init
		}
	}
	// 3. 设置到typeParam
	parse.SetTypeParamMap(typeParamsMap)
	return v
}
