package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"
)

// TypeAssert 类型断言表达式
// 用于断言接口类型的具体类型，例如 x.(T)
type TypeAssert struct {
	Content _struct.Parameter
	// 类型，
	Type _struct.Parameter
}

func (s *TypeAssert) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_TYPE_ASSERT
}

func (s *TypeAssert) GetFormula() string {
	formula := s.Content.GetFormula() + ".(" + s.Type.GetFormula() + ")"
	return formula
}

// ParseTypeAssert 解析ast
func ParseTypeAssert(expr *ast.TypeAssertExpr, context bo.ExprContext) *TypeAssert {
	typeAssert := &TypeAssert{}
	typeAssert.Content = ParseParameter(expr.X, context)
	typeAssert.Type = ParseParameter(expr.Type, context)
	return typeAssert
}
