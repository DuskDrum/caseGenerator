package faker

import (
	"math"

	"github.com/brianvoe/gofakeit/v6"
)

type Float64Faker struct {
}

func (i Float64Faker) RangeFaker(min, max float64) float64 {
	return float64Range(min, max)
}

func (i Float64Faker) RangeMinFaker(min float64) float64 {
	return float64Range(min, math.MaxInt)
}

func (i Float64Faker) RangeMaxFaker(max float64) float64 {
	return float64Range(math.MinInt, max)
}

func (i Float64Faker) FakerPtr(value float64) {
	// 不处理
}

func (i Float64Faker) FakerSizePtr(size int64, value float32) {
	// 不处理
}

func float64Range(min, max float64) float64 {
	return gofakeit.Float64Range(min, max)
}
