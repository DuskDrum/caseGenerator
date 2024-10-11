package param

import (
	"caseGenerator/common/enum"
)

type Map struct {
}

func (s *Map) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_MAP
}

func (s *Map) GetInstance() Parameter {
	return s
}

func (s *Map) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Map) GetFormula() string {
	panic("implement me")
}
