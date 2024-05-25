package enum

// AssignmentRightType 赋值类型右边的类型
type AssignmentRightType struct {
	Code string
	Desc string
}

var (
	// RIGHT_TYPE_FUNCTION 匿名函数一般不需要处理
	RIGHT_TYPE_FUNCTION  = AssignmentRightType{"function", "匿名函数"}
	RIGHT_TYPE_CALL      = AssignmentRightType{"call", "方法调用"}
	RIGHT_TYPE_COMPOSITE = AssignmentRightType{"composite", "原始数据"}
)

var AssignmentRightTypeList = []AssignmentRightType{
	RIGHT_TYPE_FUNCTION,
	RIGHT_TYPE_CALL,
	RIGHT_TYPE_COMPOSITE,
}
