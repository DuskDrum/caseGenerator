package param

import (
	"caseGenerator/common/enum"
)

// Selector 选择类型，用来表示 a.b这种特殊类型。其中的 Child 用来表示对应的子类
type Selector struct {
	BasicParam
	Child *Selector `json:"child,omitempty"`
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
