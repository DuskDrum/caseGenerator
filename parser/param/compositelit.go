package param

import (
	"caseGenerator/common/enum"
)

type CompositeLit struct {
	BasicParam
}

func (s *CompositeLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_COMPOSITELIT
}

func (s *CompositeLit) GetInstance() Parameter {
	return s
}

func (s *CompositeLit) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *CompositeLit) GetFormula() string {
	panic("implement me")
}
