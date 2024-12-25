package stmt

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
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

func (s *Switch) FormulaExpress() ([]bo.KeyFormula, map[string]enum.SpecificType) {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mocker 的值
	if s.Init != nil {
		return s.Init.FormulaExpress()
	}
	return nil, nil
}

func (s *Switch) CalculateCondition([]bo.StatementAssignment) []ConditionResult {
	return nil
}

// ParseSwitch 解析ast
// todo switch 里要考虑 else、嵌套if、嵌套switch、嵌套 type-switch之间的关系，也要考虑 return 直接跳出 condition
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
