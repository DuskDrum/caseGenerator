package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
)

// CompositeLit 复合字面量。用于表示数据初始化的复合数据结构。
// 可以是切片、数组、map、结构体
//
//		numbers := []int{1, 2, 3}
//	 ·  numbers2 := [10]int{1, 2, 3}
//		names := map[string]int{"Alice": 10, "Bob": 20}
//		person := struct{Name string}{Name: "Charlie"}
//
// 按照上面的四个例子，Type分别是:*ast.ArrayType、*ast.ArrayType、*ast.MapType、*ast.StructType
//
// 按照上面四个例子, Content分别是:[]*ast.BasicLit，[]*ast.BasicLit，[]*ast.KeyValueExpr, []*ast.KeyValueExpr
type CompositeLit struct {
	Type    _struct.Parameter
	Content []_struct.Parameter
	// x = []int{1, 2, 3}
	// y = struct{A int}{A: 10}
	Value any
}

func (s *CompositeLit) GetValue() any {
	return s.Value
}

func (s *CompositeLit) SetValue(value any) {
	s.Value = value
}

func (s *CompositeLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_COMPOSITELIT
}

func (s *CompositeLit) GetFormula() string {
	formula := ""
	formula += s.Type.GetFormula()
	formula += "{"
	for _, c := range s.Content {
		formula += c.GetFormula() + "\t"
	}
	formula += "}"
	return formula
}

// ParseCompositeLit 解析ast
func ParseCompositeLit(expr *ast.CompositeLit, af *ast.File) *CompositeLit {
	cl := &CompositeLit{}
	parameterType := ParseParameter(expr.Type, af)

	contentList := make([]_struct.Parameter, 0, 10)
	for _, elt := range expr.Elts {
		content := ParseParameter(elt, af)
		if content != nil {
			contentList = append(contentList, content)
		}
	}
	cl.Type = parameterType
	cl.Content = contentList
	return cl
}
