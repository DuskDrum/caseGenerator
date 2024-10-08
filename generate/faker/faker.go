package faker

type RangeFaker[T any] interface {
	// RangeFaker 区域型的faker
	RangeFaker(min, max T) T
	// RangeMinFaker 指定最小的faker
	RangeMinFaker(min T) T
	// RangeMaxFaker 指定最大的faker
	RangeMaxFaker(max T) T
}

type PtrFaker[T any] interface {
	// FakerPtr 指针faker
	FakerPtr(value T)
	// FakerSizePtr 指定大小的指针faker
	FakerSizePtr(size int64, value T)
}

type Faker[T any] interface {
	RangeFaker[T]
	PtrFaker[T]
}
