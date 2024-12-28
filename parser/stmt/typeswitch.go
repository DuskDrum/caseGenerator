package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"go/ast"
)

// TypeSwitch *ast.TypeSwitchStmt用于表示类型switch语句
// 类型switch语句是一种特殊的switch语句，它主要用于根据接口值的实际类型来执行不同的分支代码
type TypeSwitch struct {
	Init   *Assign // 初始化语句，可以为空
	Assign Stmt    // 赋值语句，一般是 x := y.(type) or y.(type)
	Body   *Block
}

func (t *TypeSwitch) FormulaExpress() ([]bo.KeyFormula, map[string]*expr.Call) {
	keyFormulaList := make([]bo.KeyFormula, 0, 10)
	callMap := make(map[string]*expr.Call, 10)
	if t.Init != nil {
		initF, initOuter := t.Init.FormulaExpress()
		keyFormulaList = append(keyFormulaList, initF...)
		for k, v := range initOuter {
			callMap[k] = v
		}
	}
	// Assign有可能是x := y.(type) *FormulaExpress or y.(type) *Expr
	assign, ok := t.Assign.(*Assign)
	if ok {
		aF, aOuter := assign.FormulaExpress()
		keyFormulaList = append(keyFormulaList, aF...)
		for k, v := range aOuter {
			callMap[k] = v
		}
	}
	return keyFormulaList, callMap
}

func (t *TypeSwitch) CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []bo.KeyFormula) []ConditionResult {
	return nil
}

// ParseTypeSwitch 解析ast
func ParseTypeSwitch(stmt *ast.TypeSwitchStmt) *TypeSwitch {
	ts := &TypeSwitch{}

	if stmt.Init != nil {
		as, ok := stmt.Init.(*ast.AssignStmt)
		if !ok {
			panic("switch init type is not assign")
		}
		ts.Init = ParseAssign(as)
	}
	if stmt.Assign != nil {
		ts.Assign = ParseStmt(stmt.Assign)
	}
	ts.Body = ParseBlock(stmt.Body)

	return ts
}

func (i *TypeSwitch) ParseTypeSwitchCondition() []*ConditionNodeResult {

	return nil
}
