package utils

// LoEmpty 快速得到任何类型的0值
func LoEmpty(typeStr string) string {
	return "lo.Empty[" + typeStr + "]"
}
