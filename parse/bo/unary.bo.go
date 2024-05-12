package bo

// Unary 获取一元结构的值，一元表达式
type Unary interface {
	GetUnaryValue() any
}

type BasicLitUnary struct {
	// 原始值
	BasicLitValue string
}

func (b *BasicLitUnary) GetUnaryValue() any {
	return b.BasicLitValue
}

// InvocationUnary 常量或者变量的 变量名引用
type InvocationUnary struct {
	InvocationName string
}

func (c *InvocationUnary) GetUnaryValue() any {
	return c.InvocationName
}

// CallUnary 方法调用, package.FuncName()
type CallUnary struct {
	CallValue string
}

func (f *CallUnary) GetUnaryValue() any {
	return f.CallValue
}

// FuncUnary 定义的内部函数，或者入参里的函数，里面什么也没有
type FuncUnary struct {
}

func (f *FuncUnary) GetUnaryValue() any {
	return f
}

// ParamUnary 参数的层级调用， a.b.c.d
type ParamUnary struct {
	ParamValue string
}

func (p *ParamUnary) GetUnaryValue() any {
	return p.ParamValue
}

// TypeAssertUnary 类型断言，层级调用(a.b.c.d) + 断言的Type
type TypeAssertUnary struct {
	ParamValue string
	AssertType string
}

func (t *TypeAssertUnary) GetUnaryValue() any {
	return t.ParamValue
}
