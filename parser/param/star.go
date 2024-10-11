package param

import (
	"caseGenerator/common/enum"
)

type Star struct {
}

func (s *Star) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_STAR
}

func (s *Star) GetInstance() Parameter {
	return s
}

func (s *Star) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Star) GetFormula() string {
	panic("implement me")
}
