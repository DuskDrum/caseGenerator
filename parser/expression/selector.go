package expression

import (
	"caseGenerator/parser/expr"
	"strings"
)

func ExpressSelector(param *expr.Selector) []*ExpressDetail {
	key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	selectorMap := map[string]*expr.Selector{"astSelector_" + key: param}
	elementList := []string{"astSelector_" + key}

	expression := &ExpressDetail{
		ElementList: elementList,
		SelectorMap: selectorMap,
		Expr:        strings.Join(elementList, " "),
	}
	return []*ExpressDetail{expression}
}
