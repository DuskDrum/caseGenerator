package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"caseGenerator/parser/expression/govaluate"
	"go/ast"
	"go/token"
	"strings"
)

type Stmt interface {
}

// ExpressionStmt 参数接口类型，将参数需要的方法定义出来
// ConstantsMap map[string]any-->常量map，上游传进来，常量map不会再变
// InnerVariablesMap map[string]any， 不需要解析------>变量值Map (可以通过上下文推导出的变量， 局部变量都是可推导的),
// OuterVariablesMap，map[string]类型，表明它是外部变量 ------> 变量字段Map(无法通过上下文推导出变量， 方法请求传进来的变量)====>后面需要用算法算出来
// CallMap,map[string]*ast.Call,
// FormulasList --> 公式List(Key 对应那个变量，Formula 对应的公式)
type ExpressionStmt interface {
	// FormulaExpress  生成逻辑表达式, 入参常量Map， 出参 公式List，方法调用map
	FormulaExpress() ([]govaluate.KeyFormula, map[string]*expr.Call)
}

// Condition 条件类型: if、switch、typeSwitch、return
// 需要理清楚每个分支要走到哪些逻辑公式，和对应的Condition与非condition
// if 里要考虑 else、嵌套if、嵌套switch、嵌套 type-switch之间的关系，也要考虑 return 直接跳出 condition
// 根据数据变化 + condition得到最终需要得到的对应外部变量、call变量的值
// 首先要理出来condition的树状结构， 其中包含了if、else-if、else、switch-case、switch-default、return之间的关系
// todo 考虑fori
type Condition interface {
	// CalculateCondition 解析Condition
	CalculateCondition(constantsMap, innerVariablesMap, outerVariablesMap map[string]any, keyFormulaList []govaluate.KeyFormula) []ConditionResult
}

type ConditionResult struct {
	IdentMap    map[string]IdentConditionResult
	CallMap     map[string]CallConditionResult
	SelectorMap map[string]SelectorConditionResult
}

type IdentConditionResult struct {
}
type CallConditionResult struct {
}
type SelectorConditionResult struct {
}

// ParseStmt 完整的执行单元
func ParseStmt(expr ast.Stmt, context bo.ExprContext) Stmt {
	if expr == nil {
		return nil
	}
	switch stmtType := expr.(type) {
	case *ast.DeclStmt:
		return ParseDecl(stmtType, context)
	case *ast.EmptyStmt:
		return ParseEmpty(stmtType, context)
	case *ast.LabeledStmt:
		return ParseLabeled(stmtType, context)
	case *ast.ExprStmt:
		return ParseExpr(stmtType, context)
	case *ast.SendStmt:
		return ParseSend(stmtType, context)
	case *ast.IncDecStmt:
		return ParseIncDec(stmtType, context)
	case *ast.AssignStmt:
		return ParseAssign(stmtType, context)
	case *ast.GoStmt:
		return ParseGo(stmtType, context)
	case *ast.DeferStmt:
		return ParseDefer(stmtType, context)
	case *ast.ReturnStmt:
		return ParseReturn(stmtType, context)
	case *ast.BranchStmt:
		return ParseBranch(stmtType, context)
	case *ast.BlockStmt:
		return ParseBlock(stmtType, context)
	case *ast.IfStmt:
		return ParseIf(stmtType, context)
	case *ast.CaseClause:
		return ParseCaseClause(stmtType, context)
	case *ast.SwitchStmt:
		return ParseSwitch(stmtType, context)
	case *ast.TypeSwitchStmt:
		return ParseTypeSwitch(stmtType, context)
	case *ast.CommClause:
		return ParseCommClause(stmtType, context)
	case *ast.SelectStmt:
		return ParseSelect(stmtType, context)
	case *ast.ForStmt:
		return ParseFor(stmtType, context)
	case *ast.RangeStmt:
		return ParseRange(stmtType, context)
	default:
		panic("未知类型...")
	}
}

type ConditionNode struct {
	Condition       []*govaluate.ExpressDetail // 表示当前节点的条件 (如 "A", "B", "C", "D", "E")
	ConditionResult bool                       // 表示逻辑与还是逻辑非，true代表这一个condition要是true， false代表这个condition要是false
	Relation        *ConditionNode             // 关联节点表示嵌套逻辑
	Position        token.Position
}

type ConditionNodeExpress struct {
	Condition       []*govaluate.ExpressDetail // 表示当前节点的条件 (如 "A", "B", "C", "D", "E")
	ConditionResult bool                       // 表示逻辑与还是逻辑非，true代表这一个condition要是true， false代表这个condition要是false
	Position        token.Position
}

// Offer 往线性结构末端继续添加元素
func (cn *ConditionNode) Offer(relation *ConditionNode) *ConditionNode {
	if relation == nil {
		result := cn
		return result
	}
	if cn.Relation == nil {
		result := cn
		result.Relation = relation
		return result
	} else {
		result := cn.Relation.Offer(relation)
		cn.Relation = result
		return cn
	}
}

// Negate 取反
func (cn *ConditionNode) Negate() {
	cn.ConditionResult = false
	if cn.Relation == nil {
		return
	} else {
		cn.Relation.Negate()
	}
}

// ExprString 取所有 relation 的 expr 值
func (cn *ConditionNode) ExprString() string {
	exprList := make([]string, 0, 10)
	for _, c := range cn.Condition {
		if cn.ConditionResult {
			exprList = append(exprList, c.Expr)
		} else {
			exprList = append(exprList, "!("+c.Expr+")")
		}
	}
	joinString := strings.Join(exprList, " && ")

	if cn.Relation == nil {
		return joinString
	} else {
		return joinString + " && " + cn.Relation.ExprString()
	}
}

// Add 往线性结构头部添加元素
func (cn *ConditionNode) Add(parent *ConditionNode) *ConditionNode {
	node := &ConditionNode{
		Condition:       parent.Condition,
		ConditionResult: parent.ConditionResult,
		Relation:        cn,
	}
	return node
}

type ConditionNodeResult struct {
	ConditionNode  *ConditionNode         // 条件节点，有子条件，象征着一条条件链路
	IsBreak        bool                   // 表示是否已经中断，true代表已经中断了，不需要继续处理条件， false代表这个condition还需要继续处理
	KeyFormulaList []govaluate.KeyFormula // 赋值键值对列表
	FormulaCallMap map[string]*expr.Call  // 赋值方法map
}

// ParseCondition 两个响应值， 第一个是多条条件节点，第二个是是否中断
func ParseCondition(condition Stmt) []*ConditionNodeResult {
	switch conditionType := condition.(type) {
	case *If:
		return conditionType.ParseIfCondition()
	case *Switch:
		return conditionType.ParseSwitchCondition()
	case *TypeSwitch:
		return conditionType.ParseTypeSwitchCondition()
	case *Return:
		return conditionType.ParseReturnCondition()
	}
	return nil
}

// ParseConditionKeyFormula 两个响应值， 第一个是多条KeyFormula，第二个是对应要处理的 callMap
func ParseConditionKeyFormula(condition Stmt) ([]govaluate.KeyFormula, map[string]*expr.Call) {
	// 1. 解析每个stmt的内容；得到最新的赋值公式
	// 类型断言判断
	if pes, ok := condition.(ExpressionStmt); ok {
		return pes.FormulaExpress()
	}
	return []govaluate.KeyFormula{}, map[string]*expr.Call{}
}
