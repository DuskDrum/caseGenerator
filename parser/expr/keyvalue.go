package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
)

// KeyValue *ast.KeyValueExpr 用于表示键值对表达式。一般值得是字面量初始化CompositeLit
// 例如映射或结构体的初始化表达式，如 map[string]int{"a": 1, "b": 2} 或 Person{Name: "Alice", Age: 30} 中的每一个键值对。
type KeyValue struct {
	Key   _struct.Parameter
	Value _struct.Parameter
}

func (s *KeyValue) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_KEYVALUE
}

func (s *KeyValue) GetFormula() string {
	return s.Key.GetFormula() + " : " + s.Value.GetFormula()
}

// ParseKeyValue 解析ast
func ParseKeyValue(expr *ast.KeyValueExpr, af *ast.File) *KeyValue {
	kv := &KeyValue{}
	kv.Key = ParseParameter(expr.Key, af)
	kv.Value = ParseParameter(expr.Value, af)
	return kv
}
