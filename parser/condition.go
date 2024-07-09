package parser

import (
	"caseGenerator/common/enum"
	"go/ast"
	"go/token"
)

// ConditionNode 条件语句，if、else if、else、switch、case、default等
// 同时要处理if的嵌套
// 也要处理 if中逻辑与、逻辑或、逻辑非的符合结构
type ConditionNode struct {
	NodeType  string
	Condition ExprNode
	Children  []*ConditionNode
	Init      *InitInfo
	Cond      *CondInfo
	IsElse    bool
}

// ExprNode 表达式节点，包括二元表达式、一元表达式、标识符、字面量等
type ExprNode struct {
	NodeType string // Type of the node, e.g., "binary", "unary", "identifier", "literal"
	Expr     string // The string representation of the expression
	ExprType enum.ConditionExprType
	Left     *ExprNode // Left child node
	Right    *ExprNode // Right child node
}

type InitInfo struct {
	initType  enum.ConditionInitType
	initValue any
}

type CondInfo struct {
	xParam *Param
	yParam *ParamValue
	op     token.Token
}

// 解析条件语句
func (s *SourceInfo) parseCondition(n ast.Node) *ConditionNode {
	if n == nil {
		return nil
	}
	var conditionNode = ConditionNode{}

	switch a := n.(type) {
	case *ast.IfStmt:
		// 主要要解析init、cond、else。要注意if、else的嵌套，也要解析body中的if、else
		// 1. 解析init
		if a.Init != nil {
			conditionNode.Init = s.parseIfInit(a.Init)
		}
		// 2. 解析cond
		if a.Cond != nil {
			conditionNode.Cond = s.parseIfCond(a.Cond)
		}
		// 3. 解析else
		if a.Else != nil {
			conditionNode.IsElse = true
		}
		// 4. 解析body, 继续解析其中的if和switch
		if a.Body != nil {
			list := a.Body.List
			nodes := make([]*ConditionNode, 0, 10)
			for _, v := range list {
				if ifNode, ok := v.(*ast.IfStmt); ok {
					condition := s.parseCondition(ifNode)
					nodes = append(nodes, condition)
				}
			}
			conditionNode.Children = nodes
		}

	case *ast.GenDecl:
	}

	return &conditionNode
}

func (s *SourceInfo) parseIfInit(n ast.Stmt) *InitInfo {
	var ifo InitInfo
	switch a := n.(type) {
	case *ast.AssignStmt:
		assignment := s.ParseAssignment(a)
		ifo.initType = enum.CONDITION_INIT_TYPE_ASSIGNMENT
		ifo.initValue = assignment
	default:
		panic("未找到可用类型")
	}
	return &ifo
}

func (s *SourceInfo) parseIfCond(n ast.Expr) *CondInfo {
	var ifo CondInfo
	switch a := n.(type) {
	case *ast.BinaryExpr:
		ifo.xParam = s.ParamParse(a.X)
		ifo.yParam = s.ParamParseValue(a.Y)
		ifo.op = a.Op
	}
	return &ifo
}
