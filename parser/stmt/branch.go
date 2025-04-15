package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"
	"go/token"

	"github.com/samber/lo"
)

// Branch 分支语句
// 分支语句主要是指break、continue和goto这几种可以改变程序控制流的语句
type Branch struct {
	Label expr.Ident  // 标签，这里的这个标签和 *ast.BranchStmt 中的Label关联，比如label是outerLoop， 那么break outerLoop对应了跳到这里
	Token token.Token // 标签的类型，是 token.BREAK、token.CONTINUE、token.GOTO等
}

// ParseBranch 解析ast
func ParseBranch(stmt *ast.BranchStmt, af *ast.File) *Branch {
	branch := &Branch{}
	if stmt.Label != nil {
		label := expr.ParseIdent(stmt.Label, af)
		if label != nil {
			branch.Label = lo.FromPtr(label)
		}
	}
	branch.Token = stmt.Tok

	return branch
}
