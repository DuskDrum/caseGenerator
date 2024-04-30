package parse

import (
	"strings"
	"sync"
)

var (
	// 读写锁
	mu sync.RWMutex
	// TypeParamMap 用来存储TypeParam
	TypeParamMap map[string]*ParamParseResult
	// ImportInfo 依赖信息
	ImportInfo Import
)

func SetTypeParamMap(typeParamMap map[string]*ParamParseResult) {
	mu.Lock()
	defer mu.Unlock()
	TypeParamMap = typeParamMap
}

func GetTypeParamMap() map[string]*ParamParseResult {
	mu.RLock()
	defer mu.RUnlock()
	return TypeParamMap
}

func GetImportInfo() Import {
	mu.RLock()
	defer mu.RUnlock()
	return ImportInfo
}

func InitImport() {
	mu.Lock()
	defer mu.Unlock()
	importList := make([]string, 0, 10)
	importMap := make(map[string]string, 10)
	importInfo := Import{
		ImportList:     importList,
		AliasImportMap: importMap,
	}
	ImportInfo = importInfo
}

func AppendImportList(item string) {
	mu.Lock()
	defer mu.Unlock()
	ImportInfo.ImportList = append(ImportInfo.ImportList, item)
}

func AppendAliasImport(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	ImportInfo.AliasImportMap[key] = value
}

// Import 依赖信息
type Import struct {
	ImportList     []string
	AliasImportMap map[string]string
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
