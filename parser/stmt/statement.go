package stmt

import (
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// ParseStmt 完整的执行单元
func ParseStmt(expr ast.Stmt) _struct.Stmt {
	switch stmtType := expr.(type) {
	case *ast.DeclStmt:
		ParseDecl(stmtType)
	case *ast.EmptyStmt:
		ParseEmpty(stmtType)
	case *ast.LabeledStmt:
		ParseLabeled(stmtType)
	case *ast.ExprStmt:
		ParseExpr(stmtType)
	case *ast.SendStmt:
		ParseSend(stmtType)
	case *ast.IncDecStmt:
		ParseIncDec(stmtType)
	case *ast.AssignStmt:
		ParseAssign(stmtType)
	case *ast.GoStmt:
		ParseGo(stmtType)
	case *ast.DeferStmt:
		ParseDefer(stmtType)
	case *ast.ReturnStmt:
		ParseReturn(stmtType)
	case *ast.BranchStmt:
		ParseBranch(stmtType)
	case *ast.BlockStmt:
		ParseBlock(stmtType)
	case *ast.IfStmt:
		ParseIf(stmtType)
	case *ast.CaseClause:
		ParseCaseClause(stmtType)
	case *ast.SwitchStmt:
		ParseSwitch(stmtType)
	case *ast.TypeSwitchStmt:
		ParseTypeSwitch(stmtType)
	case *ast.CommClause:
		ParseCommClause(stmtType)
	case *ast.SelectStmt:
		ParseSelect(stmtType)
	case *ast.ForStmt:
		ParseFor(stmtType)
	case *ast.RangeStmt:
		ParseRange(stmtType)
	default:
		panic("未知类型...")
	}
	return nil
}
