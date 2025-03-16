package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

func ExpressIdent(param *expr.Ident) *z3.AST {
	//elementList := []string{param.GetFormula()}
	//identMap := map[string]*expr.Ident{param.IdentName: param}

	//expression := &Z3Express{
	//ElementList: elementList,
	//IdentMap:    identMap,
	//Expr:        strings.Join(elementList, " "),
	//}
	//return []*Z3Express{expression}
	return nil
}

func ExpressTargetIdent(param *expr.Ident, targetParam _struct.Parameter) []*Z3Express {
	//elementList := []string{param.GetFormula()}
	//identMap := map[string]*expr.Ident{param.IdentName: param}

	expression := &Z3Express{
		//ElementList: elementList,
		//IdentMap:    identMap,
		//Expr:        targetParam.GetFormula() + " = " + param.GetFormula(),
	}
	return []*Z3Express{expression}
}
