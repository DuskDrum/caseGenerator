package govaluate

import (
	"caseGenerator/parser/expr"
	"go/token"
	"strings"

	"github.com/Knetic/govaluate"
)

// ExpressUnary mocker Unary
func ExpressUnary(param *expr.Unary) []*ExpressDetail {
	// 解析公式
	eList := ExpressParam(param.Content)

	elementList := make([]string, 0, 10)
	identMap := make(map[string]*expr.Ident, 10)
	callMap := make(map[string]*expr.Call, 10)
	basicList := make([]*expr.BasicLit, 0, 10)
	selectorMap := make(map[string]*expr.Selector, 10)

	if param.Op == token.NOT {
		elementList = append(elementList, govaluate.INVERT.String())
	} else if param.Op == token.SUB {
		elementList = append(elementList, govaluate.NEGATE.String())
	}

	for _, v := range eList {
		elementList = append(elementList, v.ElementList...)
		basicList = append(basicList, v.BasicList...)
		for mk, mv := range v.SelectorMap {
			selectorMap[mk] = mv
		}
		for mk, mv := range v.CallMap {
			callMap[mk] = mv
		}
		for mk, mv := range v.IdentMap {
			identMap[mk] = mv
		}
	}
	// 解析括号

	return []*ExpressDetail{{
		Expr:        strings.Join(elementList, " "),
		ElementList: elementList,
		IdentMap:    identMap,
		CallMap:     callMap,
		BasicList:   basicList,
		SelectorMap: selectorMap,
	}}
}
