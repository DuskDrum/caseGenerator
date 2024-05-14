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
	// ImportList 依赖信息
	ImportList     []string
	AliasImportMap map[string]string
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

func GetImportInfo() []string {
	mu.RLock()
	defer mu.RUnlock()
	return ImportList
}

func InitImport() {
	mu.Lock()
	defer mu.Unlock()
	ImportList = make([]string, 0, 10)
	AliasImportMap = make(map[string]string, 10)
}

func AppendImportList(item string) {
	mu.Lock()
	defer mu.Unlock()
	ImportList = append(ImportList, item)
}

func AppendAliasImport(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	AliasImportMap[key] = value
}

func GetImportPathFromAliasMap(name string) string {
	s, ok := AliasImportMap[name]
	if ok {
		result := strings.ReplaceAll(s, "\"", "")
		result = name + " \"" + result + "\""
		return result
	}
	for _, info := range ImportList {
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
	ParamNeedToMap = sync.Map{}
	receiverInfo = nil
	requestDetailList = make([]generate.RequestDetail, 0, 10)
}
