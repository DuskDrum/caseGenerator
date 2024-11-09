package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// If 代码块
// 代码块是由花括号{}包围的一系列语句，它在函数体、控制结构（如if、for、switch）等场景中广泛使用。
type If struct {
	Init      _struct.Stmt
	Condition _struct.Parameter
	Block     Block
	Else      _struct.Stmt
}

// ParseIf 解析ast
func ParseIf(stmt *ast.IfStmt) If {
	i := If{}
	if stmt.Init != nil {
		i.Init = ParseStmt(stmt.Init)
	}
	if stmt.Else != nil {
		i.Else = ParseStmt(stmt.Else)
	}
	i.Block = ParseBlock(stmt.Body)
	i.Condition = expr.ParseParameter(stmt.Cond)
	return i
}
