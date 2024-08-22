package enum

// ConditionExprType 总结下来，需要考虑：1. 用方法的返回值做逻辑运算符的条件 2. 布尔值做逻辑运算符 3. 某个类的属性做逻辑运算符 4. 多个条件做嵌套
type ConditionExprType struct {
	Name    string
	Desc    string
	Example string
}

var (
	CONDITION_EXPR_TYPE_COMPARE          = ConditionExprType{Name: "compare", Desc: "比较两个变量的值（相等、不相等、大于、小于、大于等于、小于等于）", Example: "x == y { }"}
	CONDITION_EXPR_TYPE_LOGIC            = ConditionExprType{Name: "logic", Desc: "使用逻辑运算符（与、或、非）连接多个条件", Example: "if a && b {}"}
	CONDITION_EXPR_TYPE_CALL             = ConditionExprType{Name: "call", Desc: "函数调用", Example: "strings.Contains(str, \"go\") { }"}
	CONDITION_EXPR_TYPE_GROUP            = ConditionExprType{Name: "group", Desc: "结合多个条件进行判断", Example: "if (x > 0 && y > 0) || z > 0 { }"}
	CONDITION_EXPR_TYPE_ZERO_VALUE_CHECK = ConditionExprType{Name: "zeroValueCheck", Desc: "检查变量是否为零值", Example: "if err != nil { }"}
	CONDITION_EXPR_TYPE_BOOL             = ConditionExprType{Name: "bool", Desc: "直接判断布尔变量的值", Example: "if isSuccess { }"}

	// CONDITION_EXPR_TYPE_TYPE_ASSERT 这种情况不考虑，解析成一个assignment(s, ok := i.(string)) 和一个condition(ok)
	CONDITION_EXPR_TYPE_TYPE_ASSERT = ConditionExprType{Name: "typeAssert", Desc: "检查接口变量的具体类型", Example: "if s, ok := i.(string); ok { }"}
	// CONDITION_EXPR_TYPE_CHAN 这种情况不考虑，解析成一个assignment(v, ok := <-ch) 和一个condition(ok)
	CONDITION_EXPR_TYPE_CHAN = ConditionExprType{Name: "chan", Desc: "使用通道的操作结果进行判断", Example: "if v, ok := <-ch; ok { }"}
	// CONDITION_EXPR_TYPE_ASSIGNMENT 这种情况不考虑，解析成一个assignment(x := f()) 和一个condition(x > 0)
	CONDITION_EXPR_TYPE_ASSIGNMENT = ConditionExprType{Name: "assignment", Desc: "在条件中进行赋值并立即检查", Example: "if x := f(); x > 0 { }"}
)

var ALL_CONDITION_EXPR_TYPE = map[string]ConditionExprType{
	CONDITION_EXPR_TYPE_COMPARE.Name:          CONDITION_EXPR_TYPE_COMPARE,
	CONDITION_EXPR_TYPE_LOGIC.Name:            CONDITION_EXPR_TYPE_LOGIC,
	CONDITION_EXPR_TYPE_CALL.Name:             CONDITION_EXPR_TYPE_CALL,
	CONDITION_EXPR_TYPE_TYPE_ASSERT.Name:      CONDITION_EXPR_TYPE_TYPE_ASSERT,
	CONDITION_EXPR_TYPE_GROUP.Name:            CONDITION_EXPR_TYPE_GROUP,
	CONDITION_EXPR_TYPE_ZERO_VALUE_CHECK.Name: CONDITION_EXPR_TYPE_ZERO_VALUE_CHECK,
	CONDITION_EXPR_TYPE_BOOL.Name:             CONDITION_EXPR_TYPE_BOOL,
	CONDITION_EXPR_TYPE_CHAN.Name:             CONDITION_EXPR_TYPE_CHAN,
	CONDITION_EXPR_TYPE_ASSIGNMENT.Name:       CONDITION_EXPR_TYPE_ASSIGNMENT,
}
