package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
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
	BinaryParam
}

type BinaryParam struct {
	X  _struct.Parameter
	Op token.Token
	Y  _struct.Parameter
}

func (s *Binary) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BINARY
}

func (s *Binary) GetFormula() string {
	panic("implement me")
}

// ParseBinary 解析Binary
func ParseBinary(expr *ast.BinaryExpr, af *ast.File) *Binary {
	ab := &Binary{}
	ab.BinaryParam = lo.FromPtr(ParseBinaryParam(expr, af))
	return ab
}

func ParseBinaryParam(expr *ast.BinaryExpr, af *ast.File) *BinaryParam {
	bp := &BinaryParam{}
	bp.Op = expr.Op
	// 解析X、解析Y
	// 调用这个方法会去 switch 所有 expr 类型
	// 如果类型是BinaryParam--->就会来调用ParseBinaryParam方法，形成递归
	xP := ParseParameter(expr.X, af)
	yP := ParseParameter(expr.Y, af)
	bp.X = xP
	bp.Y = yP
	return bp
}

type Mock interface {
}

// CallMock call类型的mock
type CallMock struct {
}

// ParamMock param类型的mock
type ParamMock struct {
	MockName  string
	MockValue string
}

//func MockExpr(expr _struct.Parameter) []Mock {
//	switch exprType := expr.(type) {
//	case *Binary:
//		return MockBinary(exprType)
//	case *Unary:
//		return MockUnary(exprType)
//	default:
//		panic("未知类型...")
//	}
//}

// MockBinary BinaryParam
// Op有三种类型
// 1. 逻辑与、逻辑或 LOR、LAND
// 2. 大于、等于、小于  EQL、LSS、GTR、LEQ、GEQ
// 3. 运算符。加减乘除余 ADD、SUB、MUL、QUO、REM
//func MockBinary(param *Binary) []Mock {
//	mockList := make([]Mock, 0, 10)
//	// 如果类型是||逻辑或，剪枝只处理X
//	if param.Op == token.LOR {
//		return MockExpr(param.X)
//		// 如果类型是&&逻辑与，处理X和Y
//	} else if param.Op == token.LAND {
//		x := MockExpr(param.X)
//		y := MockExpr(param.Y)
//		mockList = append(mockList, x)
//		mockList = append(mockList, y)
//		return mockList
//	}
//	// 如果类型是	== 相等，那么组装mock
//	if param.Op == token.EQL {
//		// 等于，底层一定不是binary
//		//x := MockExpr(param.X)
//		//y := MockExpr(param.Y)
//	} else if param.Op == token.LSS {
//
//	} else if param.Op == token.GTR {
//
//	} else if param.Op == token.LEQ {
//
//	} else if param.Op == token.GEQ {
//
//	}
//
//	// 暴力破解
//	if param.Op == token.ADD {
//
//	} else if param.Op == token.SUB {
//
//	} else if param.Op == token.MUL {
//
//	} else if param.Op == token.QUO {
//
//	} else if param.Op == token.REM {
//
//	}
//
//	//var a, b struct{}
//	//if a > b {
//	//
//	//}
//
//	return nil
//}

// 计算表达式  github.com/Knetic/govaluate、github.com/antonmedv/expr
