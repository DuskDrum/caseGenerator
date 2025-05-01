package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
)

func ExpressIdent(param *expr.Ident, context bo.ExpressionContext) (*z3.AST, []*z3.AST) {
	astList := make([]*z3.AST, 0, 10)
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	// 1. 根据上下文解析出变量名
	symbolName := context.TemporaryVariable.VariableName
	// 2. 判断这个变量的类型
	// 2.1 判断请求变量中有没有这个变量

	// 2.2 判断局部变量中有没有这个变量

	// 2.3 判断ident对应的变量是否在请求变量中（ident 可能指向其他的变量）

	// 2.4 判断ident对应的变量是否在局部变量中（ident 可能指向其他的变量）

	// 2.5 判断能否从这个ident的值推导出类型 (比如ident的值是: 8 那么就是int)
	// 这里需要注意2.5里ident的值有可能是全局变量，是个string

	// iType
	switch param.IdentName {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64": // rune就是int32
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.IntSort()))
	case "float32":
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.FloatSort()))
	case "float64":
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.DoubleSort()))
	case "bool":
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.BoolSort()))
	case "string":
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.StringSort()))
	default:
		// 啥也不干
	}
	return nil, astList
}
