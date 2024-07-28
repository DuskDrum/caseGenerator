package generate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser"
	"strconv"
	"strings"
)

// StandardInfo 标准信息
type StandardInfo struct {
	// 测试用例的包名
	PackageName string
	// 生成的文件名
	FileName string
	// 生成的文件路径
	FilePath string
	// 需要依赖的import
	ImportPkgPaths []string
	// case列表
	MethodDetailList []FunctionCase
	// goLink 列表
	GoLinkList []string
}

// FunctionCase 方法层面的case，一个方法可能要根据条件、入参出参、mock方法生成多次
type FunctionCase struct {
	// 方法名
	FunctionName string
	// receive类型
	ReceiverType string
	// case名称
	CaseName string
	// mock列表
	MockList []*CaseMockInfo
	// 请求列表
	RequestList []CaseRequest
	// 把请求名按照tt.args.xx，并按照","分割，最后一位没有逗号
	RequestNameString string
}

type CaseMockInfo struct {
	MockResponseParam []string
	MockFunction      string
	MockNumber        string
	// mock 返回
	MockReturns        string
	MockFunctionParam  []MockParamInfo
	MockFunctionResult []MockParamInfo
}

type MockParamInfo struct {
	// 参数名
	ParamName string
	// 参数类型
	ParamType string
	// 参数初始化值
	ParamInitValue string
	// 参数校验值
	ParamCheckValue string
	// 是否是省略号语法
	IsEllipsis bool
}

// CaseLinkedInfo 内部方法
type CaseLinkedInfo struct {
	// 请求信息
	InnerFuncRequest []CaseRequest
	// 响应信息
	InnerFuncResponse []CaseResponse
	// 代名
	InnerFuncMethodName string
	// 代名关联的方法路径
	InnerFuncLinkedName string
}

type CaseResponse struct {
	// 参数名
	ParamName string
	// 参数类型
	ParamType string
	// 参数初始化值
	ParamInitValue string
	// 参数校验值
	ParamCheckValue string
	// 是否是省略号语法
	IsEllipsis bool
}

func (rd CaseResponse) GenerateResponseContent() string {
	// 1. 使用 stringBuilder 解析请求名称 + 请求类型
	var stringBuilder strings.Builder
	stringBuilder.WriteString(rd.ParamName)
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(rd.ParamType)
	return stringBuilder.String()
}

type CaseRequest struct {
	// 请求名称
	RequestName string
	// 请求类型
	RequestType string
	// 请求值
	RequestValue string
	// 是否是省略号语法
	IsEllipsis bool
}

func (rd CaseRequest) GenerateRequestContent() string {
	// 1. 使用 stringBuilder 解析请求名称 + 请求类型
	var stringBuilder strings.Builder
	stringBuilder.WriteString(rd.RequestName)
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(rd.RequestType)
	return stringBuilder.String()
}

// GenStandardInfo 将sourceInfo 源信息转为生成模板的标准信息
func GenStandardInfo(sourceInfo *parser.SourceInfo) []*StandardInfo {
	list := sourceInfo.FileInfoList
	for i, info := range list {
		var standardInfo StandardInfo
		// 1. 组装基础信息
		standardInfo.FilePath = info.FunctionPath
		standardInfo.FileName = info.FileName
		standardInfo.PackageName = sourceInfo.PackageName
		// 2. 组装方法信息
		methodList := make([]FunctionCase, 0, 10)
		// 3. 组装goLink 信息
		goLinkList := make([]string, 0, 10)
		// 4. 组装import信息
		importList := make([]string, 0, 10)

	}

	return nil
}

// GenGoLink 根据解析出来的方法信息，生成go-linked记录, 组装go:link
// //go:linkname awxCommonConvertSettlementReportAlgorithm slp/reconcile/core/message/standard.commonConvertSettlementReportAlgorithm
// func awxCommonConvertSettlementReportAlgorithm(transactionType enums.TransactionType, createdAt time.Time, ctx context.Context, dataBody dto.AwxSettlementReportDataBody) (result []service.OrderAlgorithmResult, err error)
func GenGoLink(fileInfo *parser.FileInfo) []string {
	// 1. 使用 stringBuilder 解析go linkname
	var stringBuilder strings.Builder
	stringBuilder.WriteString("//go:linkname ")
	stringBuilder.WriteString(fp.aliasFuncName)
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(fp.moduleName)
	stringBuilder.WriteString("/")
	stringBuilder.WriteString(fp.filepath)
	stringBuilder.WriteString(".")
	stringBuilder.WriteString(fp.funcName)
	stringBuilder.WriteString("\n")
	// 2. 组装内部方法对应的结构
	stringBuilder.WriteString("func ")
	stringBuilder.WriteString(fp.aliasFuncName)
	stringBuilder.WriteString("(")
	// 2.1 解析内部方法的请求
	requestStr := ""
	for _, v := range fp.requestList {
		content := v.GenerateRequestContent()
		requestStr = requestStr + content + ", "
	}
	if len(requestStr) > 0 {
		requestStr = requestStr[0 : len(requestStr)-2]
	}
	stringBuilder.WriteString(requestStr)
	stringBuilder.WriteString(")")
	// 2.2 解析内部方法的响应
	if len(fp.responseList) > 0 {
		stringBuilder.WriteString("(")
		responseStr := ""
		for _, v := range fp.responseList {
			content := v.GenerateResponseContent()
			responseStr = responseStr + content + ", "
		}
		responseStr = responseStr[0 : len(responseStr)-2]
		stringBuilder.WriteString(responseStr)
		stringBuilder.WriteString(")")
	}

	return stringBuilder.String()
}

// GenMethodInfo 方法组装
func GenMethodInfo(fileInfo *parser.FileInfo, i int) []FunctionCase {
	var mc FunctionCase
	mockList := make([]CaseMockInfo, 0, 10)
	requestList := make([]CaseRequest, 0, 10)
	// 组装方法的信息
	mc.FunctionName = fileInfo.FunctionName
	mc.CaseName = fileInfo.FunctionName + strconv.Itoa(i)
	if fileInfo.Receiver == nil {
		mc.ReceiverType = fileInfo.Receiver.Type
	}
	if len(fileInfo.RequestList) > 0 {
		var requestNameString string
		for _, r := range fileInfo.RequestList {
			if r.AstType == enum.PARAM_AST_TYPE_Ellipsis {
				requestNameString += "tt.args." + r.Name + "... , "
			} else {
				requestNameString += "tt.args." + r.Name + ", "
			}
		}
		requestNameString = strings.TrimRight(requestNameString, ", ")
		mc.RequestNameString = requestNameString
	}
	// 组装mock信息

}
