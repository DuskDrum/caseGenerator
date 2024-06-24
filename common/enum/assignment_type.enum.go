package enum

// AssignmentType 赋值类型的类型
type AssignmentType struct {
	Code        string
	Desc        string
	ContentDesc string
}

var (
	// RIGHT_TYPE_FUNCTION 匿名函数一般不需要处理
	RIGHT_TYPE_FUNCTION = AssignmentType{Code: "function", Desc: "匿名函数", ContentDesc: ""}
	RIGHT_TYPE_CALL     = AssignmentType{Code: "call", Desc: "方法调用", ContentDesc: ""}
	// RIGHT_TYPE_COMPOSITE 构造类型，比如{1,2,3,4}
	RIGHT_TYPE_COMPOSITE = AssignmentType{Code: "composite", Desc: "构造类型", ContentDesc: "CompositeLitValue"}
	// RIGHT_TYPE_BASICLIT 基本类型，直接是值
	RIGHT_TYPE_BASICLIT = AssignmentType{Code: "basiclit", Desc: "基本类型", ContentDesc: "直接是值"}
)

var AssignmentRightTypeList = []AssignmentType{
	RIGHT_TYPE_FUNCTION,
	RIGHT_TYPE_CALL,
	RIGHT_TYPE_COMPOSITE,
	RIGHT_TYPE_BASICLIT,
}
