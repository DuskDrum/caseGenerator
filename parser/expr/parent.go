package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"
)

// Parent 带括号的表达式。这种表达式主要用于改变运算的优先级，确保括号内的表达式先进行计算。
// 比如说(a+b)
type Parent struct {
	Content _struct.Parameter
}

func (s *Parent) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_PARENT
}

func (s *Parent) GetFormula() string {
	formula := "("
	formula += s.Content.GetFormula()
	formula += ")"
	return formula
}

// ParseParent 解析ast
func ParseParent(expr *ast.ParenExpr, context bo.ExprContext) *Parent {
	p := &Parent{}
	parameter := ParseParameter(expr.X, context)
	p.Content = parameter
	return p
}
