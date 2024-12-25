package bo

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/expression"
	_struct "caseGenerator/parser/struct"
)

// StatementAssignment stmt的表达式，记录了参数的变动, 参数也可以直接重新赋值
type StatementAssignment struct {
	Name      string
	InitParam _struct.ValueAble
	Type      enum.StmtType
	// 参数变动列表
	expression.ExpressDetail
}

type KeyFormula struct {
	Key     string
	Formula string // 公式
	Type    enum.StmtType
	// 参数变动列表
	expression.ExpressDetail
}
