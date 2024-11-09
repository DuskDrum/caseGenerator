package expr

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// Selector 选择类型，用来表示 a.b这种特殊类型。其中的 Child 用来表示对应的子类
type Selector struct {
	SelectorParam
}

type SelectorParam struct {
	Child *SelectorParam `json:"child,omitempty"`
	Ident *Ident
}

func (s *Selector) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_SELECTOR
}

func (s *Selector) GetFormula() string {
	formula := ""
	current := s.SelectorParam
	// 递归遍历每一层的 child，直到没有 child
	for current.Ident != nil {
		formula = current.Ident.GetFormula() + "." + formula
		if current.Child == nil {
			break
		} else {
			current = lo.FromPtr(current.Child)
		}
	}
	return formula
}

// ParseSelector 解析ast
func ParseSelector(expr *ast.SelectorExpr) *Selector {
	selector := Selector{}
	selectorExpr := GetRelationFromSelectorExpr(expr)
	selector.SelectorParam = lo.FromPtr(selectorExpr)
	return &selector
}

func GetRelationFromSelectorExpr(se *ast.SelectorExpr) *SelectorParam {
	var sp = &SelectorParam{}
	if si, ok := se.X.(*ast.Ident); ok {
		sp.Child = nil
		sp.Ident = ParseIdent(si)

		return sp
	}
	if sse, ok := se.X.(*ast.SelectorExpr); ok {
		sp.Child = GetRelationFromSelectorExpr(sse)
		sp.Ident = ParseIdent(sse.Sel)
		return sp
	}
	return sp
}
