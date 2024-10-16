package param

import (
	"caseGenerator/common/enum"
)

type Interface struct {
	BasicParam
	BasicValue
}

func (s *Interface) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INTERFACE
}

func (s *Interface) GetInstance() Parameter {
	return s
}

func (s *Interface) GetZeroValue() Parameter {
	s.Value = nil
	return s
}

func (s *Interface) GetFormula() string {
	return "interface{}"
}
