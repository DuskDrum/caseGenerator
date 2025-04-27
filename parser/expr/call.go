package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"
	"strings"
)

// Call 表示函数或方法的调用表达式。
// 例如，表达式 fmt.Println("Hello, World!")、add(1, 2) 或 myStruct.Method(3) 都是函数调用表达式，
type Call struct {
	// 可能是 a()， 可能是 a.b()
	Function _struct.Parameter
	// 参数列表
	Args []_struct.Parameter
	// 响应详情
	ResponseList []Field
}

func (s *Call) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_CALL
}

func (s *Call) GetFormula() string {
	formula := s.Function.GetFormula() + "("
	for _, a := range s.Args {
		formula += " " + a.GetFormula() + ","
	}
	if strings.HasSuffix(formula, ",") {
		// 去掉最后一个字符
		formula = formula[:len(formula)-1]
	}
	formula += ")"
	return formula
}

// ParseCall 解析ast
func ParseCall(expr *ast.CallExpr, context bo.ExprContext) *Call {
	ca := &Call{}
	ca.Function = ParseParameter(expr.Fun, context)
	// Fun

	callList := make([]_struct.Parameter, 0, 10)
	for _, v := range expr.Args {
		call := ParseParameter(v, context)
		callList = append(callList, call)
	}
	ca.Args = callList
	// 解析响应

	return ca
}
