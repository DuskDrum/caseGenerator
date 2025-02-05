package decl

import (
	"caseGenerator/common/utils"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/stmt"
	"go/ast"

	"github.com/samber/lo"
)

// Func 解析
type Func struct {
	Receiver []expr.Field  // receiver (methods); or nil (functions)
	Name     expr.Ident    // function/method name
	Type     expr.FuncType // function signature: type and value parameters, results, and position of "func" keyword
	Body     *stmt.Block
}

// ParseFunc 解析Func
func ParseFunc(decl *ast.FuncDecl) *Func {
	fieldList := make([]expr.Field, 0, 10)
	if decl.Recv != nil {
		for _, v := range decl.Recv.List {
			pf := expr.ParseField(v)
			if pf != nil {
				fieldList = append(fieldList, lo.FromPtr(pf))
			}
		}
	}
	name := expr.ParseIdent(decl.Name)
	funcType := expr.ParseFuncType(decl.Type)
	body := stmt.ParseBlock(decl.Body)

	f := &Func{
		Receiver: fieldList,
		Name:     lo.FromPtr(name),
		Type:     lo.FromPtr(funcType),
		Body:     body,
	}
	return f
}

// ParseBody 解析方法
func ParseBody(sb *ast.BlockStmt) {
	// 常量map
	constantsMap := make(map[string]any, 10)
	// 内部变量map， 通过上下文推导出来
	innerVariablesMap := make(map[string]any, 10)
	// 外部变量map，从请求参数里拿到
	outerVariablesMap := make(map[string]any, 10)
	// 公式列表
	keyFormulaList := make([]bo.KeyFormula, 0, 10)
	callVariableMap := make(map[string]*expr.Call, 10)
	nodeResultList := make([]*stmt.ConditionNodeResult, 0, 10)

	// 1. 遍历方法里的每个 stmt，
	for _, v := range sb.List {
		// todo
		// 1.1 解析每个stmt的内容
		p := stmt.ParseStmt(v)
		// 1.2 执行每一个 stmt 的公式；得到最新的赋值公式
		if pes, ok := p.(stmt.ExpressionStmt); ok {
			kfList, callMap := pes.FormulaExpress()
			keyFormulaList = append(keyFormulaList, kfList...)
			for k, value := range callMap {
				callVariableMap[k] = value
			}
		}
		// 1.3 执行每一句的condition
		conditionResultList := stmt.ParseCondition(v)
		// 1.4 拼接每一个 case 的赋值信息和条件信息
		for _, result := range conditionResultList {
			// 深度拷贝赋值语句
			copyFormulaList := utils.CopySlice(keyFormulaList)
			copyCallMap := utils.CopyMap(callVariableMap)
			// new 新的 result
			nodeResultList = append(nodeResultList, &stmt.ConditionNodeResult{
				ConditionNode:  result.ConditionNode,
				IsBreak:        result.IsBreak,
				KeyFormulaList: copyFormulaList,
				FormulaCallMap: copyCallMap,
			})
		}
	}
	// 2. 遍历执行所有公式；
	// 按照先遍历赋值语句，再遍历条件语句的顺序，执行formula
	for _, v := range nodeResultList {
		moker.MockExpression()
	}

	// 3. 找到需要mock的值，开始生成mock结构体

	// 3. 得到mock结果

}
