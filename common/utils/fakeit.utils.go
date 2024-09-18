package utils

import (
	"go/token"

	"github.com/brianvoe/gofakeit/v6"
)

// FakeItByOp 根据Op类型Fake值
//
//	EQL    // ==
//	NEQ      // !=
//	LEQ      // <=
//	GEQ      // >=
//	LSS    // <
//	GTR    // >
func FakeItByOp(op token.Token, targetValue string) string {
	if op == token.EQL {
		return targetValue
	} else if op == token.NEQ {
		return ""
	} else if op == token.LEQ {
		return ""
	} else if op == token.GEQ {
		return ""
	} else {
		return ""
	}

}

// 不等于某个值的fake
// <= 某个值的fake
// >= 某个值的fake
// > 某个值 的fake
// < 某个值的fake
func FakeNeqTarget(targetValue string) string {
	faker := gofakeit.New(int64(len(targetValue) + 1))
	return faker.Letter()
}
