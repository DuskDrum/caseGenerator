package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"go/ast"

	"github.com/samber/lo"
)

// FuncDecl 定义方法特殊结构，方法名、请求参数列表、响应参数列表
type FuncDecl struct {
	Type     FuncType
	Receiver []Field // receiver
	Name     Ident
}

func (s *FuncDecl) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNC_DECL
}

func (s *FuncDecl) GetFormula() string {
	formulaStr := "func "
	if len(s.Receiver) != 0 {
		//todo 考虑 receiver
	}

	formulaStr += s.Name.GetFormula()

	formulaStr += s.Type.GetFormula()

	formulaStr += "{}"

	return formulaStr
}

// ParseFuncDecl 解析ast
func ParseFuncDecl(expr *ast.FuncDecl, context bo.ExprContext) *FuncDecl {
	fd := &FuncDecl{Name: Ident{}}
	if expr.Type != nil {
		funcType := ParseFuncType(expr.Type, context)
		if funcType != nil {
			fd.Type = lo.FromPtr(funcType)
		}
	}
	if expr.Recv != nil {
		fields := make([]Field, 0, 10)
		for _, recv := range expr.Recv.List {
			funcType := ParseField(recv, context)
			if funcType != nil {
				fields = append(fields, lo.FromPtr(funcType))
			}
		}
		fd.Receiver = fields
	}

	return nil
}
