package test

import (
	"caseGenerator/parser/expression/z3/test/common"
	"math"
	"math/rand"
)

func Selector() int {
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

func RequestSelector(ZeroInt string) int {
	if ZeroInt == "" {
		return 0
	}
	return 1
}

func GlobalSelector1() int {
	localVariable := GlobalReceiverName
	if localVariable == "" {
		return 0
	}
	return 1
}

func GlobalSelector2() int {
	localVariable := GlobalReceiverAge

	return localVariable
}

func GlobalSelector3() int {
	localVariable := GlobalReceiverIsPass
	if localVariable {
		return 1
	}
	return -1
}
