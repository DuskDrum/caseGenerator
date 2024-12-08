package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Expr
// 表示一个表达式作为语句来使用。
// 例如，函数调用(fmt.Println(""))等在单独作为一条语句出现时，会被表示为ast.ExprStmt类型。
type Expr struct {
	Expr _struct.Parameter
}

func (e *Expr) LogicExpression() []StatementAssignment {
	return nil
}

// ParseExpr 解析ast
func ParseExpr(stmt *ast.ExprStmt) *Expr {
	exprStmt := &Expr{}
	if stmt.X != nil {
		parameter := expr.ParseParameter(stmt.X)
		if parameter != nil {
			exprStmt.Expr = parameter
		}
	}
	return exprStmt
}
