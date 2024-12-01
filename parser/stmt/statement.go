package stmt

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/expression"
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// Stmt 参数接口类型，将参数需要的方法定义出来
type Stmt interface {
	// Express 生成表达式.
	Express() []StatementExpression
}

// Condition 条件类型: if、switch、typeSwitch
type Condition interface {
	CalculateCondition([]StatementExpression) []ConditionResult
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

// StatementExpression stmt的表达式，记录了参数的变动, 参数也可以直接重新赋值
type StatementExpression struct {
	Name      string
	InitParam _struct.ValueAble
	Type      enum.StmtType
	// 参数变动列表
	expression.Expression
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
