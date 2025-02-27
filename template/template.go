package template

import (
	"bytes"
	"caseGenerator/common/utils"
	"caseGenerator/generate_old"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

const NotEditMark = `// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
`

// NotHaveReceiveModel used as a variable because it cannot load template file after packed, params still can pass file
const NotHaveReceiveModel = NotEditMark + `
package {{.Package}}

import (
	{{range .ImportPkgPaths}} {{.}} ` + "\n" + `{{end}}
)

{{range .CaseDetailList}}
func Test_{{.FileName}}{{.CaseName}}(t *testing.T) {
	type fields struct {
	}
	type args struct {
		{{range .RequestList}}
		{{.RequestName}}  {{.RequestType}} ` + `
		{{end}}
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "{{.CaseName}}",
			args: args{
				{{range .RequestList}}
				{{.RequestName}}: {{.RequestValue}}, 
				{{end}}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				//统一处理
				if err := recover(); err != nil {
					t.Log("recover了")
				}
			}()
			{{if .ReceiverType}}
			convert := {{.ReceiverType}}
			{{if .MockList}}
			{{range .MockList}}
			Mock({{.MockDealMethod}}).Return({{.MockResponseList}}).Build()
			{{end}}
			{{end}}
			convert.{{.FunctionName}}({{.RequestNameString}})
			{{else}}
			{{if .MockList}}
			{{range .MockList}}
			{{.MockNumber}} := bytedanceMockey.Mock({{.MockFunction}}).Return({{.MockReturns}}).Build()
			defer {{.MockNumber}}.UnPatch()
			{{end}}
			{{end}}
			{{.FunctionName}}({{.RequestNameString}})
			{{end}}
		})
	}
}
{{end}}


{{range}}


`

func GenGenerateFile(data generate_old.GenMeta) {
	var buf bytes.Buffer

	caseDetails := data.CaseDetailList
	// 设置RequestNameString字段
	cdList := make([]generate_old.CaseDetail, 0, 10)
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

	data.MockList = lo.Filter(data.MockList, func(item *generate_old.MockInstruct, index int) bool {
		if item.MockFunction == "" {
			return false
		} else {
			return true
		}
	})
	mockInstructs := make([]*generate_old.MockInstruct, 0, 10)

	// mockey.Mock((*repo.ClearingPipeConfigRepo).GetAllConfigs).Return(clearingPipeConfigs).Build()
	// go:linkname awxCommonConvertSettlementReportAlgorithm slp/reconcile/core/message/standard.commonConvertSettlementReportAlgorithm
	// func awxCommonConvertSettlementReportAlgorithm(transactionType enums.TransactionType, createdAt time.Time, ctx context.Context, dataBody dto.AwxSettlementReportDataBody) (result []service.OrderAlgorithmResult, err error)
	if len(data.MockList) > 0 {
		by := lo.SliceToMap(data.MockList, func(item *generate_old.MockInstruct) (string, *generate_old.MockInstruct) {
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
			mockNumber := "mocker" + k
			// 3. 如果方法名是小写开头，且没有包名引用，说明需要go-linkname
			if utils.IsLower(v.MockFunction) && !strings.Contains(v.MockFunction, ".") {

			}

			mi := generate_old.MockInstruct{
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

	err := genRender(NotHaveReceiveModel, &buf, data)
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
