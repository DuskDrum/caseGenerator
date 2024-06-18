package parser

import (
	"go/ast"
	"strings"
	"sync"
)

// Import 依赖信息
type Import struct {
	mu sync.RWMutex

	ImportList     []string
	AliasImportMap map[string]string
}

// InitImport 初始化 import信息
func (i *Import) InitImport(af *ast.File) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ImportList = make([]string, 0, 10)
	i.AliasImportMap = make(map[string]string, 10)
	for _, importSpec := range af.Imports {
		if importSpec.Name == nil {
			i.AppendImportList(importSpec.Path.Value)
		} else {
			i.AppendAliasImport(importSpec.Name.Name, importSpec.Path.Value)
		}
	}
}

// AppendImportList 增加依赖列表
func (i *Import) AppendImportList(item string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.ImportList = append(i.ImportList, item)
}

// AppendAliasImport 增加依赖映射
func (i *Import) AppendAliasImport(key string, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.AliasImportMap[key] = value
}

// GetImportPathFromAliasMap 根据名称
func (i *Import) GetImportPathFromAliasMap(name string) string {
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
