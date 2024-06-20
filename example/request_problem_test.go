// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.

package example

import (
	"caseGenerator/common/utils"
	"testing"
)

func Test_request_problemRequestGenericProblem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		list []any

		process func([]any) []any

		batchSize int
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "b0f4ee1e-12c2-11ef-828a-7af6acbff8ed",
			args: args{

				list: make([]any, 0, 10),

				process: nil,

				batchSize: utils.Empty[int](),
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

			RequestGenericProblem(tt.args.list, tt.args.process, tt.args.batchSize)

		})
	}
}

func Test_request_problemRequestGeneric1Problem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		list []any

		process func([]any) ([]any, string)

		batchSize int
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "b0f4f4cc-12c2-11ef-828a-7af6acbff8ed",
			args: args{

				list: make([]any, 0, 10),

				process: nil,

				batchSize: utils.Empty[int](),
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

			RequestGeneric1Problem(tt.args.list, tt.args.process, tt.args.batchSize)

		})
	}
}

func Test_request_problemRequestGeneric3Problem(t *testing.T) {
	type fields struct {
	}
	type args struct {
		list []any

		process func([]any) ([]any, string)

		batchSize int
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "b0f4f54e-12c2-11ef-828a-7af6acbff8ed",
			args: args{

				list: make([]any, 0, 10),

				process: nil,

				batchSize: utils.Empty[int](),
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

			RequestGeneric3Problem(tt.args.list, tt.args.process, tt.args.batchSize)

		})
	}
}
