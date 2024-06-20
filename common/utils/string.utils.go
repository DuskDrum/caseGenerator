package utils

import "unicode"

// IsLower 判读字母是否是小写开头
func IsLower(s string) bool {
	for i, ch := range s {
		if i == 0 {
			// 判断首字母是否是小写
			return unicode.IsLower(ch)
		}
		// 如果不是首字母，返回false
		return false
	}
	return false
}
