package expression

import (
	"caseGenerator/parser/expr"
	"strings"
)

func ExpressBasicLit(param *expr.BasicLit) []*ExpressDetail {
	elementList := []string{param.GetFormula()}
	basicList := []*expr.BasicLit{param}

	expression := &ExpressDetail{
		ElementList: elementList,
		BasicList:   basicList,
		Expr:        strings.Join(elementList, " "),
	}
	return []*ExpressDetail{expression}
}
