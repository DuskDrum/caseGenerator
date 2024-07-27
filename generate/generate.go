package generate

import (
	"bytes"
	"caseGenerator/common/utils"
	"caseGenerator/parser"
	template2 "caseGenerator/template"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/samber/lo"
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
	CaseDetailList []MethodCase
	// mock列表
	MockList []*MockInstruct
	// goLink 列表
	GoLinkList []string
}

// MethodCase 方法层面的case，一个方法可能要根据条件、入参出参、mock方法生成多次
type MethodCase struct {
	// 方法名
	MethodName string
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
	// goLink 列表
	GoLinkList []string
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
func GenStandardInfo(sourceInfo *parser.SourceInfo) *StandardInfo {
	return nil
}

func GenGenerateFile(data GenMeta) {
	var buf bytes.Buffer

	caseDetails := data.CaseDetailList
	// 设置RequestNameString字段
	cdList := make([]CaseDetail, 0, 10)
	importList := make([]string, 0, 10)
	importList = append(importList, "\"testing\"")

	for _, cd := range caseDetails {
		// 如果是内部方法，那么跳过处理
		// 将字符串转换为rune类型

		cd.FileName = strings.ReplaceAll(data.FileName, ".", "")
		if len(cd.RequestList) > 0 {
			var requestNameString string
			for _, r := range cd.RequestList {
				if r.IsEllipsis {
					requestNameString += "tt.args." + r.RequestName + "... , "
				} else {
					requestNameString += "tt.args." + r.RequestName + ", "
				}
			}
			requestNameString = strings.TrimRight(requestNameString, ", ")
			cd.RequestNameString = requestNameString
		}
		cdList = append(cdList, cd)
	}
	data.CaseDetailList = cdList

	data.MockList = lo.Filter(data.MockList, func(item *MockInstruct, index int) bool {
		if item.MockFunction == "" {
			return false
		} else {
			return true
		}
	})
	mockInstructs := make([]*MockInstruct, 0, 10)

	// mockey.Mock((*repo.ClearingPipeConfigRepo).GetAllConfigs).Return(clearingPipeConfigs).Build()
	// go:linkname awxCommonConvertSettlementReportAlgorithm slp/reconcile/core/message/standard.commonConvertSettlementReportAlgorithm
	// func awxCommonConvertSettlementReportAlgorithm(transactionType enums.TransactionType, createdAt time.Time, ctx context.Context, dataBody dto.AwxSettlementReportDataBody) (result []service.OrderAlgorithmResult, err error)
	if len(data.MockList) > 0 {
		by := lo.SliceToMap(data.MockList, func(item *MockInstruct) (string, *MockInstruct) {
			return item.MockFunction, item
		})
		for k, v := range by {
			// 1. 组装mock的响应值
			str := "[]any{"
			for range v.MockResponseParam {
				str += " nil,"
			}
			// 去掉最后一个逗号
			lastCommaIndex := strings.LastIndex(str, ",")
			if lastCommaIndex != -1 {
				str = str[:lastCommaIndex]
			}
			str += "}"
			mockReturns := str
			// 2. 组装mock的返回
			mockNumber := "mock" + k
			// 3. 如果方法名是小写开头，且没有包名引用，说明需要go-linkname
			if utils.IsLower(v.MockFunction) && !strings.Contains(v.MockFunction, ".") {

			}

			mi := MockInstruct{
				MockResponseParam:  v.MockResponseParam,
				MockFunction:       v.MockFunction,
				MockNumber:         mockNumber,
				MockReturns:        mockReturns,
				MockFunctionParam:  v.MockFunctionParam,
				MockFunctionResult: v.MockFunctionResult,
			}

			mockInstructs = append(mockInstructs, &mi)
		}
		data.MockList = mockInstructs
	}

	if len(data.ImportPkgPaths) > 0 {
		for _, v := range data.ImportPkgPaths {
			if v != "" {
				importList = append(importList, v)
			}
		}
	}
	uniqImportList := lo.Uniq(importList)
	data.ImportPkgPaths = uniqImportList

	err := genRender(template2.NotHaveReceiveModel, &buf, data)
	if err != nil {
		_ = fmt.Errorf("cannot format file: %w", err)
	}
	modelFile := data.FilePath + data.FileName + "_test" + ".go"
	content := buf.Bytes()
	err = os.WriteFile(modelFile, content, os.ModePerm)
	if err != nil {
		_ = fmt.Errorf("cannot format file: %w", err)
	}
}

func genRender(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
