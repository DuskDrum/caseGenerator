package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
	"fmt"
	"strings"
)

func ExpressCall(param *expr.Call) (*z3.AST, []*z3.AST) {
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	astList := make([]*z3.AST, 0, 10)

	// call返回值的规则就用方法名+位数
	funcFormula := strings.ReplaceAll(param.Function.GetFormula(), ".", "_")

	for i, v := range param.ResponseList {
		// call返回值的规则就用方法名+位数
		symbolName := fmt.Sprintf(funcFormula+"%d", i)
		switch iType := v.Type.(type) {
		case *expr.Ident:
			// iType
			switch iType.IdentName {
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
		default:
			// 啥也不干
		}

	}

	return nil, astList
}

func ExpressTargetCall(param *expr.Call, targetParam _struct.Parameter) *Z3Express {
	key := strings.ReplaceAll(param.GetFormula(), ".", "_")
	key = strings.ReplaceAll(key, "(", "_")
	key = strings.ReplaceAll(key, ")", "")
	key = strings.ReplaceAll(key, " ", "")

	//callMap := map[string]*expr.Call{"astCall_" + key: param}

	//elementList := []string{"astCall_" + key}

	//expression := &Z3Express{
	//ElementList: elementList,
	//CallMap:     callMap,
	//Expr:        targetParam.GetFormula() + " = " + "astCall_" + key,
	//}
	//return []*Z3Express{expression}
	return nil
}
