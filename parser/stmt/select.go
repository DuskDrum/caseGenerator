package stmt

import (
	"caseGenerator/parser/bo"
	"go/ast"
)

// Select select语句
// *ast.SelectStmt代表select语句。select语句用于在多个通信操作（通道的发送和接收）中进行选择，它会阻塞直到某个通信操作可以进行
type Select struct {
	Body *Block
}

// ParseSelect 解析ast
func ParseSelect(stmt *ast.SelectStmt, context bo.ExprContext) *Select {
	s := &Select{}
	s.Body = ParseBlock(stmt.Body, context)
	return s
}
