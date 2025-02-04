package utils

func CopySlice[T any](sourceSlice []T) []T {
	targetSlice := make([]T, 0, len(sourceSlice))
	for _, v := range sourceSlice {
		targetSlice = append(targetSlice, v)
	}
	return targetSlice
}
