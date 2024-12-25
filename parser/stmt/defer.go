package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"

	"github.com/samber/lo"
)

// Defer defer语句。延迟函数
// 赋值语句可以是简单的变量赋值（如x = 5），也可以是多元赋值（如x, y = 1, 2）
type Defer struct {
	Call expr.Call
}

// ParseDefer 解析ast
func ParseDefer(stmt *ast.DeferStmt) *Defer {
	d := &Defer{}
	call := expr.ParseCall(stmt.Call)
	if call != nil {
		d.Call = lo.FromPtr(call)
	}
	return d
}
