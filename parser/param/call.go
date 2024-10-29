package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type Call struct {
	BasicParam
}

func (s *Call) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_CALL
}

func (s *Call) GetInstance() Parameter {
	return s
}

func (s *Call) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Call) GetFormula() string {
	panic("implement me")
}

// ParseCall 解析ast
func ParseCall(expr *ast.CallExpr, name string) *Call {
	return nil
}
