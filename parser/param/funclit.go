package param

import (
	"caseGenerator/common/enum"
)

type FuncLit struct {
}

func (s *FuncLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNCLIT
}

func (s *FuncLit) GetInstance() Parameter {
	return s
}

func (s *FuncLit) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *FuncLit) GetFormula() string {
	panic("implement me")
}
