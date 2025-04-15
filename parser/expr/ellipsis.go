package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
)

// Ellipsis *ast.Ellipsis 表示 Go 中的 ... 语法。主要有三个用法
// 用于函数参数列表: func example(args ...int) {}
// append 函数中的可变参数: append(slice, elements...)
// 不定长数组: [...]int
type Ellipsis struct {
	Param _struct.Parameter
}

func (s *Ellipsis) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ELLIPSIS
}

func (s *Ellipsis) GetFormula() string {
	return "[]" + s.Param.GetFormula()
}

// ParseEllipsis 解析ast
func ParseEllipsis(expr *ast.Ellipsis, af *ast.File) *Ellipsis {
	e := &Ellipsis{}
	e.Param = ParseParameter(expr.Elt, af)
	return e
}
