package param

import "caseGenerator/common/enum"

// Parameter 参数接口类型，将参数需要的方法定义出来
type Parameter interface {
	// GetType 获取类型(每个子类维护一个类型)
	GetType() enum.ParameterType
	// GetInstance 有些场景这个接口不能满足，直接返回自身去特殊处理
	GetInstance() Parameter
	// GetZeroValue 返回对应的零值
	GetZeroValue() Parameter
	// GetFormula 获取公式，用于展示
	GetFormula() string
}
