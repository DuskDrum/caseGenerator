package param

import (
	"caseGenerator/common/enum"
)

type Array struct {
	RecursionParam
	BasicValue
}

func (s *Array) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ARRAY
}

func (s *Array) GetInstance() Parameter {
	return s
}

func (s *Array) GetZeroValue() Parameter {
	s.Value = nil
	return s
}

func (s *Array) GetFormula() string {
	childFormula := s.Parameter.GetFormula()
	return "[]" + childFormula
}
