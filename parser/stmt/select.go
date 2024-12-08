package stmt

import (
	"go/ast"
)

// Select select语句
// *ast.SelectStmt代表select语句。select语句用于在多个通信操作（通道的发送和接收）中进行选择，它会阻塞直到某个通信操作可以进行
type Select struct {
	Body *Block
}

func (s *Select) LogicExpression() []StatementAssignment {
	return nil
}

func (s *Select) CalculateCondition([]StatementAssignment) []ConditionResult {
	return nil
}

// ParseSelect 解析ast
func ParseSelect(stmt *ast.SelectStmt) *Select {
	s := &Select{}
	s.Body = ParseBlock(stmt.Body)
	return s
}
