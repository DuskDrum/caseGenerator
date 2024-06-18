package parser

// Param 参数信息，请求参数、返回参数、类型断言参数、变量、常量等
type Param struct {
	Name    string
	Type    string
	AstType string
}
