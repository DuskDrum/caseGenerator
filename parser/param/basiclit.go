package param

import (
	"caseGenerator/common/enum"
)

type BasicLit struct {
	BasicParam
	BasicValue
}

func (s *BasicLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BASICLIT
}

func (s *BasicLit) GetInstance() Parameter {
	return s
}

func (s *BasicLit) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *BasicLit) GetFormula() string {
	panic("implement me")
}
