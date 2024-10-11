package param

import (
	"caseGenerator/common/enum"
)

type Index struct {
}

func (s *Index) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INDEX
}

func (s *Index) GetInstance() Parameter {
	return s
}

func (s *Index) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Index) GetFormula() string {
	panic("implement me")
}
