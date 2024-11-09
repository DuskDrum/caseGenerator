package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"

	"github.com/samber/lo"
)

// Map 映射类型
// key可能是:SelectorExpr、Ident、StarExpr、InterfaceType
// value可能是: 递归的map、递归的array等复杂类型
type Map struct {
	Value _struct.RecursionParam
	Key   _struct.Parameter
}

func (s *Map) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_MAP
}

func (s *Map) GetFormula() string {
	return "map[" + s.Key.GetFormula() + "]" + s.Value.GetFormula()
}

// ParseMap 解析ast
func ParseMap(expr *ast.MapType) *Map {
	m := &Map{}
	m.Key = ParseParameter(expr.Key)
	m.Value = lo.FromPtr(ParseRecursionValue(expr.Value))

	return m
}
