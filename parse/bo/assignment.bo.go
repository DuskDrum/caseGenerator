package bo

import "caseGenerator/parse/enum"

type AssignmentDetailInfo struct {
	// 等式左边的名字
	LeftName []string
	// 等式右边的类型
	RightType enum.AssignmentRightType
	// 等式右边的公式(如果是赋值，那么代表了值；如果是函数，代表了函数的调用；如果是类型断言，那么是类型断言的公式)
	RightFormula string
	// 等式右边的参数
	RightFunctionParam []ParamParseRequest
	// 等式左边参数的类型
	LeftResultType []string
}
