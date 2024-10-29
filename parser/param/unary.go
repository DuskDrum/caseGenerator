package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type Unary struct {
	BasicParam
}

func (s *Unary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_UNARY
}

func (s *Unary) GetInstance() Parameter {
	return s
}

func (s *Unary) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Unary) GetFormula() string {
	panic("implement me")
}

// ParseUnary 解析ast
func ParseUnary(expr *ast.UnaryExpr, name string) *Unary {
	return nil
}
