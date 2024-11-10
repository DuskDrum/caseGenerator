package stmt

import (
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Block 代码块
// 代码块是由花括号{}包围的一系列语句，它在函数体、控制结构（如if、for、switch）等场景中广泛使用。
type Block struct {
	StmtList []_struct.Stmt
}

// ParseBlock 解析ast
func ParseBlock(stmt *ast.BlockStmt) Block {
	b := Block{}
	stmtList := make([]_struct.Stmt, 0, 10)
	for _, v := range stmt.List {
		p := ParseStmt(v)
		stmtList = append(stmtList, p)
	}
	b.StmtList = stmtList
	return b
}