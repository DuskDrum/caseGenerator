package expr

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// Interface  表示接口类型的节点
type Interface struct {
	FieldList []Field //接口的方法列表
}

func (s *Interface) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INTERFACE
}

func (s *Interface) GetFormula() string {
	return "interface{}"
}

// ParseInterface 解析ast
func ParseInterface(expr *ast.InterfaceType) *Interface {
	i := &Interface{}
	if expr.Methods != nil {
		fieldList := make([]Field, 0, 10)
		for _, v := range expr.Methods.List {
			field := ParseField(v)
			if field != nil {
				fieldList = append(fieldList, lo.FromPtr(field))
			}
		}
		i.FieldList = fieldList
	}
	return i
}
