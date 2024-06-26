package vistitor

import (
	"caseGenerator/parse/bo"
	"go/ast"

	"github.com/samber/lo"
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
	typeParamsMap := make(map[string]*bo.ParamParseResult, 10)

	if typeParams != nil && len(typeParams.List) > 0 {
		// 1. 一般只有一个
		field := typeParams.List[0]
		for _, v := range field.Names {
			ident, ok := field.Type.(*ast.Ident)
			if ok && ident.Name == "comparable" {
				typeParamsMap[v.Name] = lo.ToPtr(bo.ParamParseResult{
					ParamName:      ident.Name,
					ParamType:      "string",
					ParamInitValue: "\"\"",
				})
				continue
			}
			//v.Obj.Decl
			init := ParseParamWithoutInit(field.Type, v.Name)
			typeParamsMap[v.Name] = init
		}
	}
	// 3. 设置到typeParam
	bo.SetTypeParamMap(typeParamsMap)
	return v
}
