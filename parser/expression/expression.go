package expression

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/token"
	"strings"

	"github.com/Knetic/govaluate"
)

// Expression 公式
type Expression struct {
	Expr        string                    // "a > 10"
	ElementList []string                  // ["a",">","10"]
	IdentMap    map[string]*expr.Ident    // {"a":"a"}
	SelectorMap map[string]*expr.Selector // {"astSelector_a_b_c":"a.b.c"}
	CallMap     map[string]*expr.Call     // {"astCall_funca_c_d":funca(c,d)}
	BasicList   []*expr.BasicLit          // ["10"]
}

func Express(param _struct.Parameter) []*Expression {
	expressionList := make([]*Expression, 0, 10)
	eList := ExpressParam(param)
	expressionList = append(expressionList, eList...)
	return expressionList
}

// ExpressParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
func ExpressParam(param _struct.Parameter) []*Expression {
	switch exprType := param.(type) {
	case *expr.Binary:
		return ExpressBinary(exprType)
	case *expr.Unary:
		return ExpressUnary(exprType)
	case *expr.Parent:
		return ExpressParent(exprType)
	case *expr.Ident:
		elementList := []string{param.GetFormula()}
		identMap := map[string]*expr.Ident{exprType.IdentName: exprType}

		expression := &Expression{
			ElementList: elementList,
			IdentMap:    identMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *expr.Selector:
		key := strings.ReplaceAll(param.GetFormula(), ".", "_")
		selectorMap := map[string]*expr.Selector{"astSelector_" + key: exprType}
		elementList := []string{"astSelector_" + key}

		expression := &Expression{
			ElementList: elementList,
			SelectorMap: selectorMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *expr.Call:
		key := strings.ReplaceAll(param.GetFormula(), ".", "_")
		key = strings.ReplaceAll(key, "(", "_")
		key = strings.ReplaceAll(key, ")", "")
		key = strings.ReplaceAll(key, " ", "")

		callMap := map[string]*expr.Call{"astCall_" + key: exprType}

		elementList := []string{"astCall_" + key}

		expression := &Expression{
			ElementList: elementList,
			CallMap:     callMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *expr.BasicLit:
		elementList := []string{param.GetFormula()}
		basicList := []*expr.BasicLit{exprType}

		expression := &Expression{
			ElementList: elementList,
			BasicList:   basicList,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	default:
		elementList := []string{param.GetFormula()}

		expression := &Expression{
			ElementList: elementList,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	}
}

// ExpressBinary mock Binary
func ExpressBinary(param *expr.Binary) []*Expression {
	expressionList := make([]*Expression, 0, 10)

	// 如果类型是||逻辑或，剪枝只处理X
	// 解析X
	xExpressionList := ExpressParam(param.X)
	// 解析Y
	yExpressionList := ExpressParam(param.Y)

	// 解析Op
	// 逻辑或进行剪枝处理
	if param.Op == token.LOR {
		// 剪枝处理，只处理||左边的公式
		return xExpressionList
		// 如果类型是&&逻辑与，处理X和Y
	} else if param.Op == token.LAND {
		expressionList = append(expressionList, xExpressionList...)
		expressionList = append(expressionList, yExpressionList...)
		return expressionList
	}
	// 如果类型不是逻辑与或者逻辑或，那么手动组装Expression
	identMap := make(map[string]*expr.Ident, 10)
	callMap := make(map[string]*expr.Call, 10)
	selectorMap := make(map[string]*expr.Selector, 10)
	elementList := make([]string, 0, 10)
	basicList := make([]*expr.BasicLit, 0, 10)

	for _, v := range xExpressionList {
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

	opElement := expressRelationalOperators(param)
	elementList = append(elementList, opElement)

	for _, v := range yExpressionList {
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

	return []*Expression{{
		IdentMap:    identMap,
		BasicList:   basicList,
		ElementList: elementList,
		CallMap:     callMap,
		Expr:        strings.Join(elementList, " "),
		SelectorMap: selectorMap,
	}}
}

// 关系运算符 相关处理
func expressRelationalOperators(param *expr.Binary) string {
	if param.Op == token.EQL {
		// 等于，底层一定不是binary
		// 不是govaluate.EQ.String()==>"="
		return "=="
	} else if param.Op == token.LSS {
		return govaluate.LT.String()
	} else if param.Op == token.GTR {
		return govaluate.GT.String()
	} else if param.Op == token.LEQ {
		return govaluate.LTE.String()
	} else if param.Op == token.GEQ {
		return govaluate.GTE.String()
	}

	// 暴力破解
	if param.Op == token.ADD {
		return govaluate.PLUS.String()
	} else if param.Op == token.SUB {
		return govaluate.MINUS.String()
	} else if param.Op == token.MUL {
		return govaluate.MULTIPLY.String()
	} else if param.Op == token.QUO {
		return govaluate.DIVIDE.String()
	} else if param.Op == token.REM {
		return govaluate.MODULUS.String()
	}
	panic("express op illegal")
}

// ExpressUnary mock Unary
func ExpressUnary(param *expr.Unary) []*Expression {
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

	return []*Expression{{
		Expr:        strings.Join(elementList, " "),
		ElementList: elementList,
		IdentMap:    identMap,
		CallMap:     callMap,
		BasicList:   basicList,
		SelectorMap: selectorMap,
	}}
}

// ExpressParent mock Parent
func ExpressParent(param *expr.Parent) []*Expression {
	// 解析公式
	eList := ExpressParam(param.Content)

	elementList := make([]string, 0, 10)
	identMap := make(map[string]*expr.Ident, 10)
	callMap := make(map[string]*expr.Call, 10)
	basicList := make([]*expr.BasicLit, 0, 10)
	selectorMap := make(map[string]*expr.Selector, 10)

	elementList = append(elementList, "(")
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
	elementList = append(elementList, ")")

	return []*Expression{{
		Expr:        strings.Join(elementList, " "),
		ElementList: elementList,
		IdentMap:    identMap,
		CallMap:     callMap,
		BasicList:   basicList,
		SelectorMap: selectorMap,
	}}
}
