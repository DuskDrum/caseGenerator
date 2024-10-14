package param

import (
	"caseGenerator/common/enum"

	"github.com/samber/lo"
)

// Star 指针类型，其实和其他类型类似、只是前面加了个*指针符号
type Star struct {
	BasicParam
	BasicValue
	Child *Parameter
}

func (s *Star) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_STAR
}

func (s *Star) GetInstance() Parameter {
	return s
}

// GetZeroValue 指针类型的零值是
func (s *Star) GetZeroValue() Parameter {
	// 取 child 的零值，并且转为指针. (考虑是否直接返回 nil)
	child := lo.FromPtr(s.Child)
	s.Value = lo.ToPtr(child.GetZeroValue())
	return s
}

func (s *Star) GetFormula() string {
	child := s.Child
	if child != nil {
		fromPtr := lo.FromPtr(child)
		return "*" + fromPtr.GetFormula()
	}
	return ""
}
