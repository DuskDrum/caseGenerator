package faker

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type TimeFaker struct {
}

func (i TimeFaker) RangeFaker(min, max time.Time) time.Time {
	return timeRange(min, max)
}

func (i TimeFaker) RangeMinFaker(min time.Time) time.Time {
	return timeRange(min, time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC))
}

func (i TimeFaker) RangeMaxFaker(max time.Time) time.Time {
	return timeRange(time.Time{}, max)
}

func (i TimeFaker) FakerPtr(value time.Time) {
	// 不处理
}

func (i TimeFaker) FakerSizePtr(size int64, value time.Time) {
	// 不处理
}

func timeRange(min, max time.Time) time.Time {
	return gofakeit.DateRange(min, max)
}
