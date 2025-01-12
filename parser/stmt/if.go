package stmt

import (
	"caseGenerator/common/utils"
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
			elseIfList, blockStmt := ParseElseIf(elseTyp)
			i.ElseIfCondition = elseIfList
			i.ElseCondition = blockStmt

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

func ParseElseIf(stmt *ast.IfStmt) ([]*If, *Block) {
	ifList := make([]*If, 0, 10)
	// 先解析本身，将其加到切片中
	stmtIf := ParseIf(stmt)

	ifList = append(ifList, stmtIf)
	// 递归去找else-if
	if stmt.Else == nil {
		return ifList, nil
	}
	blockType, ok := stmt.Else.(*ast.BlockStmt)
	if ok {
		// 说明不再是else if, 而是else
		elseBlock := ParseBlock(blockType)
		return ifList, elseBlock
	}
	elseStmt, ok := stmt.Else.(*ast.IfStmt)
	if !ok {
		panic("if else type is not *ast.IfStmt")
	}
	// 再递归的解析else
	elseIfList, blockStmt := ParseElseIf(elseStmt)
	ifList = append(ifList, elseIfList...)
	return ifList, blockStmt
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
	// 1. 先拿到 本 IF 的Condition的表达式 (父类if条件)
	conditionExpressionList := expression.Express(i.Condition)
	parentNode := &ConditionNodeResult{
		ConditionNode: &ConditionNode{
			Condition:       conditionExpressionList,
			ConditionResult: true,
		},
		IsBreak: false,
	}
	// 2. 遍历Block, 处理block下面多个condition的关系，1对多的关系
	blockResultList := i.ParseIfBlockCondition(parentNode)
	// 3. 再遍历else if的block, 首先要排除第二步的condition
	elseIfConditionResultList, uncleNodeList := i.ParseElseIfCondition(parentNode)
	if len(elseIfConditionResultList) > 0 {
		blockResultList = append(blockResultList, elseIfConditionResultList...)
	}
	// 4. 再遍历else的block
	elseConditionResultList := i.ParseElseCondition(uncleNodeList)
	if len(elseConditionResultList) > 0 {
		blockResultList = append(blockResultList, elseConditionResultList...)
	}
	return blockResultList
}

// ParseIfBlockCondition 解析block下的多个逻辑，父节点和子节点之间是 逻辑与关系
func (i *If) ParseIfBlockCondition(parentNode *ConditionNodeResult) []*ConditionNodeResult {
	// 1. 遍历Block, 处理多个block下面condition的关系
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	sourceNodeList := make([]*ConditionNode, 0, 10)

	for _, stmtValue := range i.Block.StmtList {
		// 2. 解析block中的语句condition，如果不是condition类型则返回nil
		conditionResultList := ParseCondition(stmtValue)
		if len(conditionResultList) > 0 {
			middleSourceNodeList := make([]*ConditionNode, 0, 10)
			// 3. 遍历将结果作为子节点挂在父 condition下
			for _, conditionResult := range conditionResultList {
				// 遍历将子节点挂在父condition下
				nodeList, nodeResultList := parseConditionResult(conditionResult, sourceNodeList, parentNode)
				middleSourceNodeList = append(middleSourceNodeList, nodeList...)
				blockResultList = append(blockResultList, nodeResultList...)
			}
			// 赋值
			sourceNodeList = middleSourceNodeList
		}
	}
	for _, v := range sourceNodeList {
		blockResultList = append(blockResultList, &ConditionNodeResult{
			ConditionNode: v,
			IsBreak:       false,
		})
	}
	// 如果没有任何的条件语句，那么把主条件放到结果中
	// 含义是: IF xxx {} 下没有任何其他的条件，那么把 xxx 这个条件放到结果中返回
	if len(blockResultList) == 0 {
		blockResultList = append(blockResultList, parentNode)
	}
	return blockResultList
}

func parseConditionResult(conditionResult *ConditionNodeResult, sourceNodeList []*ConditionNode, parentNode *ConditionNodeResult) ([]*ConditionNode, []*ConditionNodeResult) {
	middleSourceNodeList := make([]*ConditionNode, 0, 10)
	blockResultList := make([]*ConditionNodeResult, 0, 10)

	if len(sourceNodeList) == 0 {
		// 执行parentNode对应的解析
		middleSourceNode, blockResult := doParseBlockParentNode(parentNode, conditionResult)
		// 解析后的middleSourceNode，要放入到一起，下个循环使用其作为参数，继续处理。
		if middleSourceNode != nil {
			middleSourceNodeList = append(middleSourceNodeList, middleSourceNode)
		}
		// 解析后的blockResult，解析到了return，下个循环不能再使用其作为参数。
		if blockResult != nil {
			blockResultList = append(blockResultList, blockResult)
		}
	} else {
		middleSourceList, nodeResultList := doParseSourceNodeList(conditionResult, sourceNodeList)
		if len(middleSourceList) != 0 {
			middleSourceNodeList = append(middleSourceNodeList, middleSourceList...)
		}
		if len(nodeResultList) != 0 {
			blockResultList = append(blockResultList, nodeResultList...)
		}
	}
	return middleSourceNodeList, blockResultList

}

func doParseBlockParentNode(parentNode *ConditionNodeResult, conditionResult *ConditionNodeResult) (*ConditionNode, *ConditionNodeResult) {
	var negativeResult *ConditionNode = nil
	// 手动深拷贝
	parentSourceNode, err := utils.DeepCopyByJson(parentNode.ConditionNode)
	if err != nil {
		panic(err.Error())
	}
	// 在尾部加上conditionNode
	result := parentSourceNode.Offer(conditionResult.ConditionNode)

	// negative parent node
	negativeParentSourceNode, err := utils.DeepCopyByJson(parentNode.ConditionNode)
	if err != nil {
		panic(err.Error())
	}
	// 在尾部加上conditionNode
	if conditionResult.ConditionNode != nil {
		negativeResult = negativeParentSourceNode.Offer(&ConditionNode{
			Condition:       conditionResult.ConditionNode.Condition,
			ConditionResult: !conditionResult.ConditionNode.ConditionResult,
			Relation:        conditionResult.ConditionNode.Relation,
		})
	}

	// 如果 解析出来的 nodeResult是被阻塞的，那么直接返回到结果中， 不放到middleNode中了
	if conditionResult.IsBreak {
		return negativeResult, &ConditionNodeResult{
			ConditionNode: result,
			IsBreak:       true,
		}
		// 如果 解析出来的 nodeResult不是被阻塞的，那么需要继续排列组合的进行执行
	} else {
		return result, nil
	}
}

func doParseSourceNodeList(conditionResult *ConditionNodeResult, sourceNodeList []*ConditionNode) ([]*ConditionNode, []*ConditionNodeResult) {
	middleSourceNodeList := make([]*ConditionNode, 0, 10)
	middleNodeResultList := make([]*ConditionNodeResult, 0, 10)
	if conditionResult.IsBreak {
		for _, sn := range sourceNodeList {
			// 手动深拷贝， 直接return出去了
			snNode, snErr := utils.DeepCopyByJson(sn)
			if snErr != nil {
				panic(snErr.Error())
			}
			// 在尾部加上conditionNode
			result := snNode.Offer(conditionResult.ConditionNode)
			middleNodeResultList = append(middleNodeResultList, &ConditionNodeResult{
				ConditionNode: result,
				IsBreak:       true,
			})

			if conditionResult.ConditionNode == nil {
				continue
			}
			// negative parent node， 同级别的条件要越过这个condition，所以需要取反
			negativeParentSourceNode, err := utils.DeepCopyByJson(sn)
			if err != nil {
				panic(err.Error())
			}
			// 在尾部加上conditionNode
			negativeResult := negativeParentSourceNode.Offer(&ConditionNode{
				Condition:       conditionResult.ConditionNode.Condition,
				ConditionResult: !conditionResult.ConditionNode.ConditionResult,
				Relation:        conditionResult.ConditionNode.Relation,
			})
			middleSourceNodeList = append(middleSourceNodeList, negativeResult)

		}
		return nil, middleNodeResultList
		// 如果 解析出来的 nodeResult不是被阻塞的，那么需要继续排列组合的进行执行
	} else {
		for _, sn := range sourceNodeList {
			// 手动深拷贝
			snNode, snErr := utils.DeepCopyByJson(sn)
			if snErr != nil {
				panic(snErr.Error())
			}
			// 在尾部加上conditionNode
			result := snNode.Offer(conditionResult.ConditionNode)
			middleSourceNodeList = append(middleSourceNodeList, result)
		}
		return middleSourceNodeList, nil
	}
}

// ParseElseIfCondition 解析else if 列表， else if 互相有取反的逻辑
func (i *If) ParseElseIfCondition(parentNode *ConditionNodeResult) ([]*ConditionNodeResult, []*ConditionNodeResult) {
	// 1. 定义else-if下的条件关系
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	// 定义叔类节点， else前的叔类都需要是false
	uncleNodeList := make([]*ConditionNodeResult, 0, 10)
	uncleNodeList = append(uncleNodeList, parentNode)

	for _, v := range i.ElseIfCondition {
		// 2. 先将前面的叔类节点都取反
		var uncleNode = ParseUncleNodeRelation(uncleNodeList)
		// 3. 解析本condition条件
		elseConditionExpressionList := expression.Express(v.Condition)
		elseNodeResultNode := uncleNode.Offer(&ConditionNode{
			Condition:       elseConditionExpressionList,
			ConditionResult: true,
			Relation:        nil,
		})
		elseNodeResult := &ConditionNodeResult{
			ConditionNode: elseNodeResultNode,
			IsBreak:       false,
		}
		// 4. 解析本condition的block
		elseBlockConditionList := v.ParseIfBlockCondition(elseNodeResult)
		// 5. 把此节点放到 叔叔辈
		uncleNodeList = append(uncleNodeList, elseNodeResult)
		// 6.把本condition的block放在总block节点下
		blockResultList = append(blockResultList, elseBlockConditionList...)
	}
	return blockResultList, uncleNodeList
}

// ParseElseCondition 解析else节点， 前面所有的if、else-if条件都需要取反
func (i *If) ParseElseCondition(uncleNodeList []*ConditionNodeResult) []*ConditionNodeResult {
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	if i.ElseCondition == nil {
		return blockResultList
	}
	// 1. 解析uncle列表
	var uncleNode = ParseUncleNodeRelation(uncleNodeList)
	if uncleNode == nil {
		return blockResultList
	}
	// 2. 解析else的内容
	for _, stmtValue := range i.ElseCondition.StmtList {
		// 2. 解析每一条condition
		// 2. 解析block中的语句condition，如果不是condition类型则返回nil
		conditionResultList := ParseCondition(stmtValue)
		if len(conditionResultList) > 0 {
			// 3. 遍历将结果作为子节点挂在父 condition下
			for _, conditionResult := range conditionResultList {
				// 手动深拷贝
				condiNode, snErr := utils.DeepCopyByJson(uncleNode)
				if snErr != nil {
					panic(snErr.Error())
				}
				// 在尾部加上conditionNode
				condiNode.Offer(conditionResult.ConditionNode)
				middleNode := &ConditionNodeResult{
					ConditionNode: condiNode,
					IsBreak:       conditionResult.IsBreak,
				}
				blockResultList = append(blockResultList, middleNode)
			}
		}
	}
	if len(blockResultList) == 0 {
		blockResultList = append(blockResultList, &ConditionNodeResult{
			ConditionNode: uncleNode,
			IsBreak:       false,
		})
	}
	return blockResultList
}

func ParseUncleNodeRelation(uncleNodeList []*ConditionNodeResult) *ConditionNode {
	var uncleNode *ConditionNode
	for _, nodeResult := range uncleNodeList {
		if nodeResult.IsBreak {
			continue
		}
		value := nodeResult.ConditionNode
		result := &ConditionNode{
			Condition:       value.Condition,
			ConditionResult: false, // 代表这个条件要取反
		}
		if uncleNode == nil {
			uncleNode = result
		} else {
			resultNode := uncleNode.Offer(result)
			uncleNode = resultNode
		}
	}
	return uncleNode
}
