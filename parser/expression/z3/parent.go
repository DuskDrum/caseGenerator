package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
)

// ExpressParent mocker Parent
func ExpressParent(param *expr.Parent) *z3.AST {
	// 解析子公式
	return ExpressParam(param.Content)
}
