package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
)

func ExpressBasicLit(param *expr.BasicLit) (*z3.AST, []*z3.AST) {
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer ctx.Close()
	switch param.Value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64: // rune就是int32
		return ctx.Const(ctx.Symbol(param.GetFormula()), ctx.IntSort()), nil
	case float32:
		return ctx.Const(ctx.Symbol(param.GetFormula()), ctx.FloatSort()), nil
	case float64:
		return ctx.Const(ctx.Symbol(param.GetFormula()), ctx.DoubleSort()), nil
	case bool:
		return ctx.Const(ctx.Symbol(param.GetFormula()), ctx.BoolSort()), nil
	case string:
		return ctx.Const(ctx.Symbol(param.GetFormula()), ctx.StringSort()), nil
	default:
		panic("basic type not found")
	}
	return nil, nil
}
