package param

import (
	"caseGenerator/common/enum"
)

type Selector struct {
}

func (s *Selector) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_SELECTOR
}

func (s *Selector) GetInstance() Parameter {
	return s
}

func (s *Selector) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Selector) GetFormula() string {
	panic("implement me")
}
