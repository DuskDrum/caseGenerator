package stmt

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	expression2 "caseGenerator/parser/expression/govaluate"
	_struct "caseGenerator/parser/struct"
	"go/ast"
	"go/token"
	"strings"

	"github.com/Knetic/govaluate"
)

// IncDec 自增、自减语句
// 表示自增（++）和自减（--）语句
type IncDec struct {
	Content _struct.Parameter // 变量，一般是ast.Ident类型代表的变量
	Token   token.Token       // 类型，是 token.INC、token.DEC
}

func (i *IncDec) FormulaExpress() ([]bo.KeyFormula, map[string]*expr.Call) {
	stmtExpressionList := make([]bo.KeyFormula, 0, 10)
	se := bo.KeyFormula{
		Key:  i.Content.GetFormula(),
		Type: enum.STMT_TYPE_INCDEC,
	}
	// 直接取第一条即可，只有逻辑与、逻辑或才会有多条
	expressionList := expression2.Express(i.Content)
	expression := expressionList[0]
	if expression == nil {
		return nil, nil
	}
	se.IdentMap = expression.IdentMap
	se.CallMap = expression.CallMap
	se.SelectorMap = expression.SelectorMap
	if i.Token == token.INC {
		// a -- 相当于 a = a - 1
		elementList := make([]string, 0, 10)
		elementList = append(elementList, i.Content.GetFormula())
		elementList = append(elementList, govaluate.MINUS.String())
		elementList = append(elementList, "1")
		se.ElementList = elementList
		se.Expr = strings.Join(se.ElementList, " ")
	} else if i.Token == token.DEC {
		// a ++ 相当于 a = a + 1
		elementList := make([]string, 0, 10)
		elementList = append(elementList, i.Content.GetFormula())
		elementList = append(elementList, govaluate.PLUS.String())
		elementList = append(elementList, "1")
		se.ElementList = elementList
		se.Expr = strings.Join(se.ElementList, " ")
	} else {
		panic("incDec type illegal")
	}

	se.Formula = se.Expr

	stmtExpressionList = append(stmtExpressionList, se)
	return stmtExpressionList, expression.CallMap
}

// ParseIncDec 解析ast
func ParseIncDec(stmt *ast.IncDecStmt) *IncDec {
	incDec := &IncDec{}

	if stmt.X != nil {
		xp := expr.ParseParameter(stmt.X)
		if xp != nil {
			incDec.Content = xp
		}
	}
	incDec.Token = stmt.Tok

	return incDec
}
