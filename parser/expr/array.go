package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
	"strconv"

	"github.com/samber/lo"
)

// Array 数组结构, 数组有三种
// [10]int, 这是数组类型。指定了长度，不能再变动
// [...]int, 这是不定长数组类型。还是固定长度的，在申明后可以自动推测出他的长度，之后就再不能改变了。
// []int, 这是切片类型。可变
type Array struct {
	_struct.RecursionParam
	// 长度， 只有数组类型有
	Length int
	// 是否是不定长数组
	IsEllipsis bool
}

func (s *Array) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_ARRAY
}

func (s *Array) GetFormula() string {
	childFormula := s.Parameter.GetFormula()
	return "[]" + childFormula
}

// ParseArray 解析ast
func ParseArray(expr *ast.ArrayType) *Array {
	arr := &Array{}
	// 1. 判断是否是数组或者不定长数组
	if expr.Len != nil {
		switch lenType := expr.Len.(type) {
		case *ast.BasicLit:
			// 使用 strconv.Atoi 将字符串转换为 int
			i, err := strconv.Atoi(lenType.Value)
			if err != nil {
				panic("ast arrayType len convert to int err: " + err.Error())
			}
			arr.Length = i
		case *ast.Ellipsis:
			arr.IsEllipsis = true
		default:
			panic("ast arrayType len type is illegal")
		}
	}
	// 2. 解析child
	ap := ParseRecursionValue(expr.Elt)
	arr.RecursionParam = lo.FromPtr(ap)

	return arr
}
