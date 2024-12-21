package decl

import (
	"caseGenerator/parser/expr"
	"caseGenerator/parser/stmt"
	"go/ast"

	"github.com/samber/lo"
)

// Func 解析
type Func struct {
	Receiver []expr.Field  // receiver (methods); or nil (functions)
	Name     expr.Ident    // function/method name
	Type     expr.FuncType // function signature: type and value parameters, results, and position of "func" keyword
	Body     *stmt.Block
}

// ParseFunc 解析Func
func ParseFunc(decl *ast.FuncDecl) *Func {
	fieldList := make([]expr.Field, 0, 10)
	if decl.Recv != nil {
		for _, v := range decl.Recv.List {
			pf := expr.ParseField(v)
			if pf != nil {
				fieldList = append(fieldList, lo.FromPtr(pf))
			}
		}
	}
	name := expr.ParseIdent(decl.Name)
	funcType := expr.ParseFuncType(decl.Type)
	body := stmt.ParseBlock(decl.Body)

	f := &Func{
		Receiver: fieldList,
		Name:     lo.FromPtr(name),
		Type:     lo.FromPtr(funcType),
		Body:     body,
	}
	return f
}

// ParseBody 解析方法
func ParseBody(sb *ast.BlockStmt) {
	crMap := make(map[string][]stmt.ConditionResult, 10)
	// 1. 遍历方法里的每个 stmt，

	// 1.1 解析每个stmt的内容

	// 1.2 解析每个stmt 的 assignment

	// 1.3 执行每一个 stmt 的公式；得到最新的赋值公式

	// 1.4 解析每个condition需要的变量，得到不等式公式

	// 2. 得到公式的列表

	// 3. 遍历的处理这个列表

	// 4. 得到mock结果

	for _, v := range sb.List {
		// 解析代码块
		p := stmt.ParseStmt(v)
		// 生成逻辑表达式
		express := p.LogicExpression()
		// 得到mock结果
		conditionResults := p.CalculateCondition(express)
		crMap["xxx"] = conditionResults

	}
}
