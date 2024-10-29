package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type Map struct {
	BasicParam
}

func (s *Map) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_MAP
}

func (s *Map) GetInstance() Parameter {
	return s
}

func (s *Map) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Map) GetFormula() string {
	panic("implement me")
}

// ParseMap 解析ast
func ParseMap(expr *ast.MapType, name string) *Map {
	return nil
}
