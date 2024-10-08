package faker

import (
	"math"

	"github.com/brianvoe/gofakeit/v6"
)

type Float32Faker struct {
}

func (i Float32Faker) RangeFaker(min, max float32) float32 {
	return float32Range(min, max)
}

func (i Float32Faker) RangeMinFaker(min float32) float32 {
	return float32Range(min, math.MaxInt)
}

func (i Float32Faker) RangeMaxFaker(max float32) float32 {
	return float32Range(math.MinInt, max)
}

func (i Float32Faker) FakerPtr(value float32) {
	// 不处理
}

func (i Float32Faker) FakerSizePtr(size int64, value float32) {
	// 不处理
}

func float32Range(min, max float32) float32 {
	return gofakeit.Float32Range(min, max)
}
