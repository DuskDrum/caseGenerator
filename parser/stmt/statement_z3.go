package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expression/z3"
)

type Z3Stmt interface {
}

// ExpressionZ3Stmt 参数接口类型，将参数需要的方法定义出来
// 入参：ExpressionContext，包含了局部变量的列表，初始化出请求的参数
type ExpressionZ3Stmt interface {
	// Z3FormulaExpress  生成逻辑表达式, 入参为表达式上下文， 出参 z3Express
	Z3FormulaExpress(context bo.ExpressionContext) []z3.Z3Express
}
