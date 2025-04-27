package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"

	"github.com/samber/lo"
)

// Map 映射类型
// key可能是:SelectorExpr、Ident、StarExpr、InterfaceType
// value可能是: 递归的map、递归的array等复杂类型
type Map struct {
	ValueType _struct.RecursionParam
	KeyType   _struct.Parameter
}

func (m *Map) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_MAP
}

func (m *Map) GetFormula() string {
	return "map[" + m.KeyType.GetFormula() + "]" + m.ValueType.GetFormula()
}

// ParseMap 解析ast
func ParseMap(expr *ast.MapType, context bo.ExprContext) *Map {
	m := &Map{}
	m.KeyType = ParseParameter(expr.Key, context)
	m.ValueType = lo.FromPtr(ParseRecursionValue(expr.Value, context))

	return m
}
