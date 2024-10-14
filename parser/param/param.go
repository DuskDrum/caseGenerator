package param

import "caseGenerator/common/enum"

// Parameter 参数接口类型，将参数需要的方法定义出来
type Parameter interface {
	// GetType 获取类型(每个子类维护一个类型)
	GetType() enum.ParameterType
	// GetName 获取name
	GetName() string
	// GetInstance 有些场景这个接口不能满足，直接返回自身去特殊处理
	GetInstance() Parameter
	// GetZeroValue 返回对应的零值
	GetZeroValue() Parameter
	// GetFormula 获取公式，用于展示
	GetFormula() string
}

type BasicParam struct {
	ParameterType enum.ParameterType
	Name          string
}

// RecursionParam 递归的参数，比如切片
type RecursionParam struct {
	Parameter
	Child *RecursionParam
}

func (b BasicParam) GetName() string {
	return b.Name
}

// BasicValue 具体的值
type BasicValue struct {
	// 这是具体的值
	Value any
}

// SpecificType 参数的具体类型，比如说 int、uint。这个需要组装为另一个枚举
type SpecificType struct {
	enum.SpecificType
}

func (s SpecificType) GetZeroValue() any {
	return s.ZeroValue
}
