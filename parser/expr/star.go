package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"

	"github.com/samber/lo"
)

// Star 指针类型，其实和其他类型类似、只是前面加了个*指针符号。用于表示指针类型或解引用操作
type Star struct {
	Child *_struct.Parameter
}

func (s *Star) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_STAR
}

func (s *Star) GetFormula() string {
	child := s.Child
	if child != nil {
		fromPtr := lo.FromPtr(child)
		return "*" + fromPtr.GetFormula()
	}
	return ""
}

// ParseStar 解析ast
func ParseStar(expr *ast.StarExpr) *Star {
	return &Star{
		Child: lo.ToPtr(ParseParameter(expr)),
	}
}
