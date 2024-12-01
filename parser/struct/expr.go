package _struct

import (
	"caseGenerator/common/enum"
)

// Parameter 参数接口类型，将参数需要的方法定义出来
type Parameter interface {
	// GetType 获取类型(每个子类维护一个类型)
	GetType() enum.ParameterType
	// GetFormula 获取公式，用于展示
	GetFormula() string
}

// RecursionParam 递归的参数，比如切片
type RecursionParam struct {
	Parameter
	Child *RecursionParam
}

// ValueAble 可以直接取值的类型
type ValueAble interface {
	GetValue() any
	SetValue(value any)
}
