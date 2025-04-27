package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"go/ast"
)

// Ident 基本的变量类型，结构简单，只需要 name、type、value
// var a,b,c int
type Ident struct {
	IdentName string
}

func (i *Ident) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (i *Ident) GetFormula() string {
	return i.IdentName
}

// ParseIdent 解析ast
func ParseIdent(expr *ast.Ident, _ bo.ExprContext) *Ident {
	return &Ident{
		IdentName: expr.Name,
	}
}
