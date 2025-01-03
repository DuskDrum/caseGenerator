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
	for _, stmtValue := range i.Block.StmtList {
		// 2. 解析block中的语句condition，如果不是condition类型则返回nil
		conditionResultList := ParseCondition(stmtValue)
		if len(conditionResultList) > 0 {
			// 3. 遍历将结果作为子节点挂在父 condition下
			for _, conditionResult := range conditionResultList {
				// 手动深拷贝
				condiNode := parentNode.ConditionNode
				// 在尾部加上conditionNode
				result := condiNode.Offer(conditionResult.ConditionNode)
				middleNode := &ConditionNodeResult{
					ConditionNode: result,
					IsBreak:       conditionResult.IsBreak,
				}
				blockResultList = append(blockResultList, middleNode)
			}
		}
	}
	// 如果没有任何的条件语句，那么把主条件放到结果中
	// 含义是: IF xxx {} 下没有任何其他的条件，那么把 xxx 这个条件放到结果中返回
	if len(blockResultList) == 0 {
		blockResultList = append(blockResultList, parentNode)
	}
	return blockResultList
}

// ParseElseIfCondition 解析else if 列表， else if 互相有取反的逻辑
func (i *If) ParseElseIfCondition(parentNode *ConditionNodeResult) ([]*ConditionNodeResult, []*ConditionNode) {
	// 1. 定义else-if下的条件关系
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	// 定义叔类节点， else前的叔类都需要是false
	uncleNodeList := make([]*ConditionNode, 0, 10)
	uncleNodeList = append(uncleNodeList, parentNode.ConditionNode)

	for _, v := range i.ElseIfCondition {
		// 2. 先将前面的叔类节点都取反
		var uncleNode *ConditionNode
		for _, value := range uncleNodeList {
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
		// 3. 解析本condition条件
		elseConditionExpressionList := expression.Express(v.Condition)
		elseNodeResult := &ConditionNodeResult{
			ConditionNode: &ConditionNode{
				Condition:       elseConditionExpressionList,
				ConditionResult: true,
				Relation:        uncleNode,
			},
			IsBreak: false,
		}
		// 4. 解析本condition的block
		elseBlockConditionList := v.ParseIfBlockCondition(elseNodeResult)
		// 5. 把此节点放到 叔叔辈
		uncleNodeList = append(uncleNodeList, elseNodeResult.ConditionNode)
		// 6.把本condition的block放在总block节点下
		blockResultList = append(blockResultList, elseBlockConditionList...)
	}
	return blockResultList, uncleNodeList
}

// ParseElseCondition 解析else节点， 前面所有的if、else-if条件都需要取反
func (i *If) ParseElseCondition(uncleNodeList []*ConditionNode) []*ConditionNodeResult {
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	if i.ElseCondition == nil {
		return blockResultList
	}
	// 1. 解析uncle列表
	var uncleNode *ConditionNode
	for _, value := range uncleNodeList {
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
	if uncleNode == nil {
		return blockResultList
	}
	// 1. 解析else的内容
	for _, stmtValue := range i.ElseCondition.StmtList {
		// 2. 解析每一条condition
		// 2. 解析block中的语句condition，如果不是condition类型则返回nil
		conditionResultList := ParseCondition(stmtValue)
		if len(conditionResultList) > 0 {
			// 3. 遍历将结果作为子节点挂在父 condition下
			for _, conditionResult := range conditionResultList {
				// 手动深拷贝
				condiNode := uncleNode
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
	return blockResultList
}
