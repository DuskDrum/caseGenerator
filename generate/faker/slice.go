package faker

import "github.com/brianvoe/gofakeit/v6"

type SliceFaker struct {
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

func (i SliceFaker) FakerPtr(value any) {
	gofakeit.Slice(value)
}

func (i SliceFaker) FakerSizePtr(size int64, value any) {
	faker := gofakeit.New(size)
	faker.Slice(value)
}
