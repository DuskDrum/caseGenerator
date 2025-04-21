package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetRelativePath(filename string) string {
	// 向上寻找 go.mod
	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			dir = filepath.Dir(dir)
			fmt.Println("项目根目录:", dir)
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Println("未找到 go.mod，可能不在 Go 项目中")
			break
		}
		dir = parent
	}

	dir += "/"

	// 用 map 存储 B 中的字符，方便快速判断
	removeSet := make(map[rune]bool)
	for _, ch := range dir {
		removeSet[ch] = true
	}

	// 遍历 A，跳过出现在 B 中的字符
	var builder strings.Builder
	for _, ch := range filename {
		if !removeSet[ch] {
			builder.WriteRune(ch)
		}
	}
	return builder.String()
}
