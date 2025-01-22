package expression

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
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

func ExpressTargetBasicLit(param *expr.BasicLit, targetParam _struct.Parameter) []*ExpressDetail {
	elementList := []string{param.GetFormula()}
	basicList := []*expr.BasicLit{param}

	expression := &ExpressDetail{
		ElementList: elementList,
		BasicList:   basicList,
		Expr:        targetParam.GetFormula() + " = " + param.GetFormula(),
	}
	return []*ExpressDetail{expression}
}
