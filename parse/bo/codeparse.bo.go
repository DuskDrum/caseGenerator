package bo

import (
	"caseGenerator/generate"
	"caseGenerator/generate_old"
	"sync"
)

var (
	// 读写锁
	mu sync.RWMutex
	// TypeParamMap 用来存储TypeParam
	TypeParamMap map[string]*ParamParseResult
	// ImportList 依赖信息
	initImportList     []string
	initAliasImportMap map[string]string

	ImportList []string
	// ParamNeedToMap 信息
	ParamNeedToMap sync.Map
	// receiverInfo receiver信息
	receiverInfo *ReceiverInfo
	// requestDetailList 请求详情列表
	requestDetailList []generate.CaseRequest
	// assignmentInfoList  赋值list
	assignmentInfoList []AssignmentDetailInfo
	// mockInfoList mock数据列表
	mockInfoList []*generate_old.MockInstruct
)

func AppendMockInfoList(item generate_old.MockInstruct) {
	mu.Lock()
	defer mu.Unlock()
	if mockInfoList == nil {
		mockInfoList = make([]*generate_old.MockInstruct, 0, 10)
	}
	mockInfoList = append(mockInfoList, &item)
}

func GetMockInfoList() []*generate_old.MockInstruct {
	return mockInfoList
}

func AppendAssignmentDetailInfoToList(info AssignmentDetailInfo) {
	mu.Lock()
	defer mu.Unlock()
	if len(assignmentInfoList) == 0 {
		assignmentInfoList = make([]AssignmentDetailInfo, 0, 10)
	}
	assignmentInfoList = append(assignmentInfoList, info)
}

func GetAssignmentDetailInfoList() []AssignmentDetailInfo {
	return assignmentInfoList
}

func AppendRequestDetailToList(gr generate.CaseRequest) {
	mu.Lock()
	defer mu.Unlock()
	if len(requestDetailList) == 0 {
		requestDetailList = make([]generate.CaseRequest, 0, 10)
	}
	requestDetailList = append(requestDetailList, gr)
}

func GetRequestDetailList() []generate.CaseRequest {
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
	requestDetailList = make([]generate.CaseRequest, 0, 10)
}
