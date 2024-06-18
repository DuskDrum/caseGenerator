package parser

// SourceInfo 源信息
type SourceInfo struct {
	// 包信息
	Package
	// 依赖信息
	Import
	// 方法信息
	FunctionDeclare
	// 赋值信息列表
	AssignmentList []Assignment
	// 条件语句列表
	ConditionList []ConditionNode
}
