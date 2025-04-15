package expr

import (
	"caseGenerator/common/enum"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

// BasicLit 基本字面量（literal），表示源代码中的基本类型常量
// 比如数字、字符串和布尔值
type BasicLit struct {
	Value        any
	SpecificType enum.SpecificType
	Kind         token.Token
}

func (s *BasicLit) GetValue() any {
	return s.Value
}

func (s *BasicLit) SetValue(value any) {
	s.Value = value
}

func (s *BasicLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BASICLIT
}

func (s *BasicLit) GetFormula() string {
	return fmt.Sprintf("%v", s.Value)
}

// ParseBasicLit 解析ast
func ParseBasicLit(expr *ast.BasicLit, _ *ast.File) *BasicLit {
	bl := &BasicLit{}
	// 解析 value
	bl.Kind = expr.Kind

	switch expr.Kind {
	case token.INT:
		// 将 string 转为 int
		num, err := strconv.Atoi(expr.Value)
		if err != nil {
			panic("string convert to int error :" + err.Error())
		}
		bl.Value = num
		bl.SpecificType = enum.SPECIFIC_TYPE_INT
	case token.STRING, token.CHAR:
		bl.Value = expr.Value
		bl.SpecificType = enum.SPECIFIC_TYPE_STRING
	case token.FLOAT:
		// 转换为 float64
		num, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			panic("string convert to float64 error :" + err.Error())
		}
		bl.Value = num
		bl.SpecificType = enum.SPECIFIC_TYPE_FLOAT64
	default:
		panic("basic lit value unhandled default case")
	}
	return bl
}
