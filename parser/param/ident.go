package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// Ident 基本的变量类型，结构简单，只需要 name、type、value
// var a,b,c int
type Ident struct {
	BasicParam
	BasicValue
	IdentName string
}

func (s *Ident) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *Ident) GetInstance() Parameter {
	return s
}

// GetZeroValue 返回具体类型对应的零值
func (s *Ident) GetZeroValue() Parameter {
	s.Value = nil
	return s
}

func (s *Ident) GetFormula() string {
	return s.GetName()
}

// ParseIdent 解析ast
func ParseIdent(expr *ast.Ident, name string) *Ident {
	return &Ident{
		BasicParam: BasicParam{
			ParameterType: enum.PARAMETER_TYPE_INDEX,
			Name:          name,
		},
		IdentName: expr.Name,
	}
}
