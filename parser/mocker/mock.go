package mocker

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expression/govaluate"
	mockresult2 "caseGenerator/parser/expression/govaluate/mockresult"
	"caseGenerator/parser/goValuate"
	"caseGenerator/parser/stmt"
	"strings"

	"github.com/samber/lo"
)

// IMock 规范mock执行的逻辑
type IMock interface {
}

// MockExpression  mocker 表达式
//  1. 如果basicLit有值，那么有靶子了，给其他变量赋值代入表达式中
//     如果是 int、float，可能是==、!=、>、<。那么找到 basicLit 的最大或者最小值，计算得到需要
//     如果是 string，那么可能是==或者!=，也有可能是>、<
//     如果是 nil， 可能是==或者!=。 nil是属于 Ident 里的
//
// 2. 如果两边都没有靶子，那么将其中一边设置为零值，再继续用第一步的流程(ident 不知道变量类型，所以没办法处理)
func MockExpression(expression *govaluate.ExpressDetail, seList []bo.StatementAssignment) []mockresult2.MockResult {
	// 1. 如果有 basicLit，那么按照 govalue进行解析试算得到最终结果
	if len(expression.BasicList) > 0 {
		return goValuate.MockBasicExpression(expression, seList)
	}
	// 2. 如果有nil，那么将其他属性都变成 nil
	resultList := make([]mockresult2.MockResult, 0, 10)
	if len(expression.IdentMap) > 0 {
		var nilTag bool
		for _, v := range expression.IdentMap {
			if strings.EqualFold(v.IdentName, "nil") {
				nilTag = true
			}
		}
		if nilTag {
			for _, v := range expression.IdentMap {
				if v.IdentName == "nil" {
					continue
				}
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult2.IdentMockResult{
						Ident: *v,
						// todo 需要获取零值
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult2.IdentMockResult{
						Ident:     *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
			for _, v := range expression.CallMap {
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult2.CallMockResult{
						Call: *v,
						// todo 需要获取零值
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult2.CallMockResult{
						Call:      *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
			for _, v := range expression.SelectorMap {
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult2.SelectorMockResult{
						Selector: *v,
						// todo 需要获取零值
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult2.SelectorMockResult{
						Selector:  *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
		}
		return resultList
	}
	// 3. 如果没有 basic 也没有 nil，那么先不处理

	return nil
}

// MockKeyFormula 根据传参赋值语句列表和 条件语句得到 需要 mock 的记录
// todo 先不考虑外部常量
func MockKeyFormula(condition *stmt.ConditionNodeResult) []mockresult2.MockResult {
	// 1. 解析出全部的判断条件
	nodeExpressesList := AnalysisConditionNode(condition.ConditionNode)
	formulaList := condition.KeyFormulaList
	// 2. 遍历判断条件开始执行公式
	positionSnapshot := 0
	for _, v := range nodeExpressesList {
		formulaSnapshotList := make([]bo.KeyFormula, 0, 10)
		for _, f := range formulaList {
			// 根据行数确认 assignment 和 condition
			if f.Position.Line < v.Position.Line && f.Position.Line > positionSnapshot {
				formulaSnapshotList = append(formulaSnapshotList, f)
			}
		}

	}

	// 全部的条件信息, 全部赋值信息
	// 1. 找出所有条件的参数，判断是不是在赋值语句里

	// 2. 如果在赋值语句里，那么就是局部变量，执行公式即可

	// 3. 如果不在赋值语句里，分情况讨论:
	// a. 外部常量，如果是 包.xxx，那么是外部常量, 需要找到外部变量是啥
	// b. 内部常量， 要先解析出来包中定义的常量在里面那么就是内部常量
	// c. 变量，未知数，要用算法逐一从公式中算出来，变量对应的可能是从请求传进来，也可能是从方法得到

	// 4.

	return nil

}

func AnalysisConditionNode(conditionNode *stmt.ConditionNode) []*stmt.ConditionNodeExpress {
	expressList := make([]*stmt.ConditionNodeExpress, 0, 10)

	var nodeSnapshot = conditionNode
	for {
		expressList = append(expressList, &stmt.ConditionNodeExpress{
			Condition:       nodeSnapshot.Condition,
			ConditionResult: nodeSnapshot.ConditionResult,
			Position:        nodeSnapshot.Position,
		})
		if nodeSnapshot.Relation == nil {
			break
		} else {
			nodeSnapshot = nodeSnapshot.Relation
		}
	}

	return expressList
}
