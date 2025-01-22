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
	elseIfConditionResultList, uncleNode := i.ParseElseIfCondition(parentNode)
	if len(elseIfConditionResultList) > 0 {
		blockResultList = append(blockResultList, elseIfConditionResultList...)
	}
	// 4. 再遍历else的block
	elseConditionResultList := i.ParseElseCondition(uncleNode)
	if len(elseConditionResultList) > 0 {
		blockResultList = append(blockResultList, elseConditionResultList...)
	}
	return blockResultList
}

// ParseIfBlockCondition 解析block下的多个逻辑，
func (i *If) ParseIfBlockCondition(parentNode *ConditionNodeResult) []*ConditionNodeResult {
	// 1. 处理多个block下面condition的关系
	// resultList 形成的case列表
	resultList := make([]*ConditionNodeResult, 0, 10)
	// nodeList 取最近的那个条件
	nodeList := make([]*ConditionNode, 0, 10)
	nodeList = append(nodeList, parentNode.ConditionNode)
	// 2. for循环解析block里的每条语句
	for _, stmtValue := range i.Block.StmtList {
		// 2.1 解析每一条语句，得到的是这句statement解析出来的条件列表。
		middleNodeList := make([]*ConditionNode, 0, 10)

		conditionResultList := ParseCondition(stmtValue)
		// 2.2 如果是conditionAble的语句会返回结果
		if len(conditionResultList) > 0 {
			// 3. 遍历将结果作为子节点挂在父 condition下，
			for _, conditionResult := range conditionResultList {
				for _, node := range nodeList {
					// 遍历将子节点挂在父condition下
					resultNode, nodeResultList := parseConditionResult(conditionResult, node, parentNode)
					parentSourceNode, err := utils.DeepCopyByJson(resultNode)
					if err != nil {
						panic(err.Error())
					}
					node = parentSourceNode.Offer(node)
					middleNodeList = append(middleNodeList, node)
					resultList = append(resultList, nodeResultList...)
				}

			}
			nodeList = middleNodeList
		}
	}
	if len(nodeList) != 0 {
		for _, node := range nodeList {
			resultList = append(resultList, &ConditionNodeResult{
				ConditionNode: node,
				IsBreak:       false,
			})
		}
	}

	// 如果没有任何的条件语句，那么把主条件放到结果中
	// 含义是: IF xxx {} 下没有任何其他的条件，那么把 xxx 这个条件放到结果中返回
	if len(resultList) == 0 {
		resultList = append(resultList, parentNode)
	}
	return resultList
}

// parseConditionResult 解析条件语句结果，解析
// todo 怎么给方法加上注释
func parseConditionResult(conditionResult *ConditionNodeResult, sourceNode *ConditionNode, parentNode *ConditionNodeResult) (*ConditionNode, []*ConditionNodeResult) {
	var resultNode *ConditionNode
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	// 如果 sourceNodeList 列表是空
	if sourceNode == nil {
		// 执行parentNode对应的解析
		middleSourceNode, blockResult := doParseBlockParentNode(parentNode, conditionResult)
		// 解析后的middleSourceNode，要放入到一起，下个循环使用其作为参数，继续处理。
		if middleSourceNode != nil {
			resultNode = middleSourceNode
		}
		// 解析后的blockResult，解析到了return，下个循环不能再使用其作为参数。
		if blockResult != nil {
			blockResultList = append(blockResultList, blockResult)
		}
	} else {
		middleSource, nodeResultList := doParseSourceNodeList(conditionResult, sourceNode)
		if middleSource != nil {
			resultNode = middleSource
		}
		if len(nodeResultList) != 0 {
			blockResultList = append(blockResultList, nodeResultList...)
		}
	}
	return resultNode, blockResultList
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

func doParseSourceNodeList(conditionResult *ConditionNodeResult, sourceNode *ConditionNode) (*ConditionNode, []*ConditionNodeResult) {
	middleNodeResultList := make([]*ConditionNodeResult, 0, 10)
	if conditionResult.IsBreak {
		middleSourceNode, middleNodeList := parseIsBreak(sourceNode, conditionResult)
		middleNodeResultList = append(middleNodeResultList, middleNodeList...)
		return middleSourceNode, middleNodeResultList
		// 如果 解析出来的 nodeResult不是被阻塞的，那么需要继续排列组合的进行执行
	} else {
		// 手动深拷贝
		snNode, snErr := utils.DeepCopyByJson(sourceNode)
		if snErr != nil {
			panic(snErr.Error())
		}
		// 在尾部加上conditionNode
		result := snNode.Offer(conditionResult.ConditionNode)

		return result, nil
	}
}

func parseIsBreak(sn *ConditionNode, conditionResult *ConditionNodeResult) (*ConditionNode, []*ConditionNodeResult) {
	middleNodeResultList := make([]*ConditionNodeResult, 0, 10)
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
		return nil, middleNodeResultList
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
	return negativeResult, middleNodeResultList
}

// ParseElseIfCondition 解析else if 列表， else if 互相有取反的逻辑
func (i *If) ParseElseIfCondition(parentNode *ConditionNodeResult) ([]*ConditionNodeResult, *ConditionNode) {
	// 1. 定义else-if下的条件关系
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	// 定义叔类节点， else前的叔类都需要是false
	var uncleNode = parentNode.ConditionNode

	for _, v := range i.ElseIfCondition {
		// 2. 先将前面的叔类节点取反
		negateNode := negateUncleNode(uncleNode)
		// 3. 解析本condition条件
		elseConditionExpressionList := expression.Express(v.Condition)
		elseExpress := &ConditionNode{
			Condition:       elseConditionExpressionList,
			ConditionResult: true,
			Relation:        nil,
		}
		// 4. 往末端添加元素
		elseNodeResultNode := negateNode.Offer(elseExpress)
		elseNodeResult := &ConditionNodeResult{
			ConditionNode: elseNodeResultNode,
			IsBreak:       false,
		}
		// 5. 把添加完元素的节点
		uncleNode = elseNodeResultNode
		// 6. 解析本condition的block
		elseBlockConditionList := v.ParseIfBlockCondition(elseNodeResult)
		// 7.把本condition的block放在总block节点下
		blockResultList = append(blockResultList, elseBlockConditionList...)
	}
	// 8. 将叔类节点取反返回
	negateNode := negateUncleNode(uncleNode)
	return blockResultList, negateNode
}

func negateUncleNode(uncleNode *ConditionNode) *ConditionNode {
	// 1. 深度拷贝
	copyNode, err := utils.DeepCopyByJson(uncleNode)
	if err != nil {
		panic(err.Error())
	}
	// 2. 遍历将所有节点的状态取 false
	copyNode.Negate()
	// 3. 返回结果
	return copyNode
}

// ParseElseCondition 解析else节点， 前面所有的if、else-if条件都需要取反
func (i *If) ParseElseCondition(uncleNode *ConditionNode) []*ConditionNodeResult {
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	if i.ElseCondition == nil {
		return blockResultList
	}
	// 1. 解析uncle列表
	//var uncleNode = ParseUncleNodeRelation(uncleNodeList)
	//if uncleNode == nil {
	//	return blockResultList
	//}
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
