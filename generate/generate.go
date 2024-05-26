package generate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/samber/lo"
)

type GenMeta struct {
	// 测试用例的包名
	Package string
	// 生成的文件名
	FileName string
	// 生成的文件路径
	FilePath string
	// 需要依赖的import
	ImportPkgPaths []string
	// case列表
	CaseDetailList []CaseDetail
}

type CaseDetail struct {
	// 生成的文件名
	FileName string
	// 方法名
	MethodName string
	// receive类型
	ReceiverType string
	// case名称
	CaseName string
	// mock列表
	MockList []*MockInstruct
	// 请求列表
	RequestList []RequestDetail
	// 把请求名按照tt.args.xx，并按照","分割，最后一位没有逗号
	RequestNameString string
}

type MockInstruct struct {
	MockResponseParam []string
	MockFunction      string
	MockNumber        string
	// mock 返回
	MockReturns        string
	MockFunctionParam  []ParamParseResult
	MockFunctionResult []ParamParseResult
}

type ParamParseResult struct {
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

type MockDetail struct {
	// 需要mock方法的包名和方法名
	MockDealMethod    string
	MockMethodPackage string
	MockMethodName    string
	// 依赖的包名
	ImportPkgPath []string
	// 需要go:linkedname关联的方法信息
	InnerFuncList []InnerFunc
}

// InnerFunc 内部方法
type InnerFunc struct {
	// 请求信息
	InnerFuncRequest []RequestDetail
	// 响应信息
	InnerFuncResponse []ResponseDetail
	// 代名
	InnerFuncMethodName string
	// 代名关联的方法路径
	InnerFuncLinkedName string
}

type ResponseDetail struct {
	ResponseType  string
	ImportPkgPath []string
}

type RequestDetail struct {
	// 请求名称
	RequestName string
	// 请求类型
	RequestType string
	// 请求值
	RequestValue string
	// 是否是省略号语法
	IsEllipsis bool
}

func GenFile(data GenMeta) {
	var buf bytes.Buffer

	caseDetails := data.CaseDetailList
	// 设置RequestNameString字段
	cdList := make([]CaseDetail, 0, 10)
	importList := make([]string, 0, 10)
	importList = append(importList, "\"testing\"")
	for _, cd := range caseDetails {
		// 如果是内部方法，那么跳过处理
		// 将字符串转换为rune类型
		firstChar := []rune(cd.MethodName)[0]
		// 判断首字母是否是小写字母
		isLower := unicode.IsLower(firstChar)
		if isLower {
			continue
		}

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

		cd.MockList = lo.Filter(cd.MockList, func(item *MockInstruct, index int) bool {
			if item.MockFunction == "" {
				return false
			} else {
				return true
			}
		})

		cd.MockList = lo.Uniq(cd.MockList)

		if len(cd.MockList) > 0 {
			// mockey.Mock((*repo.ClearingPipeConfigRepo).GetAllConfigs).Return(clearingPipeConfigs).Build()
			fmt.Printf("接收到mock请求，请求详情为:%+v", cd.MockList)
			for i, v := range cd.MockList {
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
				v.MockReturns = str

				v.MockNumber = "mock" + strconv.Itoa(i)
			}

		}
		cdList = append(cdList, cd)
	}
	data.CaseDetailList = cdList

	if len(data.ImportPkgPaths) > 0 {
		for _, v := range data.ImportPkgPaths {
			if v != "" {
				importList = append(importList, v)
			}
		}
	}
	uniqImportList := lo.Uniq(importList)
	data.ImportPkgPaths = uniqImportList

	err := render(NotHaveReceiveModel, &buf, data)
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

func render(tmpl string, wr io.Writer, data interface{}) error {
	t, err := template.New(tmpl).Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}
