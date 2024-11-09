package stmt

import (
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
)

// Assign 赋值语句
// 赋值语句可以是简单的变量赋值（如x = 5），也可以是多元赋值（如x, y = 1, 2）
type Assign struct {
	Token           token.Token // 赋值类型，一般是 token.DEFINE、token.ASSIGN
	AssignParamList []AssignParam
}

type AssignParam struct {
	Left  _struct.Parameter
	Right _struct.Parameter
}

// ParseAssign 解析ast
func ParseAssign(stmt *ast.AssignStmt) Assign {
	assign := Assign{}
	// 赋值的左右一定是数量一样的
	rhs := stmt.Rhs
	lhs := stmt.Lhs
	if len(rhs) != len(lhs) {
		panic("assign len is not equal")
	}
	list := make([]AssignParam, 0, 10)
	for i, l := range lhs {
		lp := expr.ParseParameter(l)
		rp := expr.ParseParameter(rhs[i])
		ap := AssignParam{
			Left:  lp,
			Right: rp,
		}
		list = append(list, ap)
	}
	assign.AssignParamList = list
	assign.Token = stmt.Tok

	return assign
}
