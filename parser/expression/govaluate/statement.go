package govaluate

import (
	"caseGenerator/common/enum"
	_struct "caseGenerator/parser/struct"
	"go/token"
)

// StatementAssignment stmt的表达式，记录了参数的变动, 参数也可以直接重新赋值
type StatementAssignment struct {
	Name      string
	InitParam _struct.ValueAble
	Type      enum.StmtType
	// 参数变动列表
	ExpressDetail
}

// KeyFormula 赋值键值对
type KeyFormula struct {
	Key     string
	Formula string // 公式
	Type    enum.StmtType
	// 参数变动列表
	ExpressDetail
	//
	Position token.Position
}
