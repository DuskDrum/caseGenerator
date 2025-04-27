package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/struct"
	"go/ast"
)

// Chan 通道类型
// chan int、chan<- int、<-chan int 等
type Chan struct {
	Dir   ast.ChanDir // 通道方向，ast.SEND 发送 chan<-，ast.RECV 接收 <-chan， 默认是双向 chan
	Param _struct.Parameter
}

func (s *Chan) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_CHAN
}

func (s *Chan) GetFormula() string {
	// 输出通道的方向
	switch s.Dir {
	case ast.SEND:
		return "chan<-" + s.Param.GetFormula()
	case ast.RECV:
		return "<-chan" + s.Param.GetFormula()
	default:
		return "chan" + s.Param.GetFormula()
	}
}

// ParseChan 解析ast
func ParseChan(expr *ast.ChanType, context bo.ExprContext) *Chan {
	ch := &Chan{}
	ch.Dir = expr.Dir
	ch.Param = ParseParameter(expr.Value, context)
	return ch
}
