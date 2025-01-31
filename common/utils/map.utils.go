package utils

func PutAll[T any](sourceMap map[string]T, targetMap map[string]T) {
	for k, v := range targetMap {
		sourceMap[k] = v
	}
}
