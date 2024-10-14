package enum

import (
	"encoding/json"
)

type ParameterType BaseEnum

var (
	PARAMETER_TYPE_SELECTOR     = ParameterType{Name: "Selector", Desc: "选择结构,类似于a.b的结构"}
	PARAMETER_TYPE_IDENT        = ParameterType{Name: "Ident", Desc: "比如报名,函数名,变量名"}
	PARAMETER_TYPE_STAR         = ParameterType{Name: "Star", Desc: "指针类型"}
	PARAMETER_TYPE_FUNC         = ParameterType{Name: "Func", Desc: "函数类型"}
	PARAMETER_TYPE_INTERFACE    = ParameterType{Name: "Interface", Desc: "接口类型"}
	PARAMETER_TYPE_ARRAY        = ParameterType{Name: "Array", Desc: "数组类型"}
	PARAMETER_TYPE_MAP          = ParameterType{Name: "Map", Desc: "Map 类型"}
	PARAMETER_TYPE_ELLIPSIS     = ParameterType{Name: "Ellipsis", Desc: "省略号表达式,比如参数列表的最后一个可以写成arg..."}
	PARAMETER_TYPE_CHAN         = ParameterType{Name: "Chan", Desc: "管道类型"}
	PARAMETER_TYPE_INDEX        = ParameterType{Name: "Index", Desc: "下标类型，expr[expr]"}
	PARAMETER_TYPE_BASICLIT     = ParameterType{Name: "BasicLit", Desc: "基本字面值,数字或者字符串"}
	PARAMETER_TYPE_FUNCLIT      = ParameterType{Name: "FuncLit", Desc: "函数定义"}
	PARAMETER_TYPE_COMPOSITELIT = ParameterType{Name: "CompositeLit", Desc: "构造类型,比如{1,2,3,4}"}
	PARAMETER_TYPE_CALL         = ParameterType{Name: "Call", Desc: "expr()"}
	PARAMETER_TYPE_KEYVALUE     = ParameterType{Name: "KeyValue", Desc: "key:value"}
	PARAMETER_TYPE_UNARY        = ParameterType{Name: "Unary", Desc: "一元表达式"}
	PARAMETER_TYPE_BINARY       = ParameterType{Name: "Binary", Desc: "二元表达式"}
	PARAMETER_TYPE_PARENT       = ParameterType{Name: "Parent", Desc: "括号表达式"}
)

var ALL_PARAMETER_TYPE = map[string]ParameterType{
	PARAMETER_TYPE_SELECTOR.Name:     PARAMETER_TYPE_SELECTOR,
	PARAMETER_TYPE_IDENT.Name:        PARAMETER_TYPE_IDENT,
	PARAMETER_TYPE_STAR.Name:         PARAMETER_TYPE_STAR,
	PARAMETER_TYPE_FUNC.Name:         PARAMETER_TYPE_FUNC,
	PARAMETER_TYPE_INTERFACE.Name:    PARAMETER_TYPE_INTERFACE,
	PARAMETER_TYPE_ARRAY.Name:        PARAMETER_TYPE_ARRAY,
	PARAMETER_TYPE_MAP.Name:          PARAMETER_TYPE_MAP,
	PARAMETER_TYPE_ELLIPSIS.Name:     PARAMETER_TYPE_ELLIPSIS,
	PARAMETER_TYPE_CHAN.Name:         PARAMETER_TYPE_CHAN,
	PARAMETER_TYPE_INDEX.Name:        PARAMETER_TYPE_INDEX,
	PARAMETER_TYPE_BASICLIT.Name:     PARAMETER_TYPE_BASICLIT,
	PARAMETER_TYPE_FUNCLIT.Name:      PARAMETER_TYPE_FUNCLIT,
	PARAMETER_TYPE_COMPOSITELIT.Name: PARAMETER_TYPE_COMPOSITELIT,
	PARAMETER_TYPE_CALL.Name:         PARAMETER_TYPE_CALL,
	PARAMETER_TYPE_KEYVALUE.Name:     PARAMETER_TYPE_KEYVALUE,
	PARAMETER_TYPE_UNARY.Name:        PARAMETER_TYPE_UNARY,
	PARAMETER_TYPE_BINARY.Name:       PARAMETER_TYPE_BINARY,
	PARAMETER_TYPE_PARENT.Name:       PARAMETER_TYPE_PARENT,
}

func (at *ParameterType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Name)
}

func (at *ParameterType) UnmarshalJSON(data []byte) error {
	// 自定义序列化逻辑
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	assignmentType, ok := ALL_PARAMETER_TYPE[name]
	if ok {
		at.Name = assignmentType.Name
		at.Desc = assignmentType.Desc
	}
	return nil
}
