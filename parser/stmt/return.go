package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Return return语句
// 用于从函数中返回值，一个函数可以有一个或多个return语句
type Return struct {
	ReturnList []_struct.Parameter
}

// todo 需要考虑return
func (r *Return) CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []bo.KeyFormula) []ConditionResult {
	return nil
}

// ParseReturn 解析ast
func ParseReturn(stmt *ast.ReturnStmt) *Return {
	r := &Return{}
	resultList := make([]_struct.Parameter, 0, 10)
	for _, v := range stmt.Results {
		result := expr.ParseParameter(v)
		resultList = append(resultList, result)
	}
	r.ReturnList = resultList
	return r
}

func (i *Return) ParseReturnCondition() []*ConditionNodeResult {
	return []*ConditionNodeResult{{
		ConditionNode: nil,
		IsBreak:       true,
	}}
}
