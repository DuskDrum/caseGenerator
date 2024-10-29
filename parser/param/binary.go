package param

import (
	"caseGenerator/common/enum"
	"go/ast"
	"go/token"

	"github.com/samber/lo"
)

// Binary 二元表达式, 举例表示：
// a + b + c
// a + b + 1
// a || b && c.d
// 其中的元素要么是*ast.BasicLit、要么是 *ast.SelectorExpr、要么是*ast.CallExpr、要么是*ast.BinaryExpr、 要么是*ast.UnaryExpr
type Binary struct {
	BasicParam
	BinaryParam
}

type BinaryParam struct {
	X  Parameter
	Op token.Token
	Y  Parameter
}

func (s *Binary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BINARY
}

func (s *Binary) GetInstance() Parameter {
	return s
}

func (s *Binary) GetZeroValue() Parameter {
	return s
}

func (s *Binary) GetFormula() string {
	panic("implement me")
}

// ParseBinary 解析Binary
func ParseBinary(expr *ast.BinaryExpr, name string) *Binary {
	ab := Binary{
		BasicParam: BasicParam{
			ParameterType: enum.PARAMETER_TYPE_BINARY,
			Name:          name,
		},
	}
	ab.BinaryParam = lo.FromPtr(ParseBinaryParam(expr))
	return nil
}

func ParseBinaryParam(expr *ast.BinaryExpr) *BinaryParam {
	bp := &BinaryParam{}
	bp.Op = expr.Op
	// 解析X、解析Y
	// 调用这个方法会去 switch 所有 expr 类型
	// 如果类型是BinaryParam--->就会来调用ParseBinaryParam方法，形成递归
	xP := ParseParameter(expr.X)
	yP := ParseParameter(expr.Y)
	bp.X = xP
	bp.Y = yP
	return bp
}
