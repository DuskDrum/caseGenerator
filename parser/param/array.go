package param

import (
	"caseGenerator/common/enum"
)

type Array struct {
}

func (s *Array) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ARRAY
}

func (s *Array) GetInstance() Parameter {
	return s
}

func (s *Array) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Array) GetFormula() string {
	panic("implement me")
}
