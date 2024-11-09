package stmt

import (
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// TypeSwitch *ast.TypeSwitchStmt用于表示类型switch语句
// 类型switch语句是一种特殊的switch语句，它主要用于根据接口值的实际类型来执行不同的分支代码
type TypeSwitch struct {
	Init   _struct.Stmt // 初始化语句，可以为空
	Assign _struct.Stmt // 赋值语句，一般是 x := y.(type) or y.(type)
	Body   Block
}

// ParseTypeSwitch 解析ast
func ParseTypeSwitch(stmt *ast.TypeSwitchStmt) TypeSwitch {
	ts := TypeSwitch{}

	if stmt.Init != nil {
		ts.Init = ParseStmt(stmt.Init)
	}
	if stmt.Assign != nil {
		ts.Assign = ParseStmt(stmt.Assign)
	}
	ts.Body = ParseBlock(stmt.Body)

	return ts
}
