package bo

import (
	"caseGenerator/generate"
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
	// ParamNeedToMap 信息
	ParamNeedToMap sync.Map
	// receiverInfo receiver信息
	receiverInfo *ReceiverInfo
	// requestDetailList 请求详情列表
	requestDetailList []generate.RequestDetail
)

func AppendRequestDetailToList(gr generate.RequestDetail) {
	mu.Lock()
	defer mu.Unlock()
	if len(requestDetailList) == 0 {
		requestDetailList = make([]generate.RequestDetail, 0, 10)
	}
	requestDetailList = append(requestDetailList, gr)
}

func GetRequestDetailList() []generate.RequestDetail {
	mu.RLock()
	defer mu.RUnlock()
	return requestDetailList
}

func SetTypeParamMap(typeParamMap map[string]*ParamParseResult) {
	mu.Lock()
	defer mu.Unlock()
	TypeParamMap = typeParamMap
}

func GetTypeParamMap() map[string]*ParamParseResult {
	mu.RLock()
	defer mu.RUnlock()
	if TypeParamMap == nil {
		return make(map[string]*ParamParseResult, 10)
	} else {
		return TypeParamMap
	}
}

// Import 依赖信息
type Import struct {
	ImportList     []string
	AliasImportMap map[string]string
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

func AddParamNeedToMapDetail(paramName string, p Param) {
	ParamNeedToMap.Store(paramName, p)
}

func SetReceiverInfo(r *ReceiverInfo) {
	receiverInfo = r
}

func GetReceiverInfo() *ReceiverInfo {
	return receiverInfo
}

func ClearBo() {
	mu.Lock()
	defer mu.Unlock()
	TypeParamMap = make(map[string]*ParamParseResult, 10)
	ImportInfo = Import{
		ImportList:     make([]string, 0, 10),
		AliasImportMap: make(map[string]string, 10),
	}
	ParamNeedToMap = sync.Map{}
	receiverInfo = nil
	requestDetailList = make([]generate.RequestDetail, 0, 10)
}
