package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type Parent struct {
	BasicParam
}

func (s *Parent) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_PARENT
}

func (s *Parent) GetInstance() Parameter {
	return s
}

func (s *Parent) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Parent) GetFormula() string {
	panic("implement me")
}

// ParseParent 解析ast
func ParseParent(expr *ast.ParenExpr, name string) *Parent {
	return nil
}
