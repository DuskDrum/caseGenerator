package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

func ExpressIdent(param *expr.Ident) (*z3.AST, []*z3.AST) {
	astList := make([]*z3.AST, 0, 10)
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	symbolName := ""

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
	return nil, nil
}

func ExpressTargetIdent(param *expr.Ident, targetParam _struct.Parameter) []*Z3Express {
	//elementList := []string{param.GetFormula()}
	//identMap := map[string]*expr.Ident{param.IdentName: param}

	expression := &Z3Express{
		//ElementList: elementList,
		//IdentMap:    identMap,
		//Expr:        targetParam.GetFormula() + " = " + param.GetFormula(),
	}
	return []*Z3Express{expression}
}
