package expr

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// FuncLit *ast.FuncLit 表示一个匿名函数（也称为函数字面量）
// 例如，func(a int) int { return a + 1 } 就是一个函数字面量
type FuncLit struct {
	FuncType FuncType
	// Body 匿名函数里面的函数 body 解析。具体类型需要等 combination结构定义之后
	// Body
}

func (s *FuncLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNCLIT
}

func (s *FuncLit) GetFormula() string {
	formula := ""
	formula = s.FuncType.GetFormula() + "{}"
	return formula
}

// ParseFuncLit 解析ast
func ParseFuncLit(expr *ast.FuncLit, af *ast.File) *FuncLit {
	fl := &FuncLit{}
	funcType := ParseFuncType(expr.Type, af)
	if funcType != nil {
		fl.FuncType = lo.FromPtr(funcType)
	}
	return fl
}
