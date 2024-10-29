package param

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// Index 选择类型，a.b.c.d
type Index struct {
	BasicParam
	// Index的SpecificType 代表的是"a.b.c.d"中的"d"对应的类型
	SpecificType
	BasicValue
	IndexReference
}

// IndexReference 选择类型的引用，用来实现 a.b.c.d
type IndexReference struct {
	IndexType string
	Child     *IndexReference
}

func (s *Index) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INDEX
}

func (s *Index) GetInstance() Parameter {
	return s
}

func (s *Index) GetZeroValue() Parameter {
	zeroValue := s.SpecificType.ZeroValue
	s.Value = zeroValue
	return s
}

func (s *Index) GetFormula() string {
	reference := s.IndexReference
	return GetFormulaFromReference(reference)
}

func GetFormulaFromReference(reference IndexReference) string {
	child := reference.Child
	if child == nil {
		return reference.IndexType
	} else {
		childFormula := GetFormulaFromReference(lo.FromPtr(reference.Child))
		return reference.IndexType + "." + childFormula
	}
}

// ParseIndex 解析ast
func ParseIndex(expr *ast.IndexExpr, name string) *Index {
	return nil
}

// ParseIndexList 解析ast
func ParseIndexList(expr *ast.IndexListExpr, name string) *Index {
	return nil
}
