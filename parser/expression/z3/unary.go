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
		config := z3.NewConfig()
		ctx := z3.NewContext(config)
		return ctx.Int(0, ctx.IntSort()).Sub(ast)
	}

	return nil
}
