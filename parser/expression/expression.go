package expression

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"strings"
)

// ExpressDetail 公式
type ExpressDetail struct {
	Expr        string                    // "a > 10"
	ElementList []string                  // ["a",">","10"]
	IdentMap    map[string]*expr.Ident    // {"a":"a"}
	SelectorMap map[string]*expr.Selector // {"astSelector_a_b_c":"a.b.c"}
	CallMap     map[string]*expr.Call     // {"astCall_funca_c_d":funca(c,d)}
	BasicList   []*expr.BasicLit          // ["10"]
}

func Express(param _struct.Parameter) []*ExpressDetail {
	expressionList := make([]*ExpressDetail, 0, 10)
	eList := ExpressParam(param)
	expressionList = append(expressionList, eList...)
	return expressionList
}

// ExpressParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
func ExpressParam(param _struct.Parameter) []*ExpressDetail {
	switch exprType := param.(type) {
	case *expr.Binary:
		return ExpressBinary(exprType)
	case *expr.Unary:
		return ExpressUnary(exprType)
	case *expr.Parent:
		return ExpressParent(exprType)
	case *expr.Ident:
		return ExpressIdent(exprType)
	case *expr.Selector:
		return ExpressSelector(exprType)
	case *expr.Call:
		return ExpressCall(exprType)
	case *expr.BasicLit:
		return ExpressBasicLit(exprType)
	default:
		elementList := []string{param.GetFormula()}

		expression := &ExpressDetail{
			ElementList: elementList,
			Expr:        strings.Join(elementList, " "),
		}
		return []*ExpressDetail{expression}
	}
}
