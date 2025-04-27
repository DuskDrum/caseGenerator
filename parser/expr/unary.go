package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

type Unary struct {
	Op      token.Token
	Content _struct.Parameter
}

func (s *Unary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_UNARY
}

func (s *Unary) GetFormula() string {
	formula := ""
	formula += s.Op.String() + s.Content.GetFormula()
	return formula
}

// ParseUnary 解析ast
func ParseUnary(expr *ast.UnaryExpr, context bo.ExprContext) *Unary {
	u := &Unary{}
	u.Op = expr.Op
	u.Content = ParseParameter(expr.X, context)
	return u
}

// MockUnary UnaryParam
//func MockUnary(param *Unary) []Mock {
//	//mockList := make([]Mock, 0, 10)
//	// 如果类型是||逻辑或，剪枝只处理X
//	if param.Op == token.NOT {
//		return MockExpr(param.Content)
//		// 如果类型是&&逻辑与，处理X和Y
//	}
//
//	return nil
//}
