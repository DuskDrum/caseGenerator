package test

import (
	"caseGenerator/parser/expression/z3/test/common"
	"math"
	"math/rand"
)

var ZeroInt = 0

func Ident() int {
	i := rand.Int()
	maxInt := math.MaxInt
	if i < ZeroInt {
		return -1
	}
	if i > maxInt {
		return -1
	}
	if i > common.TEST_INIT {
		return -1
	}
	return 1
}
