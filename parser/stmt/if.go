package stmt

import (
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression"
	"caseGenerator/parser/mocker"
	_struct "caseGenerator/parser/struct"
	"fmt"
	"go/ast"
)

// If 条件语句
// 代码块是由花括号{}包围的一系列语句，它在函数体、控制结构（如if、for、switch）等场景中广泛使用。
type If struct {
	Init      *Assign
	Condition _struct.Parameter
	Block     *Block
	Else      Stmt
}

func (i *If) Express() []StatementExpression {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mocker 的值
	if i.Init != nil {
		return i.Init.Express()
	}
	return nil
}

func (i *If) CalculateCondition(seList []StatementExpression) []ConditionResult {
	// 1. 先拿到 Condition的表达式
	conditionExpressionList := expression.Express(i.Condition)
	// 2. 找表达式中的变量,去遍历找表达式中的变化记录
	for _, v := range conditionExpressionList {
		mockList := mocker.MockExpression(v, seList)
		if len(mockList) > 0 {
			fmt.Printf("mock结果列表: %v\n", mockList)
		}
	}

	return nil
}

// ParseIf 解析ast
func ParseIf(stmt *ast.IfStmt) *If {
	i := &If{}
	if stmt.Init != nil {
		as, ok := stmt.Init.(*ast.AssignStmt)
		if !ok {
			panic("switch init type is not assign")
		}
		i.Init = ParseAssign(as)
	}
	if stmt.Else != nil {
		i.Else = ParseStmt(stmt.Else)
	}
	i.Block = ParseBlock(stmt.Body)
	i.Condition = expr.ParseParameter(stmt.Cond)
	return i
}
