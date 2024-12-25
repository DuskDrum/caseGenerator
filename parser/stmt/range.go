package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// Range 范围语句
// 比如for range循环语句
type Range struct {
	Key     _struct.Parameter // 有可能为空，为_
	Value   _struct.Parameter // 有可能为空，为_
	Token   token.Token       // token.ASSIGN, token.DEFINE
	Content _struct.Parameter // range后跟着的部分，比如说range []int{}。 这个就代表了[]int{}
	Body    *Block
}

// ParseRange 解析ast
func ParseRange(stmt *ast.RangeStmt) *Range {
	r := &Range{}
	if stmt.Key != nil {
		r.Key = expr.ParseParameter(stmt.Key)
	}
	if stmt.Value != nil {
		r.Value = expr.ParseParameter(stmt.Value)
	}
	if stmt.X != nil {
		r.Content = expr.ParseParameter(stmt.X)
	}
	r.Token = stmt.Tok
	r.Body = ParseBlock(stmt.Body)

	return r
}
