package stmt

import (
	"caseGenerator/common/utils"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression"
	_struct "caseGenerator/parser/struct"
	"encoding/json"
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
	blockResultList := make([]*ConditionNodeResult, 0, 10)
	// 2. 解析 if下面 遍历Block, 处理block下面多个condition的关系，1对多的关系
	ifResultList := i.ParseIfBlockCondition(parentNode)
	if len(ifResultList) > 0 {
		blockResultList = append(blockResultList, ifResultList...)
	}
	// 3. 再遍历else if的block, 首先要排除第二步的condition
	elseIfResultList, uncleNode := i.ParseElseIfCondition(parentNode)
	if len(elseIfResultList) > 0 {
		blockResultList = append(blockResultList, elseIfResultList...)
	}
	// 4. 再遍历else的block
	elseResultList := i.ParseElseCondition(uncleNode)
	if len(elseResultList) > 0 {
		blockResultList = append(blockResultList, elseResultList...)
	}
	// 5. 如果if对应的结果列表全部都是 break 结尾，那么需要再给自己取个反，否则就阻塞了
	if len(elseIfResultList) == 0 && len(elseResultList) == 0 {
		isBreak := true
		for _, v := range ifResultList {
			if !v.IsBreak {
				isBreak = false
				break
			}
		}
		// 如果是break的，那么取[0]进行取反额外插入
		if isBreak {
			negateResult := ifResultList[0]
			// deep-copy
			copyNegateResult := DeepCopyErrorByJson(negateResult)
			// 只需要处理 最外层的条件取反
			copyNegateResult.ConditionNode.ConditionResult = !copyNegateResult.ConditionNode.ConditionResult
			copyNegateResult.IsBreak = false
			copyNegateResult.ConditionNode.Relation = nil

			blockResultList = append(blockResultList, copyNegateResult)
		}
	}

	return blockResultList
}

// ParseIfBlockCondition 解析block下的逻辑，
// 举例说明：
//
//		if b>0 {
//			if a == 5 {
//				fmt.Println("a ==5 && b > 0")
//	     	} else if a < 5 {
//				fmt.Println("a < 5 && b > 0")
//			} else if a > 10 {
//				fmt.Println("a > 5 && b > 0")
//				return
//			}
//			if b == 10 {
//				fmt.Println("b == 10")
//				return
//			}
//			if b < 10 {
//				fmt.Println("b < 10")
//			}
//		}
//
// parentNode
//
//	第一次: b>0
//	第二次: a == 5
//	第三次: b == 10
func (i *If) ParseIfBlockCondition(parentNode *ConditionNodeResult) []*ConditionNodeResult {
	// 1. 处理多个block下面condition的关系
	// resultList 形成的case列表
	// resultList 第一次: 1. a == 5 ;    (第一次处理同个语句中的 if、elseif、else)
	//					 2. a != 5 && a < 5 ;
	//					 3. a != 5 && !(a < 5) && a>10 (break)
	//			  第二次: 1. b == 10  (break)   (和第一次同级)
	//			  		 2. b != 10
	//			  第三次: 1. b < 10 (和第一次同级)
	//					 2. b >= 10
	//			  第三次: 1. a == 5 && b == 10;   是第一、第二次的父级 (将第一步和第二步的数据排列组合出来，分别正向和反向排列，遇到break的不需要取反了)
	//				     2. a == 5 && b != 10 && b < 10;
	//		      		 3. a != 5 && a < 5 && b == 10;
	//		      		 4. a != 5 && a < 5 && b != 10 && b < 10;
	//		      		 5. a != 5 && !(a < 5) && a > 10;
	//			  第四次: 1. b > 0 && a == 5 && b == 10 ;  是第一、第二次的父级 ，(处理本身的 if、 elseif、else)
	//				     2. b > 0 && a == 5 && b != 10 && b < 10 ;
	//		      		 3. b > 0 && a != 5 && a < 5 && b == 10 ;
	//		      		 4. b > 0 && a < 5 && b != 10 && b < 10;
	//		      		 5. b > 0 && a != 5 && !(a < 5) && a > 10

	// 二维数组，用来记录平级的多条语句
	abreastMatrixList := make([][]*ConditionNodeResult, 0, 10)

	// 记录赋值键值对列表
	keyFormulaList := make([]bo.KeyFormula, 0, 10)
	formulaCallMap := make(map[string]*expr.Call, 10)

	// nodeList 取最近的那个条件
	//leftNodeList := make([]*ConditionNode, 0, 10)
	//leftNodeList = append(leftNodeList, parentNode.ConditionNode)
	// 2. for循环解析block里的每条语句
	for _, stmtValue := range i.Block.StmtList {
		// 3. 解析每一句是否是赋值语句，并将其存储起来
		formulaList, callMap := ParseConditionKeyFormula(stmtValue)
		keyFormulaList = append(keyFormulaList, formulaList...)
		utils.PutAll(formulaCallMap, callMap)
		// 4. 解析每一条条件语句，得到的是这句statement解析出来的条件列表。
		conditionResultList := ParseCondition(stmtValue)
		// 5. 将解析得到的结果放入到二维数组中, 合并每个条件语句前的赋值语句
		if len(conditionResultList) > 0 {
			for _, v := range conditionResultList {
				copyFormulaList := utils.CopySlice(keyFormulaList)
				copyFormulaList = append(copyFormulaList, v.KeyFormulaList...)
				v.KeyFormulaList = copyFormulaList
			}
			abreastMatrixList = append(abreastMatrixList, conditionResultList)
		}
	}
	// 6. 解析二维数组，将第一步和第二步的数据排列组合出来，分别正向和反向排列，遇到break的不需要取反了
	// 存放最终结果
	var resultList []*ConditionNodeResult

	// 调用递归函数, 得到 resultList
	if len(abreastMatrixList) > 0 {
		// 临时存放当前组合
		current := make([]*ConditionNodeResult, len(abreastMatrixList))
		generateCombinations(abreastMatrixList, 0, current, &resultList)
	}

	// 如果没有任何的条件语句，那么把主条件放到结果中
	// 含义是: IF xxx {} 下没有任何其他的条件，那么把 xxx 这个条件放到结果中返回
	if len(resultList) == 0 {
		resultList = append(resultList, parentNode)
	} else {
		parentResultList := make([]*ConditionNodeResult, 0, 10)
		for _, v := range resultList {
			copyNode := DeepCopyErrorByJson(parentNode)
			resultNode := copyNode.ConditionNode.Offer(v.ConditionNode)
			parentResultList = append(parentResultList, &ConditionNodeResult{
				ConditionNode:  resultNode,
				IsBreak:        v.IsBreak,
				KeyFormulaList: v.KeyFormulaList,
				FormulaCallMap: formulaCallMap,
			})
		}
		return parentResultList
	}
	return resultList
}

