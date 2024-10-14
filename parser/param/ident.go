package param

import (
	"caseGenerator/common/enum"
)

// Ident 基本的变量类型，结构简单，只需要 name、type、value
type Ident struct {
	BasicParam
	BasicValue
	SpecificType
}

func (s *Ident) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *Ident) GetInstance() Parameter {
	return s
}

// GetZeroValue 返回具体类型对应的零值
func (s *Ident) GetZeroValue() Parameter {
	zeroValue := s.SpecificType.ZeroValue
	s.Value = zeroValue
	return s
}

func (s *Ident) GetFormula() string {
	return s.GetName()
}
