package stmt

import (
	"go/ast"
)

// TypeSwitch *ast.TypeSwitchStmt用于表示类型switch语句
// 类型switch语句是一种特殊的switch语句，它主要用于根据接口值的实际类型来执行不同的分支代码
type TypeSwitch struct {
	Init   *Assign // 初始化语句，可以为空
	Assign Stmt    // 赋值语句，一般是 x := y.(type) or y.(type)
	Body   *Block
}

func (t *TypeSwitch) Express() []StatementExpression {
	stmtExpressionList := make([]StatementExpression, 0, 10)
	if t.Init != nil {
		initExpression := t.Init.Express()
		stmtExpressionList = append(stmtExpressionList, initExpression...)
	}
	// Assign有可能是x := y.(type) *Assign or y.(type) *Expr
	assign, ok := t.Assign.(*Assign)
	if ok {
		assignExpression := assign.Express()
		stmtExpressionList = append(stmtExpressionList, assignExpression...)
	}
	return stmtExpressionList
}

func (i *TypeSwitch) CalculateCondition([]StatementExpression) []ConditionResult {
	return nil
}

// ParseTypeSwitch 解析ast
func ParseTypeSwitch(stmt *ast.TypeSwitchStmt) *TypeSwitch {
	ts := &TypeSwitch{}

	if stmt.Init != nil {
		as, ok := stmt.Init.(*ast.AssignStmt)
		if !ok {
			panic("switch init type is not assign")
		}
		ts.Init = ParseAssign(as)
	}
	if stmt.Assign != nil {
		ts.Assign = ParseStmt(stmt.Assign)
	}
	ts.Body = ParseBlock(stmt.Body)

	return ts
}
