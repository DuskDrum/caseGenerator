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
	Body     stmt.Block
}

// ParseFunc 解析ast
func ParseFunc(decl *ast.FuncDecl) *Func {
	f := &Func{}
	fieldList := make([]expr.Field, 0, 10)
	for _, v := range decl.Recv.List {
		pf := expr.ParseField(v)
		if pf != nil {
			fieldList = append(fieldList, lo.FromPtr(pf))
		}

	}

	return f
}
