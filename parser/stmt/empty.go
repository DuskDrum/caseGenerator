package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"
)

// Empty 空语句.
// *ast.EmptyStmt 表示一个空语句。空语句通常用于表示分号（;）一个常见的场景是 for 循环中，语句之间的分隔。
// 比如在 for ; ; {} 中，两个分号表示空语句。还有可能在控制结构中使用空语句来占位，保证语法的完整性。
type Empty struct {
	expr.ValueSpec
}

func (e *Empty) LogicExpression() []StatementAssignment {
	return nil
}

func ParseEmpty(_ *ast.EmptyStmt) *Empty {
	empty := &Empty{}

	return empty
}
