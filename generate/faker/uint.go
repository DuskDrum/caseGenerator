package faker

import (
	"math"

	"github.com/brianvoe/gofakeit/v6"
)

type UintFaker struct {
}

func (i UintFaker) RangeFaker(min, max uint) uint {
	return uintRange(min, max)
}

func (i UintFaker) RangeMinFaker(min uint) uint {
	return uintRange(min, math.MaxInt)
}

func (i UintFaker) RangeMaxFaker(max uint) uint {
	return uintRange(math.MinInt, max)
}

func (i UintFaker) FakerPtr(value uint) {
	// 不处理
}

func (i UintFaker) FakerSizePtr(size int64, value uint) {
	// 不处理
}

func uintRange(min, max uint) uint {
	return gofakeit.UintRange(min, max)
}
