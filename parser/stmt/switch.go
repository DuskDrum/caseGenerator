package stmt

import (
	"caseGenerator/common/utils"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression"
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

func (s *Switch) FormulaExpress() ([]bo.KeyFormula, map[string]*expr.Call) {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mocker 的值
	if s.Init != nil {
		return s.Init.FormulaExpress()
	}
	return nil, nil
}

func (s *Switch) CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []bo.KeyFormula) []ConditionResult {
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

// ParseSwitchCondition  每个switch-case天然就是else的，每个case只需要关心自己要对应的值即可
// default则不同，default需要把之前所有的case都取反
func (s *Switch) ParseSwitchCondition() []*ConditionNodeResult {
	results := make([]*ConditionNodeResult, 0, 10)
	caseDetailNodeList := make([]*ConditionNode, 0, 10)
	// 1. 首先解析tag
	tag := s.Tag
	// 2. 解析Body， Body是一个[]stmt.Stmt，里面都是*stmt.CaseClause
	for _, stmtDetail := range s.Body.StmtList {
		clause, ok := stmtDetail.(*CaseClause)
		if !ok {
			panic("can't parse switch stmtList item")
		}
		caseList := clause.CaseList
		bodyList := clause.BodyList
		// 遍历的解析caseList、bodyList
		if len(caseList) > 0 {
			for _, caseDetail := range caseList {
				list, node := parseBodyList(bodyList, tag, caseDetail)
				results = append(results, list...)
				caseDetailNodeList = append(caseDetailNodeList, node)
			}
		} else {
			// clause如果是没有CaseList的，代表是else。将caseDetailNodeList里的所有内容取反
			var cn *ConditionNode
			for _, v := range caseDetailNodeList {
				conditionNode := &ConditionNode{
					Condition:       v.Condition,
					ConditionResult: false,
				}
				cn = conditionNode.Offer(cn)
			}
			results = append(results, &ConditionNodeResult{
				ConditionNode: cn,
				IsBreak:       false,
			})
		}
	}
	return results
}

// parseBodyList 解析body列表
func parseBodyList(bodyList []Stmt, tag _struct.Parameter, caseDetail _struct.Parameter) ([]*ConditionNodeResult, *ConditionNode) {
	results := make([]*ConditionNodeResult, 0, 10)

	conditionExpressionList := expression.ExpressTargetParam(caseDetail, tag)
	cn := &ConditionNode{
		Condition:       conditionExpressionList,
		ConditionResult: true,
	}
	// 定义中间变量记录一个body中多个condition之间的关系
	middleNodeResultList := make([]*ConditionNodeResult, 0, 10)
	if len(middleNodeResultList) == 0 {
		results = append(results, &ConditionNodeResult{
			ConditionNode: cn,
			IsBreak:       false,
		})
	}
	for _, bodyDetail := range bodyList {
		// 遍历解析
		conditionResultList := ParseCondition(bodyDetail)
		if len(conditionResultList) != 0 {
			for _, middleConditionResultDetail := range conditionResultList {
				// 手动深拷贝
				sourceNode, err := utils.DeepCopyByJson(cn)
				if err != nil {
					panic(err.Error())
				}
				// 向后添加元素
				offerNode := sourceNode.Offer(middleConditionResultDetail.ConditionNode)
				// 返回出去
				results = append(results, &ConditionNodeResult{
					ConditionNode: offerNode,
					IsBreak:       false,
				})
			}
		}
	}

	return results, cn
}
