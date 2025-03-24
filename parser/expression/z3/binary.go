package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	"go/token"
)

func ExpressBinary(param *expr.Binary) *z3.AST {
	// 解析X
	xExpression := ExpressParam(param.X)
	// 解析Y
	yExpression := ExpressParam(param.Y)

	// 解析Op
	if param.Op == token.LOR {
		return xExpression.Xor(yExpression)
		// 如果类型是&&逻辑与，处理X和Y
	} else if param.Op == token.LAND {
		return xExpression.And(yExpression)
	}

	// 下面是除了逻辑与、逻辑或的其他运算符
	return expressRelation(param, xExpression, yExpression)

}

func expressRelation(param *expr.Binary, xExpression *z3.AST, yExpression *z3.AST) *z3.AST {
	if param.Op == token.EQL {
		return xExpression.Eq(yExpression)
	} else if param.Op == token.LSS {
		return xExpression.Lt(yExpression)
	} else if param.Op == token.GTR {
		return xExpression.Gt(yExpression)
	} else if param.Op == token.LEQ {
		return xExpression.Le(yExpression)
	} else if param.Op == token.GEQ {
		return xExpression.Ge(yExpression)
	}

	// 需要判断类型是int还是float
	if param.Op == token.ADD {
		return xExpression.Add(yExpression)
	} else if param.Op == token.SUB {
		return xExpression.Sub(yExpression)
	} else if param.Op == token.MUL {
		return xExpression.Mul(yExpression)
	} else if param.Op == token.QUO {
		return xExpression.Div(yExpression)
	} else if param.Op == token.REM {
		return xExpression.Rem(yExpression)
	}
	panic("express op illegal")
}
