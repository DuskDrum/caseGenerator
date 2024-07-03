package enum

import "encoding/json"

// ConditionInitType 条件初始化类型
type ConditionInitType struct {
	Code        string
	Desc        string
	ContentDesc string
}

var (
	// CONDITION_INIT_TYPE_ASSIGNMENT 赋值初始化
	CONDITION_INIT_TYPE_ASSIGNMENT = ConditionInitType{Code: "assignment", Desc: "赋值初始化", ContentDesc: "string"}
	// CONDITION_INIT_TYPE_COMPOSITE 构造类型，比如{1,2,3,4}
	CONDITION_INIT_TYPE_COMPOSITE = ConditionInitType{Code: "composite", Desc: "构造类型", ContentDesc: "CompositeLitValue"}
)

var ConditionInitTypeMap = map[string]ConditionInitType{
	CONDITION_INIT_TYPE_ASSIGNMENT.Code: CONDITION_INIT_TYPE_ASSIGNMENT,
	CONDITION_INIT_TYPE_COMPOSITE.Code:  CONDITION_INIT_TYPE_COMPOSITE,
}

func (at *ConditionInitType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Code)
}

func (at *ConditionInitType) UnmarshalJSON(data []byte) error {
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
