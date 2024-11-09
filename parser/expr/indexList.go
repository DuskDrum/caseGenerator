package expr

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// IndexList 索引列表表达式,通常用于表示多维数组的索引
// 例如 a[i][j][k]
type IndexList struct {
}

func (s *IndexList) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_INDEX
}

func (s *IndexList) GetFormula() string {
	return ""
}

// ParseIndexList 解析ast， 暂不处理
func ParseIndexList(expr *ast.IndexListExpr, name string) *Index {
	return nil
}
