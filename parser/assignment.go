package parser

type Assignment struct {
	Param
	// 赋值的位置
	ParamIndex int
	// 赋值的目标值的生成策略
	TargetValueStrategy string
}
