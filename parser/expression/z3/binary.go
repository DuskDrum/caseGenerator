package z3

import (
	"caseGenerator/parser/expr"
	"go/token"

	"github.com/Knetic/govaluate"
)

func ExpressBinary(param *expr.Binary) []*Z3Express {
	expressionList := make([]*Z3Express, 0, 10)

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
	expressDetails := make([]*Z3Express, 0, 10)
	relationList := expressRelation(param, xExpressionList, yExpressionList)
	expressDetails = append(expressDetails, relationList...)
	return expressDetails

}

func expressRelation(param *expr.Binary, xExpressionList []*Z3Express, yExpressionList []*Z3Express) []*Z3Express {
	// 如果类型不是逻辑与或者逻辑或，那么手动组装Expression
	//identMap := make(map[string]*expr.Ident, 10)
	//callMap := make(map[string]*expr.Call, 10)
	//selectorMap := make(map[string]*expr.Selector, 10)
	elementList := make([]string, 0, 10)
	//basicList := make([]*expr.BasicLit, 0, 10)
	//
	//for _, v := range xExpressionList {
	//	elementList = append(elementList, v.ElementList...)
	//	basicList = append(basicList, v.BasicList...)
	//	for mk, mv := range v.SelectorMap {
	//		selectorMap[mk] = mv
	//	}
	//	for mk, mv := range v.CallMap {
	//		callMap[mk] = mv
	//	}
	//	for mk, mv := range v.IdentMap {
	//		identMap[mk] = mv
	//	}
	//}

	opElement := expressRelationalOperators(param)
	elementList = append(elementList, opElement)

	//for _, v := range yExpressionList {
	//	elementList = append(elementList, v.ElementList...)
	//	basicList = append(basicList, v.BasicList...)
	//	for mk, mv := range v.SelectorMap {
	//		selectorMap[mk] = mv
	//	}
	//	for mk, mv := range v.CallMap {
	//		callMap[mk] = mv
	//	}
	//	for mk, mv := range v.IdentMap {
	//		identMap[mk] = mv
	//	}
	//}

	return []*Z3Express{{
		//IdentMap:    identMap,
		//BasicList:   basicList,
		//ElementList: elementList,
		//CallMap:     callMap,
		//Expr:        strings.Join(elementList, " "),
		//SelectorMap: selectorMap,
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

// z3.AST 关系运算符 相关处理
//func expressZ3RelationalOperators(param *expr.Binary) z3.AST {
//	if param.Op == token.EQL {
//		// 等于，底层一定不是binary
//		// 不是govaluate.EQ.String()==>"="
//		return "=="
//	} else if param.Op == token.LSS {
//		return govaluate.LT.String()
//	} else if param.Op == token.GTR {
//		return govaluate.GT.String()
//	} else if param.Op == token.LEQ {
//		return govaluate.LTE.String()
//	} else if param.Op == token.GEQ {
//		return govaluate.GTE.String()
//	}
//
//	// 暴力破解
//	if param.Op == token.ADD {
//		return govaluate.PLUS.String()
//	} else if param.Op == token.SUB {
//		return govaluate.MINUS.String()
//	} else if param.Op == token.MUL {
//		return govaluate.MULTIPLY.String()
//	} else if param.Op == token.QUO {
//		return govaluate.DIVIDE.String()
//	} else if param.Op == token.REM {
//		return govaluate.MODULUS.String()
//	}
//	panic("express op illegal")
//}
