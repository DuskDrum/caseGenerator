package param

import (
	"caseGenerator/common/enum"
)

type Ellipsis struct {
	BasicParam
}

func (s *Ellipsis) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ELLIPSIS
}

func (s *Ellipsis) GetInstance() Parameter {
	return s
}

func (s *Ellipsis) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Ellipsis) GetFormula() string {
	panic("implement me")
}
