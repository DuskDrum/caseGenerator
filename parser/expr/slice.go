package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"

	"github.com/samber/lo"
)

// Slice 切片操作的表达式，
// 例如二参数切片 a[1:3] 或 b[:4]
//   - low：切片的起始索引（包含），如果省略则默认为 0。
//   - high：切片的结束索引（不包含），如果省略则默认为数组或切片的长度。
//
// 特殊的结构是：三参数切片。a[low:high:max](这里的元素不一定只是数字 BasicLit，也有可能是 Ident 或者 Selector)
//   - low：切片的起始索引（包含）。
//   - high：切片的结束索引（不包含）。
//   - max：切片的最大索引（不包含）。
type Slice struct {
	Content _struct.Parameter // 代表了上面的a或者b
	// 切片起始索引(包含),有可能是 BasicLit 、 Ident 或者 Selector
	Low *_struct.Parameter
	// 切片结束索引(不包含),有可能是 BasicLit 、 Ident 或者 Selector
	High *_struct.Parameter
	// 切片的最大索引(不包含)，主要是为了指定切片后的容量，防止越界访问了原始数组的其他元素
	Max *_struct.Parameter
}

func (s *Slice) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_SLICE
}

func (s *Slice) GetFormula() string {
	formula := s.Content.GetFormula()
	formula += "["
	if s.Low != nil {
		formula += lo.FromPtr(s.Low).GetFormula()
	}
	formula += ":"
	if s.High != nil {
		formula += lo.FromPtr(s.High).GetFormula()
	}
	if s.Max != nil {
		formula += ":"
		formula += lo.FromPtr(s.Max).GetFormula()
	}
	return formula
}

// ParseSlice 解析ast
func ParseSlice(expr *ast.SliceExpr, context bo.ExprContext) *Slice {
	slice := &Slice{}
	slice.Content = ParseParameter(expr.X, context)
	if expr.Low != nil {
		low := ParseParameter(expr.Low, context)
		if low != nil {
			slice.Low = lo.ToPtr(low)
		}
	}
	if expr.High != nil {
		high := ParseParameter(expr.High, context)
		if high != nil {
			slice.High = lo.ToPtr(high)
		}
	}
	if expr.Max != nil {
		maxExpr := ParseParameter(expr.Max, context)
		if maxExpr != nil {
			slice.Max = lo.ToPtr(maxExpr)
		}
	}
	return slice
}
