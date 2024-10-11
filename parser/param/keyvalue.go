package param

import (
	"caseGenerator/common/enum"
)

type KeyValue struct {
}

func (s *KeyValue) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_KEYVALUE
}

func (s *KeyValue) GetInstance() Parameter {
	return s
}

func (s *KeyValue) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *KeyValue) GetFormula() string {
	panic("implement me")
}
