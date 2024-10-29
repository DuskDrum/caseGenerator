package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type Chan struct {
	BasicParam
}

func (s *Chan) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_CHAN
}

func (s *Chan) GetInstance() Parameter {
	return s
}

func (s *Chan) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Chan) GetFormula() string {
	panic("implement me")
}

// ParseChan 解析ast
func ParseChan(expr *ast.ChanType, name string) *Chan {
	return nil
}
