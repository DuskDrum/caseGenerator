package enum

import "encoding/json"

// AssignmentType 赋值类型的类型
type AssignmentType struct {
	Code        string
	Desc        string
	ContentDesc string
}

var (
	// ASSIGNMENT_TYPE_FUNCTION 匿名函数一般不需要处理，
	ASSIGNMENT_TYPE_FUNCTION = AssignmentType{Code: "function", Desc: "匿名函数", ContentDesc: "string"}
	// ASSIGNMENT_TYPE_COMPOSITE 构造类型，比如{1,2,3,4}
	ASSIGNMENT_TYPE_COMPOSITE = AssignmentType{Code: "composite", Desc: "构造类型", ContentDesc: "CompositeLitValue"}
	// ASSIGNMENT_TYPE_BASICLIT 基本类型，直接是值
	ASSIGNMENT_TYPE_BASICLIT = AssignmentType{Code: "basiclit", Desc: "基本类型", ContentDesc: "直接是值"}
	// ASSIGNMENT_TYPE_CALL 有结果的call，比如说 x1,x2 := aa.bb(cc,dd,ee)
	ASSIGNMENT_TYPE_CALL = AssignmentType{Code: "call", Desc: "方法调用", ContentDesc: "CallLitValue"}
	// ASSIGNMENT_TYPE_EXPRSTMT 没有结果的call，比如说fmt.Println()
	ASSIGNMENT_TYPE_EXPRSTMT = AssignmentType{Code: "exprstmt", Desc: "没有结果的方法调用", ContentDesc: "CallLitValue"}
	// ASSIGNMENT_TYPE_IDENT 关联的变量
	ASSIGNMENT_TYPE_IDENT = AssignmentType{Code: "ident", Desc: "关联的变量", ContentDesc: "关联变量的名字"}
	// ASSIGNMENT_TYPE_STAR 指针
	ASSIGNMENT_TYPE_STAR = AssignmentType{Code: "starExpr", Desc: "指针", ContentDesc: "关联变量的类型"}
	// ASSIGNMENT_TYPE_SELECTOR 层级调用
	ASSIGNMENT_TYPE_SELECTOR = AssignmentType{Code: "selectorExpr", Desc: "层级调用", ContentDesc: "关联调用的类型"}
	// ASSIGNMENT_TYPE_UNARY 一元表达式
	ASSIGNMENT_TYPE_UNARY = AssignmentType{Code: "unaryExpr", Desc: "一元表达式", ContentDesc: "一元的表达式"}
	// ASSIGNMENT_TYPE_BINARY 二元表达式
	ASSIGNMENT_TYPE_BINARY    = AssignmentType{Code: "binaryExpr", Desc: "二元表达式", ContentDesc: "AssignmentBinaryNode"}
	ASSIGNMENT_TYPE_SLICE     = AssignmentType{Code: "sliceExpr", Desc: "下标结构", ContentDesc: "AssignmentBinaryNode"}
	ASSIGNMENT_TYPE_VALUESPEC = AssignmentType{Code: "ValueSpec", Desc: "var类型", ContentDesc: "var类型"}
)

var AssignmentRightTypeMap = map[string]AssignmentType{
	ASSIGNMENT_TYPE_FUNCTION.Code:  ASSIGNMENT_TYPE_FUNCTION,
	ASSIGNMENT_TYPE_CALL.Code:      ASSIGNMENT_TYPE_CALL,
	ASSIGNMENT_TYPE_COMPOSITE.Code: ASSIGNMENT_TYPE_COMPOSITE,
	ASSIGNMENT_TYPE_BASICLIT.Code:  ASSIGNMENT_TYPE_BASICLIT,
	ASSIGNMENT_TYPE_EXPRSTMT.Code:  ASSIGNMENT_TYPE_EXPRSTMT,
	ASSIGNMENT_TYPE_IDENT.Code:     ASSIGNMENT_TYPE_IDENT,
	ASSIGNMENT_TYPE_STAR.Code:      ASSIGNMENT_TYPE_STAR,
	ASSIGNMENT_TYPE_SELECTOR.Code:  ASSIGNMENT_TYPE_SELECTOR,
	ASSIGNMENT_TYPE_UNARY.Code:     ASSIGNMENT_TYPE_UNARY,
	ASSIGNMENT_TYPE_BINARY.Code:    ASSIGNMENT_TYPE_BINARY,
}

func (at *AssignmentType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Code)
}

func (at *AssignmentType) UnmarshalJSON(data []byte) error {
	// 自定义序列化逻辑
	var code string
	err := json.Unmarshal(data, &code)
	if err != nil {
		return err
	}
	assignmentType, ok := AssignmentRightTypeMap[code]
	if ok {
		at.Code = assignmentType.Code
		at.Desc = assignmentType.Desc
		at.ContentDesc = assignmentType.ContentDesc
	}
	return nil
}
