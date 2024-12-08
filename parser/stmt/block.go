package stmt

import (
	"go/ast"
)

// Block 代码块
// 代码块是由花括号{}包围的一系列语句，它在函数体、控制结构（如if、for、switch）等场景中广泛使用。
type Block struct {
	StmtList []Stmt
}

func (b *Block) LogicExpression() []StatementAssignment {
	return nil
}

func (b *Block) CalculateCondition([]StatementAssignment) []ConditionResult {
	return nil
}

// ParseBlock 解析ast
func ParseBlock(stmt *ast.BlockStmt) *Block {
	b := &Block{}
	stmtList := make([]Stmt, 0, 10)
	for _, v := range stmt.List {
		p := ParseStmt(v)
		stmtList = append(stmtList, p)
	}
	b.StmtList = stmtList
	return b
}
