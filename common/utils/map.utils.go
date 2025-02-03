package utils

func PutAll[T any](sourceMap map[string]T, targetMap map[string]T) {
	for k, v := range targetMap {
		sourceMap[k] = v
	}
}

func CopyMap[T any](sourceMap map[string]T) map[string]T {
	targetMap := make(map[string]T, len(sourceMap))
	for k, v := range sourceMap {
		targetMap[k] = v
	}
	return targetMap
}
