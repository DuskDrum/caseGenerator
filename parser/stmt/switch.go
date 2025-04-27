package stmt

import (
	"caseGenerator/common/utils"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression/govaluate"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// Switch switch语句
// switch语句用于根据不同的条件执行不同的代码块
type Switch struct {
	Tag            _struct.Parameter // Tag是一个表达式（Expr），它是switch语句要判断的值。例如，在switch num中，num对应的表达式就是Tag
	Init           *Assign
	CaseClauseList []*CaseClause // case列表
	DefaultCase    *CaseClause   // 默认的case列表，List值为nil
	Position       token.Pos     // 代码的行数，同一个文件里比对才有意义
}

func (s *Switch) FormulaExpress() ([]govaluate.KeyFormula, map[string]*expr.Call) {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mocker 的值
	if s.Init != nil {
		return s.Init.FormulaExpress()
	}
	return nil, nil
}

func (s *Switch) CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []govaluate.KeyFormula) []ConditionResult {
	return nil
}

// ParseSwitch 解析ast
func ParseSwitch(stmt *ast.SwitchStmt, af *ast.File) *Switch {
	s := &Switch{}
	if stmt.Init != nil {
		as, ok := stmt.Init.(*ast.AssignStmt)
		if !ok {
			panic("switch init type is not assign")
		}
		s.Init = ParseAssign(as, af)
	}
	s.Tag = expr.ParseParameter(stmt.Tag, af)
	s.Position = stmt.Pos()

	caseClauseList := make([]*CaseClause, 0, 10)
	defaultCaseList := make([]*CaseClause, 0, 10)
	// 解析 body, body是ast.BlockStmt类型，其中元素都应该是*ast.CaseClause,否则不合法
	for _, v := range stmt.Body.List {
		clause, ok := v.(*ast.CaseClause)
		if !ok {
			panic("switch clause type is not case")
		}
		if clause.List != nil {
			caseClauseList = append(caseClauseList, ParseCaseClause(clause, af))
		} else {
			defaultCaseList = append(defaultCaseList, ParseCaseClause(clause, af))
		}
	}
	// default 只能有一个，解析出多个就报错
	if len(defaultCaseList) > 1 {
		panic("default clause list size more than 1")
	}
	if len(defaultCaseList) == 1 {
		s.DefaultCase = defaultCaseList[0]
	}
	s.CaseClauseList = caseClauseList
	return s
}

// ParseSwitchCondition  每个switch-case天然就是else的，每个case只需要关心自己要对应的值即可
// todo 应该把KeyFormulaList打到每个条件里
// default则不同，default需要把之前所有的case都取反
func (s *Switch) ParseSwitchCondition() []*ConditionNodeResult {
	results := make([]*ConditionNodeResult, 0, 10)
	uncleNodeList := make([]*ConditionNode, 0, 10)
	// 1. 首先解析tag
	tag := s.Tag
	// 2. 解析caseList, 每个元素都是*stmt.CaseClause
	for _, clause := range s.CaseClauseList {
		caseList := clause.CaseList
		bodyList := clause.BodyList
		// 遍历的解析caseList、bodyList
		for _, caseDetail := range caseList {
			list, node := parseBodyList(bodyList, tag, caseDetail)
			results = append(results, list...)
			uncleNodeList = append(uncleNodeList, node)
		}
	}
	// 3. 解析 defaultCase
	// clause如果是没有CaseList的，代表是else。将caseDetailNodeList里的所有内容取反
	if s.DefaultCase != nil {
		defaultTag := parseDefaultTag(uncleNodeList, s.DefaultCase)
		results = append(results, defaultTag...)
	}
	return results
}

