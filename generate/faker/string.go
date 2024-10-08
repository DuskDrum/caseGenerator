package faker

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
)

type StringFaker struct {
}

//func (i SliceFaker) RangeFaker(min, max uint) uint {
//	return uintRange(min, max)
//}
//
//func (i SliceFaker) RangeMinFaker(min uint) uint {
//	return uintRange(min, math.MaxInt)
//}
//
//func (i SliceFaker) RangeMaxFaker(max uint) uint {
//	return uintRange(math.MinInt, max)
//}

func (i StringFaker) FakerPtr(value *string) {
	value = lo.ToPtr(gofakeit.Letter())
}

func (i StringFaker) FakerSizePtr(size int64, value *string) {
	value = lo.ToPtr(gofakeit.LetterN(uint(size)))

}
