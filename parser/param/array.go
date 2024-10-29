package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// Array 数组结构
type Array struct {
	RecursionParam
	BasicValue
}

func (s *Array) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ARRAY
}

func (s *Array) GetInstance() Parameter {
	return s
}

func (s *Array) GetZeroValue() Parameter {
	s.Value = nil
	return s
}

func (s *Array) GetFormula() string {
	childFormula := s.Parameter.GetFormula()
	return "[]" + childFormula
}

// ParseArray 解析ast
func ParseArray(expr *ast.ArrayType, name string) *Array {
	return nil
}
