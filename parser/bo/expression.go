package bo

import "caseGenerator/common/enum"

// ExpressionContext  expression上下文
type ExpressionContext struct {
	VariableParamMap  map[string]enum.BasicParameterType // 解析出变量map
	RequestParamMap   map[string]enum.BasicParameterType // 方法请求参数信息
	TemporaryVariable TemporaryVariable                  // 临时参数
	ExprContext
}

type TemporaryVariable struct {
	VariableIndex int    // 变量的下标，指的是赋值语句中比如一个call返回两个参数，那么这里要标识此参数是第几个参数，默认是0
	VariableName  string // 参数名
}

// GetKnownParamType 获取已知变量的类型
// go语法中，如果方法请求和局部变量名是同一个，那么取方法请求中，所以先在requestMap里找
func (e ExpressionContext) GetKnownParamType(paramName string) enum.BasicParameterType {
	if parameterType, ok := e.RequestParamMap[paramName]; ok {
		return parameterType
	}
	if parameterType, ok := e.VariableParamMap[paramName]; ok {
		return parameterType
	}
	return enum.BASIC_PARAMETER_TYPE_UNKNOWN
}
