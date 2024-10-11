package param

import (
	"caseGenerator/common/enum"
)

type Func struct {
}

func (s *Func) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNC
}

func (s *Func) GetInstance() Parameter {
	return s
}

func (s *Func) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Func) GetFormula() string {
	panic("implement me")
}
