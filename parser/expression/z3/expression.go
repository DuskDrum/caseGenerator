package z3

import (
	"caseGenerator/go-z3"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	_struct "caseGenerator/parser/struct"
)

// Z3Express 公式
type Z3Express struct {
	Z3Ast *z3.AST // z3的 ast
}

func Express(param _struct.Parameter, context bo.ExpressionContext) *Z3Express {
	// 解析得到condition和变量列表
	ast, _ := ExpressParam(param, context)
	return &Z3Express{Z3Ast: ast}
}

// ExpressParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
// *z3.AST : 第一个参数返回condition
// []z3.AST: 返回变量列表
func ExpressParam(param _struct.Parameter, context bo.ExpressionContext) (*z3.AST, []*z3.AST) {
	switch exprType := param.(type) {
	case *expr.Binary:
		return ExpressBinary(exprType, context)
	case *expr.Unary:
		return ExpressUnary(exprType, context)
	case *expr.Parent:
		return ExpressParent(exprType, context)
	case *expr.Ident: // 只处理局部变量、方法请求变量(不处理常量)
		return ExpressIdent(exprType, context)
	case *expr.Selector: // 只处理局部变量、方法请求变量(不处理常量)
		return ExpressSelector(exprType, context)
	case *expr.Call:
		return ExpressCall(exprType, context)
	case *expr.BasicLit:
		return ExpressBasicLit(exprType, context)
	default:
		panic("unknown expression type")
	}
}

//func ExpressTarget(param, targetParam _struct.Parameter) []*Z3Express {
//	expressionList := make([]*Z3Express, 0, 10)
//	eList := ExpressTargetParam(param, targetParam)
//	expressionList = append(expressionList, eList...)
//	return expressionList
//}

// ExpressTargetParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
//func ExpressTargetParam(param, targetParam _struct.Parameter) []*Z3Express {
//	switch exprType := param.(type) {
//	case *expr.Binary:
//		return ExpressBinary(exprType)
//	case *expr.Unary:
//		return ExpressUnary(exprType)
//	case *expr.Parent:
//		return ExpressParent(exprType)
//	case *expr.Ident:
//		return ExpressTargetIdent(exprType, targetParam)
//	case *expr.Selector:
//		return ExpressTargetSelector(exprType, targetParam)
//	case *expr.Call:
//		return ExpressTargetCall(exprType, targetParam)
//	case *expr.BasicLit:
//		return ExpressTargetBasicLit(exprType, targetParam)
//	default:
//		//elementList := []string{param.GetFormula()}
//		expression := &Z3Express{
//			//ElementList: elementList,
//			//Expr:        strings.Join(elementList, " "),
//		}
//		return []*Z3Express{expression}
//	}
//}
