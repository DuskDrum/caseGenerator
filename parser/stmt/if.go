package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// If 条件语句
// 代码块是由花括号{}包围的一系列语句，它在函数体、控制结构（如if、for、switch）等场景中广泛使用。
type If struct {
	Init            *Assign
	Condition       _struct.Parameter
	Block           *Block
	ElseIfCondition []*If  // else-if列表
	ElseCondition   *Block // else记录
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
	i.Block = ParseBlock(stmt.Body)
	i.Condition = expr.ParseParameter(stmt.Cond)

	if stmt.Else != nil {
		switch elseTyp := stmt.Else.(type) {
		case *ast.IfStmt:
			elseIfList := ParseElseIf(elseTyp)
			i.ElseIfCondition = elseIfList
		case *ast.BlockStmt:
			// 解析else
			elseBlock := ParseBlock(elseTyp)
			i.ElseCondition = elseBlock
		default:
			panic("parse condition error, can't get type")
		}

	}
	return i
}

func ParseElseIf(stmt *ast.IfStmt) []*If {
	ifList := make([]*If, 0, 10)
	// 递归去找else-if
	if stmt.Else == nil {
		return ifList
	}
	ifStmt, ok := stmt.Else.(*ast.IfStmt)
	if !ok {
		panic("if else type is not *ast.IfStmt")
	}
	elseIfList := ParseElseIf(ifStmt)
	ifList = append(ifList, elseIfList...)
	return ifList
}

func (i *If) FormulaExpress() ([]bo.KeyFormula, map[string]*expr.Call) {
	// switch 和 if 是分成两部分的， 1. init 部分组装 expression；2.整个公式来计算得到需要 mocker 的值
	if i.Init != nil {
		return i.Init.FormulaExpress()
	}
	return nil, nil
}

func (i *If) CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []bo.KeyFormula) []ConditionResult {
	// 1. 先拿到 Condition的表达式
	//conditionExpressionList := expression.Express(i.Condition)
	// 2. 找表达式中的变量,去遍历找表达式中的变化记录
	//for _, v := range conditionExpressionList {
	//	mockList := mocker.MockExpression(v, seList)
	//	if len(mockList) > 0 {
	//		fmt.Printf("mock结果列表: %v\n", mockList)
	//	}
	//}

	return nil
}

func (i *If) ParseIfCondition() []*ConditionNodeResult {
	// 1. 先拿到 Condition的表达式
	conditionExpressionList := expression.Express(i.Condition)
	parentNode := &ConditionNodeResult{
		ConditionNode: &ConditionNode{
			Condition:       conditionExpressionList,
			ConditionResult: true,
		},
		IsBreak: false,
	}

	// 2. 遍历Block, 处理多个block下面condition的关系
	// 最外层的condition, 和父类同级的条件列表
	uncleNodeList := make([]*ConditionNodeResult, 0, 10)
	uncleNodeList = append(uncleNodeList, parentNode)

	blockResultList := make([]*ConditionNodeResult, 0, 10)
	blockResultList = append(blockResultList, parentNode)
	for _, stmtValue := range i.Block.StmtList {
		// 2.1 解析condition
		conditionResultList := ParseCondition(stmtValue)
		if len(conditionResultList) > 0 {
			blockMiddleResultList := make([]*ConditionNodeResult, 0, 10)
			for _, result := range blockResultList {
				for _, conditionResult := range conditionResultList {
					condiNode := result.ConditionNode
					condiNode.Offer(conditionResult.ConditionNode)
					middleNode := &ConditionNodeResult{
						ConditionNode: condiNode,
						IsBreak:       conditionResult.IsBreak,
					}
					blockMiddleResultList = append(blockMiddleResultList, middleNode)
				}
			}
		}
	}
	// 3. 再遍历else if的block, 首先要排除第二步的condition
	for _, v := range i.ElseIfCondition {
		elseConditionExpressionList := expression.Express(v.Condition)
		elseNodResult := &ConditionNodeResult{
			ConditionNode: &ConditionNode{
				Condition:       elseConditionExpressionList,
				ConditionResult: true,
			},
			IsBreak: false,
		}
		// 递归调用
		conditionResultList := v.ParseIfCondition()
		elseResultNodeList := make([]*ConditionNodeResult, 0, 10)

		// 遍历把前面几个original都取逻辑非放在最前面节点
		var uncleNode *ConditionNode
		for _, value := range uncleNodeList {
			result := &ConditionNode{
				Condition:       value.ConditionNode.Condition,
				ConditionResult: false, // 代表这个条件要取反
			}
			if uncleNode == nil {
				uncleNode = result
			} else {
				uncleNode.Offer(result)
			}
		}

		for _, elseCondition := range conditionResultList {
			if uncleNode != nil {
				addNode := uncleNode.Add(elseCondition.ConditionNode)
				addResult := &ConditionNodeResult{
					ConditionNode: addNode,
					IsBreak:       elseCondition.IsBreak,
				}
				elseResultNodeList = append(elseResultNodeList, addResult)
			} else {
				panic("uncleNode error")
			}
		}
		// 把此节点放到 叔叔辈
		uncleNodeList = append(uncleNodeList, elseNodResult)
		// 把解析好的节点 放在block节点下
		blockResultList = append(blockResultList, elseResultNodeList...)
	}
	// 4. 再遍历else的block
	for _, stmtValue := range i.ElseCondition.StmtList {
		// 4.1 解析condition
		conditionResultList := ParseCondition(stmtValue)
		// 遍历把前面几个original都取逻辑非放在最前面节点
		var uncleNode *ConditionNode
		for _, value := range uncleNodeList {
			result := &ConditionNode{
				Condition:       value.ConditionNode.Condition,
				ConditionResult: false, // 代表这个条件要取反
			}
			if uncleNode == nil {
				uncleNode = result
			} else {
				uncleNode.Offer(result)
			}
		}
		if uncleNode != nil && len(conditionResultList) > 0 {
			for _, result := range conditionResultList {
				addNode := uncleNode.Add(result.ConditionNode)
				nodeResult := &ConditionNodeResult{
					ConditionNode: addNode,
					IsBreak:       result.IsBreak,
				}
				// 把解析好的节点 放在block节点下
				blockResultList = append(blockResultList, nodeResult)
			}
		}
	}

	return blockResultList
}
