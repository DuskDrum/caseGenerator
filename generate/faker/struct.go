package faker

import "github.com/brianvoe/gofakeit/v6"

type StructFaker struct {
}

//func (i StructFaker) RangeFaker(min, max uint) uint {
//	return uintRange(min, max)
//}
//
//func (i StructFaker) RangeMinFaker(min uint) uint {
//	return uintRange(min, math.MaxInt)
//}
//
//func (i StructFaker) RangeMaxFaker(max uint) uint {
//	return uintRange(math.MinInt, max)
//}

func (i StructFaker) FakerPtr(value any) {
	err := gofakeit.Struct(value)
	if err != nil {
		panic("catch err, detail is: " + err.Error())
	}
}

func (i StructFaker) FakerSizePtr(size int64, value any) {
	faker := gofakeit.New(size)
	err := faker.Struct(value)
	if err != nil {
		panic("catch err, detail is: " + err.Error())
	}
}
