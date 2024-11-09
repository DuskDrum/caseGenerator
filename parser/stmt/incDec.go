package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// IncDec 自增、自减语句
// 表示自增（++）和自减（--）语句
type IncDec struct {
	Content _struct.Parameter // 变量，一般是ast.Ident类型代表的变量
	Token   token.Token       // 类型，是 token.INC、token.DEC
}

// ParseIncDec 解析ast
func ParseIncDec(stmt *ast.IncDecStmt) IncDec {
	incDec := IncDec{}

	if stmt.X != nil {
		xp := expr.ParseParameter(stmt.X)
		if xp != nil {
			incDec.Content = xp
		}
	}
	incDec.Token = stmt.Tok

	return incDec
}
