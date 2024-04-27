package parse

import (
	"go/ast"
	"log"
)

type RequestAssertionTypeVisitor struct {
}

func (v *RequestAssertionTypeVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	switch node := n.(type) {
	case *ast.TypeAssertExpr:
		var tau TypeAssertUnary
		switch tae := node.X.(type) {
		case *ast.Ident:
			tau.ParamValue = tae.Name
		case *ast.SelectorExpr:
			tau.ParamValue = GetRelationFromSelectorExpr(tae)
		default:
			log.Fatalf("不支持此类型")
		}
		switch nr := node.Type.(type) {
		case *ast.Ident:
			tau.AssertType = nr.Name
		case *ast.SelectorExpr:
			tau.AssertType = GetRelationFromSelectorExpr(nr)
		default:
			log.Fatalf("不支持此类型")
		}
	}
	return v
}
