package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
)

// Index 索引表达式
// 通常用于数组、切片和映射的索引操作
// 例如 a[i]、m["key"]
type Index struct {
	// Subject代表了是上面的 a 或者 m
	Subject _struct.Parameter
	// Index代表的是上面的 i 或者"key"。 一般是 *ast.BasicLit， int 或者 string
	Index _struct.Parameter
}

func (s *Index) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INDEX
}

func (s *Index) GetFormula() string {
	return s.Subject.GetFormula() + "[" + s.Index.GetFormula() + "]"
}

// ParseIndex 解析ast
func ParseIndex(expr *ast.IndexExpr) *Index {
	i := &Index{}
	i.Subject = ParseParameter(expr.X)
	i.Index = ParseParameter(expr.Index)

	return i
}
