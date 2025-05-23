package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"go/token"
)

// ExpressUnary mocker Unary
func ExpressUnary(param *expr.Unary, context bo.ExpressionContext) (*z3.AST, []*z3.AST) {
	// 解析公式
	ast, _ := ExpressParam(param.Content, context)

	if param.Op == token.NOT {
		return ast.Not(), nil
	} else if param.Op == token.SUB {
		return ast.UnaryMinus(), nil
	} else if param.Op == token.XOR { // ^ unary代表了取反
		return ast.BvNot(), nil
	}

	return ast, nil
}
