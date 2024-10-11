package param

import (
	"caseGenerator/common/enum"
)

type Ident struct {
}

func (s *Ident) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *Ident) GetInstance() Parameter {
	return s
}

func (s *Ident) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Ident) GetFormula() string {
	panic("implement me")
}
