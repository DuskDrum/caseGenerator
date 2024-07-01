package parser

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// ConditionNode 条件语句，if、else if、else、switch、case、default等
// 同时要处理if的嵌套
// 也要处理 if中逻辑与、逻辑或、逻辑非的符合结构
type ConditionNode struct {
	NodeType  string
	Condition ExprNode
	Children  []*ConditionNode
}

// ExprNode 表达式节点，包括二元表达式、一元表达式、标识符、字面量等
type ExprNode struct {
	NodeType string // Type of the node, e.g., "binary", "unary", "identifier", "literal"
	Expr     string // The string representation of the expression
	ExprType enum.ConditionExprType
	Left     *ExprNode // Left child node
	Right    *ExprNode // Right child node
}

// 解析条件语句
func (s *SourceInfo) parseCondition(n ast.Node) *ConditionNode {
	if n == nil {
		return nil
	}
	conditionNode := &ConditionNode{}

	switch a := n.(type) {
	case *ast.IfStmt:
		// 主要要解析init、cond、else。要注意if、else的嵌套，也要解析body中的if、else
		_ = a.If
	case *ast.GenDecl:
	case *ast.DeclStmt:
	// 这种是没有响应值的function
	case *ast.ExprStmt:

	}

	return conditionNode
}
