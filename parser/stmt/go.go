package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"

	"github.com/samber/lo"
)

// Go go关键字启动协程的语句
// 在 Go 代码中使用go关键字来并发地执行一个函数或者一个函数调用表达式时，这个语句在 AST 中就会被表示为*ast.GoStmt。
type Go struct {
	Call expr.Call
}

// ParseGo 解析ast
func ParseGo(stmt *ast.GoStmt, af *ast.File) *Go {
	g := &Go{}

	call := expr.ParseCall(stmt.Call, af)
	if call != nil {
		g.Call = lo.FromPtr(call)
	}

	return g
}
