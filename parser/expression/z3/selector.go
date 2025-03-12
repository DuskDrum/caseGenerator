package z3

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

func ExpressSelector(param *expr.Selector) []*Z3Express {
	//key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	//selectorMap := map[string]*expr.Selector{"astSelector_" + key: param}
	//elementList := []string{"astSelector_" + key}

	expression := &Z3Express{
		//ElementList: elementList,
		//SelectorMap: selectorMap,
		//Expr:        strings.Join(elementList, " "),
	}
	return []*Z3Express{expression}
}

func ExpressTargetSelector(param *expr.Selector, targetParam _struct.Parameter) []*Z3Express {
	//key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	//selectorMap := map[string]*expr.Selector{"astSelector_" + key: param}
	//elementList := []string{"astSelector_" + key}

	expression := &Z3Express{
		//ElementList: elementList,
		//SelectorMap: selectorMap,
		//Expr:        targetParam.GetFormula() + " = " + "astSelector_" + key,
	}
	return []*Z3Express{expression}
}
