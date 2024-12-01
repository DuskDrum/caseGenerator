package expression

import (
	"caseGenerator/parser/expr"
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
