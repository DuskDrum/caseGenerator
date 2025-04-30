package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// For for循环语句
// 重复执行一段代码，直到满足特定的条件。
type For struct {
	Init Stmt              // 初始化
	Cond _struct.Parameter // 条件
	Post Stmt              //每次循环迭代后执行的语句， 比如i++
	Body *Block
}

// ParseFor 解析ast
func ParseFor(stmt *ast.ForStmt, context bo.ExprContext) *For {
	f := &For{}
	if stmt.Init != nil {
		f.Init = ParseStmt(stmt.Init, context)
	}
	if stmt.Post != nil {
		f.Post = ParseStmt(stmt.Post, context)
	}
	if stmt.Cond != nil {
		f.Cond = expr.ParseParameter(stmt.Cond, context)
	}
	if stmt.Body != nil {
		f.Body = ParseBlock(stmt.Body, context)
	}
	return f
}
