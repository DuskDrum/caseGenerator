package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
	"strings"

	"github.com/samber/lo"
)

// Field
//  1. 函数参数和返回值：在 *ast.FuncType 的 Params 和 Results 字段中，*ast.Field 表示每个参数或返回值。
//  2. 结构体字段：在 *ast.StructType 的 Fields 中，每个 *ast.Field 表示结构体中的一个字段。
//  3. 接口方法：在 *ast.InterfaceType 的 Methods 中，每个 *ast.Field 表示一个接口方法的签名。
type Field struct {
	FiledNames []Ident
	Type       _struct.Parameter
	Tag        BasicLit
}

func (s *Field) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FIELD
}

// GetFormula 返回类似于这样 a, b int
func (s *Field) GetFormula() string {
	formulaNameList := make([]string, 0, 10)
	for _, v := range s.FiledNames {
		formulaNameList = append(formulaNameList, v.IdentName)
	}
	resultFormula := strings.Join(formulaNameList, ", ")

	resultFormula += " "
	resultFormula += s.Type.GetFormula()

	return resultFormula
}

// ParseField 解析ast
func ParseField(expr *ast.Field) *Field {
	identList := make([]Ident, 0, 10)
	if expr.Names != nil {
		for _, v := range expr.Names {
			ident := ParseIdent(v)
			if ident != nil {
				identList = append(identList, lo.FromPtr(ident))
			}
		}
	}

	filed := &Field{
		FiledNames: identList,
		Type:       ParseParameter(expr.Type),
	}
	tag := ParseBasicLit(expr.Tag)

	if tag != nil {
		filed.Tag = lo.FromPtr(tag)
	}

	return filed
}
