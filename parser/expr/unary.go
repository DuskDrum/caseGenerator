package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

type Unary struct {
	Op      token.Token
	Content _struct.Parameter
}

func (s *Unary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_UNARY
}

func (s *Unary) GetFormula() string {
	formula := ""
	formula += s.Op.String() + s.Content.GetFormula()
	return formula
}

// ParseUnary 解析ast
func ParseUnary(expr *ast.UnaryExpr) *Unary {
	u := &Unary{}
	u.Op = expr.Op
	u.Content = ParseParameter(expr.X)
	return u
}
