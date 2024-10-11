package param

import (
	"caseGenerator/common/enum"
)

type Interface struct {
}

func (s *Interface) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INTERFACE
}

func (s *Interface) GetInstance() Parameter {
	return s
}

func (s *Interface) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Interface) GetFormula() string {
	panic("implement me")
}
