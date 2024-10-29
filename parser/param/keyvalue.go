package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type KeyValue struct {
	BasicParam
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

// ParseKeyValue 解析ast
func ParseKeyValue(expr *ast.KeyValueExpr, name string) *KeyValue {
	return nil
}
