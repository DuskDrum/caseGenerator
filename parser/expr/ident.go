package expr

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// Ident 基本的变量类型，结构简单，只需要 name、type、value
// var a,b,c int
type Ident struct {
	IdentName string
}

func (s *Ident) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *Ident) GetFormula() string {
	return s.IdentName
}

// ParseIdent 解析ast
func ParseIdent(expr *ast.Ident) *Ident {
	return &Ident{
		IdentName: expr.Name,
	}
}
