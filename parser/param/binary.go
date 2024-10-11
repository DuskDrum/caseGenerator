package param

import (
	"caseGenerator/common/enum"
)

type Binary struct {
}

func (s *Binary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BINARY
}

func (s *Binary) GetInstance() Parameter {
	return s
}

func (s *Binary) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Binary) GetFormula() string {
	panic("implement me")
}
