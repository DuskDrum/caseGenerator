package expression

import (
	"caseGenerator/parser/expr"
	"strings"
)

func ExpressCall(param *expr.Call) []*ExpressDetail {
	key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	key = strings.ReplaceAll(key, "(", "_")
	key = strings.ReplaceAll(key, ")", "")
	key = strings.ReplaceAll(key, " ", "")

	callMap := map[string]*expr.Call{"astCall_" + key: param}

	elementList := []string{"astCall_" + key}

	expression := &ExpressDetail{
		ElementList: elementList,
		CallMap:     callMap,
		Expr:        strings.Join(elementList, " "),
	}
	return []*ExpressDetail{expression}
}
