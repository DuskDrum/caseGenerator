package faker

import (
	"math"

	"github.com/brianvoe/gofakeit/v6"
)

type IntFaker struct {
}

func (i IntFaker) RangeFaker(min, max int) int {
	return intRange(min, max)
}

func (i IntFaker) RangeMinFaker(min int) int {
	return intRange(min, math.MaxInt)
}

func (i IntFaker) RangeMaxFaker(max int) int {
	return intRange(math.MinInt, max)
}

func (i IntFaker) FakerPtr(value int) {
	// 不处理
}

func (i IntFaker) FakerSizePtr(size int64, value int) {
	// 不处理
}

// intRange intRange
func intRange(min, max int) int {
	// 调整整数范围
	randomNumber := gofakeit.Number(min, max)
	return randomNumber
}
