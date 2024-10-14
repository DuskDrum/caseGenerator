package param

import (
	"caseGenerator/common/enum"
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
