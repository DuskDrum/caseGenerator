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
	Children []*ConditionNode
	Init     *InitInfo
	Cond     *CondInfo
	IsElse   bool
	Else     *ConditionNode
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
	InitType  enum.ConditionInitType
	InitValue any
}

type CondInfo struct {
	XParam *Param
	YParam *ParamValue
	Op     token.Token
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
		conditionNode.IsElse = false

		// 3. 解析else
		if a.Else != nil {
			var conditionElseNode = ConditionNode{}
			conditionElseNode.IsElse = true
			conditionElseNode.Cond = conditionNode.Cond
			conditionElseNode.Init = conditionNode.Init
			switch elseBodyList := a.Else.(type) {
			case *ast.BlockStmt:
				list := elseBodyList.List
				nodes := make([]*ConditionNode, 0, 10)
				for _, v := range list {
					if ifNode, ok := v.(*ast.IfStmt); ok {
						condition := s.parseCondition(ifNode)
						if condition != nil {
							nodes = append(nodes, condition)
						}
					}
				}
				conditionElseNode.Children = nodes
			default:
				panic("未找到的body类型")
			}
			conditionNode.Else = &conditionElseNode
		}
		// 4. 解析body, 继续解析其中的if和switch
		if a.Body != nil {
			list := a.Body.List
			nodes := make([]*ConditionNode, 0, 10)
			for _, v := range list {
				if ifNode, ok := v.(*ast.IfStmt); ok {
					condition := s.parseCondition(ifNode)
					if condition != nil {
						nodes = append(nodes, condition)
					}
				}
			}
			conditionNode.Children = nodes
		}
		return &conditionNode
	case *ast.SwitchStmt:
		// 1. 解析init
		if a.Init != nil {
			conditionNode.Init = s.parseIfInit(a.Init)
		}
		// 2. 解析switch
		var ifo CondInfo
		//ifo.XParam = s.ParamParse(a.Switch)
		ifo.YParam = s.ParamParseValue(a.Tag)
		ifo.Op = token.EQL
		conditionNode.Cond = &ifo

		conditionNode.IsElse = false
		// 3. 解析body, 继续解析其中的if和switch
		if a.Body != nil {
			list := a.Body.List
			// 4. 先解析这个switch的case,每个case中可能包含if、switch
			for _, v := range list {
				if caseNode, ok := v.(*ast.CaseClause); ok {

				}
			}

			nodes := make([]*ConditionNode, 0, 10)
			for _, v := range list {
				if ifNode, ok := v.(*ast.IfStmt); ok {
					condition := s.parseCondition(ifNode)
					if condition != nil {
						nodes = append(nodes, condition)
					}
				}
			}
			conditionNode.Children = nodes
		}
		return &conditionNode
	default:
		panic("未找到可用类型")
	}

	return nil
}

func (s *SourceInfo) parseIfInit(n ast.Stmt) *InitInfo {
	var ifo InitInfo
	switch a := n.(type) {
	case *ast.AssignStmt:
		assignment := s.ParseAssignment(a)
		ifo.InitType = enum.CONDITION_INIT_TYPE_ASSIGNMENT
		ifo.InitValue = assignment
	default:
		panic("未找到可用类型")
	}
	return &ifo
}

func (s *SourceInfo) parseIfCond(n ast.Expr) *CondInfo {
	var ifo CondInfo
	switch a := n.(type) {
	case *ast.BinaryExpr:
		ifo.XParam = s.ParamParse(a.X)
		ifo.YParam = s.ParamParseValue(a.Y)
		ifo.Op = a.Op
	default:
		panic("未知的condition类型")
	}
	return &ifo
}
