package param

import (
	"caseGenerator/common/enum"
	"strings"
)

// Func 定义方法特殊结构，方法名、请求参数列表、响应参数列表
// todo 确认是否有Receiver
type Func struct {
	BasicParam
	BasicValue
	RequestList  []Parameter
	ResponseList []Parameter
}

func (s *Func) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNC
}

func (s *Func) GetInstance() Parameter {
	return s
}

func (s *Func) GetZeroValue() Parameter {
	s.Value = nil
	return s
}

func (s *Func) GetFormula() string {
	var paramType = "func("
	// 解析func的入参
	list := s.RequestList
	for _, v := range list {
		paramType = paramType + v.GetType().Name + ", "
	}
	// 去掉最后一个逗号
	lastCommaIndex := strings.LastIndex(paramType, ",")
	if lastCommaIndex != -1 {
		paramType = paramType[:lastCommaIndex]
	}
	paramType = paramType + ")"
	// 解析func的出参
	if s.ResponseList == nil {
		return paramType
	}
	fields := s.ResponseList
	if len(fields) > 0 {
		paramType = paramType + " ("
		for _, v := range fields {
			paramType = paramType + v.GetType().Name + ", "
		}
		// 去掉最后一个逗号
		respLastCommaIndex := strings.LastIndex(paramType, ",")
		if respLastCommaIndex != -1 {
			paramType = paramType[:respLastCommaIndex]
		}
		paramType = paramType + ")"
	}
	return paramType
}