func parseDefaultTag(uncleNodeList []*ConditionNode, defaultClause *CaseClause) []*ConditionNodeResult {
	results := make([]*ConditionNodeResult, 0, 10)
	// 1. 给 uncle 所有条件取反
	var cn *ConditionNode
	for _, v := range uncleNodeList {
		conditionNode := &ConditionNode{
			Condition:       v.Condition,
			ConditionResult: !v.ConditionResult,
		}
		cn = conditionNode.Offer(cn)
	}
	// 2.记录赋值键值对列表
	keyFormulaList := make([]govaluate.KeyFormula, 0, 10)
	formulaCallMap := make(map[string]*expr.Call, 10)
	// 3. 解析 default 中的每个条件
	for _, bodyDetail := range defaultClause.BodyList {
		// 4. 解析每一句是否是赋值语句，并将其存储起来
		formulaList, callMap := ParseConditionKeyFormula(bodyDetail)
		keyFormulaList = append(keyFormulaList, formulaList...)
		utils.PutAll(formulaCallMap, callMap)
		// 遍历解析
		conditionResultList := ParseCondition(bodyDetail)
		if len(conditionResultList) != 0 {
			for _, middleConditionResultDetail := range conditionResultList {
				// 手动深拷贝
				sourceNode, err := utils.DeepCopyByJson(cn)
				if err != nil {
					panic(err.Error())
				}
				// 组装赋值语句
				copyFormulaList := utils.CopySlice(keyFormulaList)
				copyFormulaList = append(copyFormulaList, middleConditionResultDetail.KeyFormulaList...)

				copyCallMap := utils.CopyMap(callMap)
				// 向后添加元素
				offerNode := sourceNode.Offer(middleConditionResultDetail.ConditionNode)
				// 返回出去
				results = append(results, &ConditionNodeResult{
					ConditionNode:  offerNode,
					IsBreak:        middleConditionResultDetail.IsBreak,
					KeyFormulaList: copyFormulaList,
					FormulaCallMap: copyCallMap,
				})
			}
		}
	}
	// 如果 default 中没有子条件，那么返回所有叔节点取反
	if len(results) == 0 {
		results = append(results, &ConditionNodeResult{
			ConditionNode: cn,
			IsBreak:       false,
		})
	}
	return results
}

// parseBodyList 解析body列表
func parseBodyList(bodyList []Stmt, tag _struct.Parameter, caseDetail _struct.Parameter) ([]*ConditionNodeResult, *ConditionNode) {
	results := make([]*ConditionNodeResult, 0, 10)
	// 1. 解析出 tag
	conditionExpressionList := govaluate.ExpressTargetParam(caseDetail, tag)
	cn := &ConditionNode{
		Condition:       conditionExpressionList,
		ConditionResult: true,
	}

	// 2.记录赋值键值对列表
	keyFormulaList := make([]govaluate.KeyFormula, 0, 10)
	formulaCallMap := make(map[string]*expr.Call, 10)

	for _, bodyDetail := range bodyList {
		// 3. 解析每一句是否是赋值语句，并将其存储起来
		formulaList, callMap := ParseConditionKeyFormula(bodyDetail)
		keyFormulaList = append(keyFormulaList, formulaList...)
		utils.PutAll(formulaCallMap, callMap)

		// 4.遍历解析条件语句
		conditionResultList := ParseCondition(bodyDetail)
		if len(conditionResultList) != 0 {
			for _, middleConditionResultDetail := range conditionResultList {
				// 手动深拷贝
				sourceNode, err := utils.DeepCopyByJson(cn)
				if err != nil {
					panic(err.Error())
				}
				// 组装赋值语句
				// 深度拷贝
				copyFormulaList := utils.CopySlice(keyFormulaList)
				copyFormulaList = append(copyFormulaList, middleConditionResultDetail.KeyFormulaList...)
				//middleConditionResultDetail.KeyFormulaList = copyFormulaList
				copyCallMap := utils.CopyMap(callMap)
				// 向后添加元素
				offerNode := sourceNode.Offer(middleConditionResultDetail.ConditionNode)
				// 返回出去
				results = append(results, &ConditionNodeResult{
					ConditionNode:  offerNode,
					IsBreak:        middleConditionResultDetail.IsBreak,
					KeyFormulaList: copyFormulaList,
					FormulaCallMap: copyCallMap,
				})
			}
		}
	}
	// 如果没有子节点，那么返回当前节点
	if len(results) == 0 {
		results = append(results, &ConditionNodeResult{
			ConditionNode: cn,
			IsBreak:       false,
		})
	}

	return results, cn
}
