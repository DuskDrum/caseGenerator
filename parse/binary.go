package parse

import "go/token"

// Binary 二元表达式，赋值或者判断。 左右两边都认为是Unary
type Binary interface {
	// GetX 获取左边
	GetX() Unary
	// GetY 获取右边
	GetY() Unary
}

// CondBinary 条件判断的二元表达式
type CondBinary struct {
	X Unary
	Y Unary
	// 判断条件
	Op token.Token
	// 位置信息
	PointerSite token.Pos
}

func (c *CondBinary) GetX() Unary {
	return c.X
}

func (c *CondBinary) GetY() Unary {
	return c.Y
}

// AssignmentBinary 赋值二元
type AssignmentBinary struct {
	// 左边，一定是属性层级调用
	X ParamUnary
	// 右边
	Y Unary
}

func (a *AssignmentBinary) GetX() Unary {
	return &a.X
}

func (a *AssignmentBinary) GetY() Unary {
	return a.Y
}

// DefineBinary 请求参数，或者 var x X 这种定义类型的
type DefineBinary struct {
	// 左边，常量或者变量的 变量名引用
	X InvocationUnary
	// 右边
	Y Unary
}

func (r *DefineBinary) GetX() Unary {
	return &r.X
}

func (r *DefineBinary) GetY() Unary {
	return r.Y
}
