package param

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// Selector 选择类型，用来表示 a.b这种特殊类型。其中的 Child 用来表示对应的子类
type Selector struct {
	BasicParam
	SelectorParam
	SelectorName string
}

type SelectorParam struct {
	Child     *SelectorParam `json:"child,omitempty"`
	ParamName string
}

func (s *Selector) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_SELECTOR
}

func (s *Selector) GetInstance() Parameter {
	return s
}

func (s *Selector) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *Selector) GetFormula() string {
	panic("implement me")
}

// ParseSelector 解析ast
func ParseSelector(expr *ast.SelectorExpr, name string) *Selector {
	selector := Selector{
		BasicParam: BasicParam{
			ParameterType: enum.PARAMETER_TYPE_SELECTOR,
			Name:          name,
		},
	}
	selectorExpr, selectorName := GetRelationFromSelectorExpr(expr)
	selector.SelectorName = selectorName
	selector.SelectorParam = lo.FromPtr(selectorExpr)
	return &selector
}

func GetRelationFromSelectorExpr(se *ast.SelectorExpr) (*SelectorParam, string) {
	var sp = &SelectorParam{
		Child:     nil,
		ParamName: "",
	}
	if si, ok := se.X.(*ast.Ident); ok {
		sp.ParamName = si.Name
		return sp, si.Name + "." + se.Sel.Name
	}
	if sse, ok := se.X.(*ast.SelectorExpr); ok {
		expr, s := GetRelationFromSelectorExpr(sse)
		sp.Child = expr
		sp.ParamName = se.Sel.Name
		return sp, s + "." + se.Sel.Name
	}
	return sp, se.Sel.Name
}
