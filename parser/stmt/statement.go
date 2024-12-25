package stmt

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"go/ast"
)

type Stmt interface {
}

// ExpressionStmt 参数接口类型，将参数需要的方法定义出来
// ConstantsMap map[string]any-->常量map，上游传进来，常量map不会再变
// InnerVariablesMap map[string]any， 不需要解析------>变量值Map (可以通过上下文推导出的变量， 局部变量都是可推导的),
// 解析出OuterVariablesMap，map[string]类型，表明它是外部变量 ------> 变量字段Map(无法通过上下文推导出变量， 方法请求传进来的变量，或者用方法得到的)====>后面需要用算法算出来
// FormulasList --> 公式List(Key 对应那个变量，Formula 对应的公式)
type ExpressionStmt interface {
	// FormulaExpress  生成逻辑表达式, 入参常量Map， 出参 公式List，外部变量map
	FormulaExpress() ([]bo.KeyFormula, map[string]enum.SpecificType)
}

// Condition 条件类型: if、switch、typeSwitch
type Condition interface {
	// CalculateCondition 解析Condition
	CalculateCondition([]bo.StatementAssignment) []ConditionResult
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
func ParseStmt(expr ast.Stmt) Stmt {
	if expr == nil {
		return nil
	}
	switch stmtType := expr.(type) {
	case *ast.DeclStmt:
		return ParseDecl(stmtType)
	case *ast.EmptyStmt:
		return ParseEmpty(stmtType)
	case *ast.LabeledStmt:
		return ParseLabeled(stmtType)
	case *ast.ExprStmt:
		return ParseExpr(stmtType)
	case *ast.SendStmt:
		return ParseSend(stmtType)
	case *ast.IncDecStmt:
		return ParseIncDec(stmtType)
	case *ast.AssignStmt:
		return ParseAssign(stmtType)
	case *ast.GoStmt:
		return ParseGo(stmtType)
	case *ast.DeferStmt:
		return ParseDefer(stmtType)
	case *ast.ReturnStmt:
		return ParseReturn(stmtType)
	case *ast.BranchStmt:
		return ParseBranch(stmtType)
	case *ast.BlockStmt:
		return ParseBlock(stmtType)
	case *ast.IfStmt:
		return ParseIf(stmtType)
	case *ast.CaseClause:
		return ParseCaseClause(stmtType)
	case *ast.SwitchStmt:
		return ParseSwitch(stmtType)
	case *ast.TypeSwitchStmt:
		return ParseTypeSwitch(stmtType)
	case *ast.CommClause:
		return ParseCommClause(stmtType)
	case *ast.SelectStmt:
		return ParseSelect(stmtType)
	case *ast.ForStmt:
		return ParseFor(stmtType)
	case *ast.RangeStmt:
		return ParseRange(stmtType)
	default:
		panic("未知类型...")
	}
}
