package stmt

import (
	"go/ast"
)

// CommClause *ast.CommClause主要用于表示select语句中的case子句
// select语句用于在多个通信操作（通道的发送和接收）中进行选择
type CommClause struct {
	Comm     Stmt
	BodyList []Stmt
}

func (c *CommClause) Express() []StatementExpression {
	return nil
}

// ParseCommClause 解析ast
func ParseCommClause(stmt *ast.CommClause) *CommClause {
	cc := &CommClause{}
	if stmt.Comm != nil {
		cc.Comm = ParseStmt(stmt.Comm)
	}
	bodyList := make([]Stmt, 0, 10)
	for _, v := range stmt.Body {
		ps := ParseStmt(v)
		bodyList = append(bodyList, ps)

	}
	cc.BodyList = bodyList

	return cc
}
