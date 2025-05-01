package bo

// ExpressionContext  expression上下文
type ExpressionContext struct {
	VariableParamMap  map[string]string // 解析出变量map
	RequestParamMap   map[string]any    // 方法请求参数信息
	TemporaryVariable TemporaryVariable // 临时参数
}

type TemporaryVariable struct {
	VariableIndex int    // 变量的下标，指的是赋值语句中比如一个call返回两个参数，那么这里要标识此参数是第几个参数，默认是0
	VariableName  string // 参数名
}
