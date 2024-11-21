package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Switch switch语句
// switch语句用于根据不同的条件执行不同的代码块
type Switch struct {
	Tag  _struct.Parameter // Tag是一个表达式（Expr），它是switch语句要判断的值。例如，在switch num中，num对应的表达式就是Tag
	Init *Assign
	Body *Block
}

func (s *Switch) Express() []StatementExpression {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mock 的值
	if s.Init != nil {
		return s.Init.Express()
	}
	return nil
}

func (i *Switch) CalculateCondition([]StatementExpression) []ConditionResult {
	return nil
}

// ParseSwitch 解析ast
func ParseSwitch(stmt *ast.SwitchStmt) *Switch {
	s := &Switch{}
	if stmt.Init != nil {
		as, ok := stmt.Init.(*ast.AssignStmt)
		if !ok {
			panic("switch init type is not assign")
		}
		s.Init = ParseAssign(as)
	}
	s.Tag = expr.ParseParameter(stmt.Tag)
	s.Body = ParseBlock(stmt.Body)
	return s
}
