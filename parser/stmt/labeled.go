package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"
)

// Labeled 标签语句
// 主要用于goto语句和break、continue语句的跳转目标
type Labeled struct {
	Label expr.Ident // 标签，这里的这个标签和 *ast.BranchStmt 中的Label关联，比如label是outerLoop， 那么break outerLoop对应了跳到这里
	Block Stmt       // 标签后跟的语句
}

func (l *Labeled) LogicExpression() []StatementAssignment {
	return nil
}

func (l *Labeled) CalculateCondition([]StatementAssignment) []ConditionResult {
	return nil
}

func ParseLabeled(stmt *ast.LabeledStmt) *Labeled {
	labeled := &Labeled{}

	return labeled
}
