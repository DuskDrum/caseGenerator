package example

// RequestGenericProblem 泛型
func RequestGenericProblem[T, R any](list []T, process func([]T) []R, batchSize int) {

}
