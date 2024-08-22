package generate

// Request 记录方法要请求的信息
type Request struct {
	// receiver相关
	ReceiverName string
	ReceiverType string
	// 参数
	ParamList []Param
}

type Param struct {
	ParamDetail ParamDetail
}
type ParamDetail struct {
	ParamName string
	ParamType string
	//递归子类型
	ChildParamList *ParamDetail
}

// 记录receiver，方法名，请求参数
