package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	"go/token"
)

// ExpressUnary mocker Unary
func ExpressUnary(param *expr.Unary) *z3.AST {
	// 解析公式
	ast := ExpressParam(param.Content)

	if param.Op == token.NOT {
		return ast.Not()
	} else if param.Op == token.SUB {
		return ast.UnaryMinus()
	} else if param.Op == token.XOR { // ^ unary代表了取反
		return ast.BvNot()
	}

	return nil
}
