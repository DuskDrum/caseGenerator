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

func FakeNeqTarget(targetValue string) string {
	gofakeit.RandomString()
}
