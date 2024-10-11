package param

import (
	"caseGenerator/common/enum"
)

type Unary struct {
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
