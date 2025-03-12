package z3

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

func ExpressBasicLit(param *expr.BasicLit) []*Z3Express {
	//elementList := []string{param.GetFormula()}
	//basicList := []*expr.BasicLit{param}

	expression := &Z3Express{
		//ElementList: elementList,
		//BasicList:   basicList,
		//Expr:        strings.Join(elementList, " "),
	}
	return []*Z3Express{expression}
}

func ExpressTargetBasicLit(param *expr.BasicLit, targetParam _struct.Parameter) []*Z3Express {
	//elementList := []string{param.GetFormula()}
	//basicList := []*expr.BasicLit{param}

	expression := &Z3Express{
		//ElementList: elementList,
		//BasicList:   basicList,
		//Expr:        targetParam.GetFormula() + " = " + param.GetFormula(),
	}
	return []*Z3Express{expression}
}
