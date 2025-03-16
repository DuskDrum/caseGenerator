package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

// Z3Express 公式
type Z3Express struct {
	Z3Ast *z3.AST // z3的 ast
}

func Express(param _struct.Parameter) *Z3Express {
	ast := ExpressParam(param)
	return &Z3Express{Z3Ast: ast}
}

// ExpressParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
func ExpressParam(param _struct.Parameter) *z3.AST {
	switch exprType := param.(type) {
	case *expr.Binary:
		return ExpressBinary(exprType)
	case *expr.Unary:
		return ExpressUnary(exprType)
	case *expr.Parent:
		return ExpressParent(exprType)
	case *expr.Ident:
		return ExpressIdent(exprType)
	case *expr.Selector:
		return ExpressSelector(exprType)
	case *expr.Call:
		return ExpressCall(exprType)
	case *expr.BasicLit:
		return ExpressBasicLit(exprType)
	default:
		//elementList := []string{param.GetFormula()}
		expression := &Z3Express{
			//ElementList: elementList,
			//Expr:        strings.Join(elementList, " "),
		}
		return []*Z3Express{expression}
	}
}

func ExpressTarget(param, targetParam _struct.Parameter) []*Z3Express {
	expressionList := make([]*Z3Express, 0, 10)
	eList := ExpressTargetParam(param, targetParam)
	expressionList = append(expressionList, eList...)
	return expressionList
}

// ExpressTargetParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
func ExpressTargetParam(param, targetParam _struct.Parameter) []*Z3Express {
	switch exprType := param.(type) {
	case *expr.Binary:
		return ExpressBinary(exprType)
	case *expr.Unary:
		return ExpressUnary(exprType)
	case *expr.Parent:
		return ExpressParent(exprType)
	case *expr.Ident:
		return ExpressTargetIdent(exprType, targetParam)
	case *expr.Selector:
		return ExpressTargetSelector(exprType, targetParam)
	case *expr.Call:
		return ExpressTargetCall(exprType, targetParam)
	case *expr.BasicLit:
		return ExpressTargetBasicLit(exprType, targetParam)
	default:
		//elementList := []string{param.GetFormula()}
		expression := &Z3Express{
			//ElementList: elementList,
			//Expr:        strings.Join(elementList, " "),
		}
		return []*Z3Express{expression}
	}
}
