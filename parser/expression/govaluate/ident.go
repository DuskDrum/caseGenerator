package govaluate

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"strings"
)

func ExpressIdent(param *expr.Ident) []*ExpressDetail {
	elementList := []string{param.GetFormula()}
	identMap := map[string]*expr.Ident{param.IdentName: param}

	expression := &ExpressDetail{
		ElementList: elementList,
		IdentMap:    identMap,
		Expr:        strings.Join(elementList, " "),
	}
	return []*ExpressDetail{expression}
}

func ExpressTargetIdent(param *expr.Ident, targetParam _struct.Parameter) []*ExpressDetail {
	elementList := []string{param.GetFormula()}
	identMap := map[string]*expr.Ident{param.IdentName: param}

	expression := &ExpressDetail{
		ElementList: elementList,
		IdentMap:    identMap,
		Expr:        targetParam.GetFormula() + " = " + param.GetFormula(),
	}
	return []*ExpressDetail{expression}
}
