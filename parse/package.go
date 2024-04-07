package parse

import "strings"

// Package 包信息
type Package struct {
	PackagePath string
	PackageName string
}

// Import 依赖信息
type Import struct {
	ImportList     []string
	AliasImportMap map[string]string
}

// Method 方法
type Method struct {
	MethodName string
}

// GetImportPath  获取import信息
func (i Import) GetImportPath(name string) string {
	s, ok := i.AliasImportMap[name]
	if ok {
		result := strings.ReplaceAll(s, "\"", "")
		result = name + " \"" + result + "\""
		return result
	}
	for _, info := range i.ImportList {
		xx := strings.ReplaceAll(info, "\"", "")
		if strings.HasSuffix(xx, name) {
			result := "\"" + xx + "\""
			return result
		}
	}
	return name
}
