package _struct

import (
	"caseGenerator/common/enum"
)

// Stmt 参数接口类型，将参数需要的方法定义出来
type Stmt interface {
	// GetType 获取类型(每个子类维护一个类型)
	GetType() enum.ParameterType
	// GetFormula 获取公式，用于展示
	GetFormula() string
}
