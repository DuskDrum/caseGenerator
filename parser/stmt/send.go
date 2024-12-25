package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// Send 发送语句
// 用于在 Go 的通道（channel）操作中，将一个值发送到通道里
type Send struct {
	Chan  _struct.Parameter // 通道表达式, 是发送操作所针对的通道。
	Value _struct.Parameter // 表示要发送到通道的值
	Arrow token.Pos         // 通道的方向 token.Arrow
}

// ParseSend 解析ast
func ParseSend(stmt *ast.SendStmt) *Send {
	send := &Send{}
	if stmt.Chan != nil {
		cp := expr.ParseParameter(stmt.Chan)
		if cp != nil {
			send.Chan = cp
		}
	}
	if stmt.Value != nil {
		vp := expr.ParseParameter(stmt.Value)
		if vp != nil {
			send.Value = vp
		}
	}
	send.Arrow = stmt.Arrow
	return send
}
