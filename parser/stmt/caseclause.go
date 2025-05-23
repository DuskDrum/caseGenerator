package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// CaseClause 表示switch语句中的一个case子句
// 它包含了一个或多个值（用于匹配switch表达式的值）和一个语句块（当匹配成功时执行的代码）
type CaseClause struct {
	CaseList []_struct.Parameter //例如，在case 1, 2:中，1和2对应的表达式就在List中。list为nil那么代表了是default
	BodyList []Stmt
	Position token.Pos // 代码的行数，同一个文件里比对才有意义
}

// ParseCaseClause 解析ast
func ParseCaseClause(stmt *ast.CaseClause, context bo.ExprContext) *CaseClause {
	cc := &CaseClause{}
	cl := make([]_struct.Parameter, 0, 10)
	for _, v := range stmt.List {
		c := expr.ParseParameter(v, context)
		cl = append(cl, c)
	}
	bl := make([]Stmt, 0, 10)
	for _, b := range stmt.Body {
		ps := ParseStmt(b, context)
		bl = append(bl, ps)
	}
	cc.CaseList = cl
	cc.BodyList = bl
	cc.Position = stmt.Pos()
	return cc
}
