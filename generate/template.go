package generate

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
			convert.{{.MethodName}}({{.RequestNameString}})
			{{else}}
			{{if .MockList}}
			{{range .MockList}}
			Mock({{.MockDealMethod}}).Return({{.MockResponseList}}).Build()
			{{end}}
			{{end}}
			{{.MethodName}}({{.RequestNameString}})
			{{end}}
		})
	}
}
{{end}}


`
