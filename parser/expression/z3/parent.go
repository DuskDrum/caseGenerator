package z3

import (
	"caseGenerator/parser/expr"
)

// ExpressParent mocker Parent
func ExpressParent(param *expr.Parent) []*Z3Express {
	// 解析公式
	//eList := ExpressParam(param.Content)

	elementList := make([]string, 0, 10)
	//identMap := make(map[string]*expr.Ident, 10)
	//callMap := make(map[string]*expr.Call, 10)
	//basicList := make([]*expr.BasicLit, 0, 10)
	//selectorMap := make(map[string]*expr.Selector, 10)

	elementList = append(elementList, "(")
	//for _, v := range eList {
	//elementList = append(elementList, v.ElementList...)
	//basicList = append(basicList, v.BasicList...)
	//for mk, mv := range v.SelectorMap {
	//	selectorMap[mk] = mv
	//}
	//for mk, mv := range v.CallMap {
	//	callMap[mk] = mv
	//}
	//for mk, mv := range v.IdentMap {
	//	identMap[mk] = mv
	//}
	//}
	// 解析括号
	elementList = append(elementList, ")")

	return []*Z3Express{{
		//Expr:        strings.Join(elementList, " "),
		//ElementList: elementList,
		//IdentMap:    identMap,
		//CallMap:     callMap,
		//BasicList:   basicList,
		//SelectorMap: selectorMap,
	}}
}
