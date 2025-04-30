package stmt

import (
	"caseGenerator/parser/bo"
	"go/ast"
)

// CommClause *ast.CommClause主要用于表示select语句中的case子句
// select语句用于在多个通信操作（通道的发送和接收）中进行选择
type CommClause struct {
	Comm     Stmt
	BodyList []Stmt
}

// ParseCommClause 解析ast
func ParseCommClause(stmt *ast.CommClause, context bo.ExprContext) *CommClause {
	cc := &CommClause{}
	if stmt.Comm != nil {
		cc.Comm = ParseStmt(stmt.Comm, context)
	}
	bodyList := make([]Stmt, 0, 10)
	for _, v := range stmt.Body {
		ps := ParseStmt(v, context)
		bodyList = append(bodyList, ps)

	}
	cc.BodyList = bodyList

	return cc
}