// generateCombinations 递归生成二维数组的排列组合，遇到break的不需要继续处理了
func generateCombinations(arr [][]*ConditionNodeResult, index int, current []*ConditionNodeResult, resultList *[]*ConditionNodeResult) {
	if index == len(arr) {
		// 当递归到达数组末尾，将组合结果添加到结果集
		// 遍历current， 组装 offer，遇到break 就跳出处理
		var result *ConditionNodeResult
		for _, v := range current {
			// 深度拷贝
			copyNode := DeepCopyErrorByJson(v)

			if result == nil {
				result = copyNode
			} else {
				node := result.ConditionNode.Offer(copyNode.ConditionNode)
				result.ConditionNode = node
				result.IsBreak = copyNode.IsBreak
			}
			// 如果已经是 break，那么就不需要 offer 后面的记录
			if copyNode.IsBreak {
				break
			}
		}

		if result != nil {
			// append 前要先判断是否已经存在了
			conflictTag := isResultConflict(*resultList, result)
			if !conflictTag {
				*resultList = append(*resultList, result)
			}
		}
		return
	}

	// 遍历当前子数组
	for _, value := range arr[index] {
		current[index] = value
		generateCombinations(arr, index+1, current, resultList)
	}
}

func isResultConflict(targetList []*ConditionNodeResult, source *ConditionNodeResult) bool {
	// 如果source只有isBreak字段，那么跳出
	if source.ConditionNode == nil {
		return false
	}
	for _, target := range targetList {
		// 如果次条target只有 isBreak字段，那么 continue
		if target.ConditionNode == nil {
			continue
		}
		// 如果source和 target的ConditionResult不一致，跳过
		if source.ConditionNode.ConditionResult != target.ConditionNode.ConditionResult {
			continue
		}
		// 如果 source和 target的length不一样长，那么跳过
		if len(source.ConditionNode.Condition) != len(target.ConditionNode.Condition) {
			continue
		}

		// 遍历的比较 condition每个值， 默认是一致的
		sourceExpr := source.ConditionNode.ExprString()
		targetExpr := target.ConditionNode.ExprString()

		// 如果每一条 expr都一致，那么返回 true，代表了有冲突
		if sourceExpr == targetExpr {
			return true
		}
	}

	// 默认是不冲突的
	return false
}

// todo 怎么给方法加上注释

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

func DeepCopyErrorByJson[T any](p T) T {
	data, err := json.Marshal(p)
	if err != nil {
		panic(err.Error())
	}
	var newP T
	err = json.Unmarshal(data, &newP)
	if err != nil {
		panic(err.Error())
	}
	return newP
}
