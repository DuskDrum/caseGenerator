package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Switch switch语句
// switch语句用于根据不同的条件执行不同的代码块
type Switch struct {
	Tag  _struct.Parameter // Tag是一个表达式（Expr），它是switch语句要判断的值。例如，在switch num中，num对应的表达式就是Tag
	Init _struct.Stmt
	Body Block
}

// ParseSwitch 解析ast
func ParseSwitch(stmt *ast.SwitchStmt) Switch {
	s := Switch{}
	s.Init = ParseStmt(stmt.Init)
	s.Tag = expr.ParseParameter(stmt.Tag)
	s.Body = ParseBlock(stmt.Body)
	return s
}
