package stmt

import (
	"caseGenerator/common/constants"
	"caseGenerator/common/enum"
	"caseGenerator/parser/expr"
	expression2 "caseGenerator/parser/expression/govaluate"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
	"strings"

	"github.com/Knetic/govaluate"
)

// Assign 赋值语句
// 赋值语句可以是简单的变量赋值（如x = 5），也可以是多元赋值（如x, y = 1, 2）
type Assign struct {
	Token token.Token // 赋值类型，
	// 一般是 token.DEFINE (:=)、
	// token.ASSIGN(=)、
	// token.ADD_ASSIGN // +=
	// token.SUB_ASSIGN // -=
	// token.MUL_ASSIGN // *=
	// token.QUO_ASSIGN // /=
	// token.REM_ASSIGN // %=
	// token.AND_ASSIGN     // &=
	// token.OR_ASSIGN      // |=
	// token.XOR_ASSIGN     // ^=
	// token.SHL_ASSIGN     // <<=
	// token.SHR_ASSIGN     // >>=
	// token.AND_NOT_ASSIGN // &^=
	AssignParamList []AssignParam
	Position        token.Pos // 代码的行数，同一个文件里比对才有意义
}

func (a *Assign) FormulaExpress() ([]expression2.KeyFormula, map[string]*expr.Call) {
	formulasList := make([]expression2.KeyFormula, 0, 10)
	callMap := make(map[string]*expr.Call, 10)
	for _, param := range a.AssignParamList {
		se := expression2.KeyFormula{
			Key:  param.Left.GetFormula(),
			Type: enum.STMT_TYPE_ASSIGN,
		}
		// 直接取第一条即可，只有逻辑与、逻辑或才会有多条
		expressionList := expression2.Express(param.Right)
		expression := expressionList[0]
		if expression == nil {
			continue
		}
		se.IdentMap = expression.IdentMap
		se.CallMap = expression.CallMap
		se.SelectorMap = expression.SelectorMap

		for k, v := range se.CallMap {
			callMap[k] = v
		}

		switch a.Token {
		case token.DEFINE, token.ASSIGN:
			// 设置初始化的值
			se.Expr = expression.Expr
			se.ElementList = expression.ElementList
		case token.ADD_ASSIGN:
			// a += b 相当于 a = a + b
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.PLUS.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.SUB_ASSIGN:
			// a -= b 相当于 a = a - b
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.MINUS.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.MUL_ASSIGN:
			// a *= b 相当于 a = a * b
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.MULTIPLY.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.QUO_ASSIGN:
			// a /= b 相当于 a = a / b
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.DIVIDE.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.REM_ASSIGN:
			// a %= b 相当于 a = a % b
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.MODULUS.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.AND_ASSIGN:
			// a &= b 相当于 a = a & b
			// 按位与运算 满足交换律
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.BITWISE_AND.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.OR_ASSIGN:
			// a |= b 相当于 a = a | b
			// 按位或运算 满足交换律
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.BITWISE_OR.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.XOR_ASSIGN:
			// a ^= b 相当于 a = a ^ b
			// 按位异或运算 满足交换律
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.BITWISE_XOR.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.SHL_ASSIGN:
			// a <<= b 相当于 a = a << b
			// 左移赋值运算符 左移相当于对数字进行乘以 2^b 的运算。
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.BITWISE_LSHIFT.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.SHR_ASSIGN:
			// a >>= b 相当于 a = a >> b
			// 右移赋值运算符 左移相当于对数字进行乘以 2^b 的运算。
			elementList := make([]string, 0, 10)
			elementList = append(elementList, param.Left.GetFormula())
			elementList = append(elementList, govaluate.BITWISE_RSHIFT.String())
			elementList = append(elementList, constants.CLAUSE)
			elementList = append(elementList, expression.ElementList...)
			elementList = append(elementList, constants.CLAUSE_CLOSE)
			se.ElementList = elementList
			se.Expr = strings.Join(se.ElementList, " ")
		case token.AND_NOT_ASSIGN:
			// a &^= b 相当于 a = a &^ b
			// 按位清除赋值运算符 用于按位清除并赋值。简单来说，它会根据另一个操作数将变量中对应的位清零
			// 不支持
			panic("token type illegal")
		default:
			panic("token type illegal")
		}
		formulasList = append(formulasList, se)
	}
	return formulasList, callMap
}

type AssignParam struct {
	Left  _struct.Parameter
	Right _struct.Parameter
}

// ParseAssign 解析ast
func ParseAssign(stmt *ast.AssignStmt, af *ast.File) *Assign {
	assign := &Assign{}
	// 赋值的左右一定是数量一样的
	rhs := stmt.Rhs
	lhs := stmt.Lhs
	if len(rhs) != len(lhs) {
		panic("assign len is not equal")
	}
	list := make([]AssignParam, 0, 10)
	for i, l := range lhs {
		lp := expr.ParseParameter(l, af)
		rp := expr.ParseParameter(rhs[i], af)
		ap := AssignParam{
			Left:  lp,
			Right: rp,
		}
		list = append(list, ap)
	}
	assign.AssignParamList = list
	assign.Token = stmt.Tok
	assign.Position = stmt.Pos()
	return assign
}
