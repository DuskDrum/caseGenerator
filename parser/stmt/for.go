package stmt

import (
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

func (f *For) Express() []StatementExpression {
	return nil
}

// ParseFor 解析ast
func ParseFor(stmt *ast.ForStmt) *For {
	f := &For{}
	if stmt.Init != nil {
		f.Init = ParseStmt(stmt.Init)
	}
	if stmt.Post != nil {
		f.Post = ParseStmt(stmt.Post)
	}
	if stmt.Cond != nil {
		f.Cond = expr.ParseParameter(stmt.Cond)
	}
	if stmt.Body != nil {
		f.Body = ParseBlock(stmt.Body)
	}
	return f
}
