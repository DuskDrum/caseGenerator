package z3

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"strings"
)

func ExpressCall(param *expr.Call) []*Z3Express {
	key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	key = strings.ReplaceAll(key, "(", "_")
	key = strings.ReplaceAll(key, ")", "")
	key = strings.ReplaceAll(key, " ", "")

	//callMap := map[string]*expr.Call{"astCall_" + key: param}

	//elementList := []string{"astCall_" + key}

	expression := &Z3Express{
		//ElementList: elementList,
		//CallMap:     callMap,
		//Expr:        strings.Join(elementList, " "),
	}
	return []*Z3Express{expression}
}

func ExpressTargetCall(param *expr.Call, targetParam _struct.Parameter) []*Z3Express {
	key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	key = strings.ReplaceAll(key, "(", "_")
	key = strings.ReplaceAll(key, ")", "")
	key = strings.ReplaceAll(key, " ", "")

	//callMap := map[string]*expr.Call{"astCall_" + key: param}

	//elementList := []string{"astCall_" + key}

	expression := &Z3Express{
		//ElementList: elementList,
		//CallMap:     callMap,
		//Expr:        targetParam.GetFormula() + " = " + "astCall_" + key,
	}
	return []*Z3Express{expression}
}
