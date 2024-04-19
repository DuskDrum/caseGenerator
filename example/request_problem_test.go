// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.

package example

import (
	"caseGenerator/example/dict"
	"context"
	"testing"
	"time"
)

func Test_request_problemRequestResponseSimpleFuncProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (string, int, int32, float32, float64, bool, any, error)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed30d2a-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseSimpleFuncProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseStructProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() Example
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed3102c-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseStructProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponsePackStructProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (dict.ExampleDict, context.Context, time.Time)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed3109a-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponsePackStructProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseFuncProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() func(string, Example, context.Context) (dict.ExampleDict, error)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed316bc-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseFuncProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseBlankFuncProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() func()
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed3177a-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseBlankFuncProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseArrayProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() ([]string, []int, []bool, []Example, []dict.ExampleDict, [][]string, [][][][]Example, [][]*dict.ExampleDict, [][][]map[string]string, [][][][][][]map[*Example][][][][]*dict.ExampleDict)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed319a0-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseArrayProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseVariableParamProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() string
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed31c70-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseVariableParamProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseMapProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (map[string]string, map[Example]dict.ExampleDict, map[*Example]*dict.ExampleDict, map[context.Context][]string, map[string][][][][][]*Example, map[string]map[*Example]map[context.Context]map[time.Time]bool)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed31f86-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseMapProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseStarProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (*Example, *dict.ExampleDict, []*Example, *[]Example, *[][][][]*Example, *map[string]string, *map[*Example]map[*dict.ExampleDict][][][][]*Example)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed32030-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseStarProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseChanProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (<-chan string, chan<- string, <-chan Example, <-chan dict.ExampleDict, chan<- *dict.ExampleDict, chan<- [][][][][][][]*Example, chan<- [][][]map[Example][][][][]*dict.ExampleDict)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed320bc-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseChanProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseGenericProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() ([]any, func([]any) []any, int)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed32666-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseGenericProblem(tt.args.param0)

		})
	}
}

func Test_request_problemRequestResponseGenericValueProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		param0 func() (*T, bool)
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fed32f12-fe5c-11ee-bcc4-7af6acbff8ec",
			args: args{

				param0: nil,
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

			RequestResponseGenericValueProblem(tt.args.param0)

		})
	}
}
