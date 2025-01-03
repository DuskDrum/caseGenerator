package decl

import (
	"caseGenerator/parser/bo"
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
	// 常量map
	constantsMap := make(map[string]any, 10)
	// 内部变量map， 通过上下文推导出来
	innerVariablesMap := make(map[string]any, 10)
	// 外部变量map，从请求参数里拿到
	outerVariablesMap := make(map[string]any, 10)
	// 公式列表
	keyFormulaList := make([]bo.KeyFormula, 0, 10)
	callVariableMap := make(map[string]*expr.Call, 10)

	// 1. 遍历方法里的每个 stmt，
	for _, v := range sb.List {
		// todo
		// 1.1 解析每个stmt的内容
		p := stmt.ParseStmt(v)
		// 1.2 执行每一个 stmt 的公式；得到最新的赋值公式
		if pes, ok := p.(stmt.ExpressionStmt); ok {
			kfList, callMap := pes.FormulaExpress()
			keyFormulaList = append(keyFormulaList, kfList...)
			for k, value := range callMap {
				callVariableMap[k] = value
			}
		}
		// 1.3 解析每个condition需要的变量，得到不等式公式
		condition := p.CalculateCondition(express)
		crList = append(crList, condition...)
	}
	// 2. 遍历执行所有公式、

	// 3. 得到mock结果

}
