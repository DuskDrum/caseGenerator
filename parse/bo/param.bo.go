package bo

type ParamParseResult struct {
	// 参数名
	ParamName string
	// 参数类型
	ParamType string
	// 参数初始化值
	ParamInitValue string
	// 参数校验值
	ParamCheckValue string
	// 参数引入的依赖包
	//ImportPkgPath []string
	// 是否是省略号语法
	IsEllipsis bool
}

type Param interface {
	GetParamName() string
	UnmarshalerInfo(jsonString string)
}
