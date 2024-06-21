package parser

// FunctionBasic 函数基础信息， 方法名、路径名、泛型的映射、receiver接受者信息、文件名
type FunctionBasic struct {
	FunctionPath string
	FileName     string
	GenericsMap  map[string]Param
	Receiver     Param
}

// FunctionDeclare 函数的声明，请求列表，响应列表，也包含了基础信息
type FunctionDeclare struct {
	FunctionBasic
	FunctionName string
	RequestList  []*Param
	ResponseList []*Param
}
