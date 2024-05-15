package example

// RequestGenericProblem 泛型
func RequestGenericProblem[T, R any](list []T, process func([]T) []R, batchSize int) {

}

// RequestGeneric1Problem 泛型
func RequestGeneric1Problem[T, R any](list []T, process func([]T) ([]R, string), batchSize int) {

}

// RequestGeneric3Problem 泛型
func RequestGeneric3Problem[T, R any](list []T, process func([]T) (result []R, str string), batchSize int) {

}

//func RequestResponseGenericValueProblem[T int | uint | int8 | int16 | int32 | int64 | float32 | float64 | string | bool](func() (p *T, s T)) {
//}
