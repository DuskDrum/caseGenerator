// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.

package assignment

import (
	"testing"
)

func Test_built_in_functionBuiltInFunctionLen(t *testing.T) {
	type fields struct {
	}
	type args struct {
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "BuiltInFunctionLen",
			args:    args{},
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

			BuiltInFunctionLen()

		})
	}
}

func Test_built_in_functionBuiltInFunctionNew(t *testing.T) {
	type fields struct {
	}
	type args struct {
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "BuiltInFunctionNew",
			args:    args{},
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

			BuiltInFunctionNew()

		})
	}
}

func Test_built_in_functionBuiltInFunctionAppend(t *testing.T) {
	type fields struct {
	}
	type args struct {
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "BuiltInFunctionAppend",
			args:    args{},
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

			BuiltInFunctionAppend()

		})
	}
}

func Test_built_in_functionBuiltInFunctionDelete(t *testing.T) {
	type fields struct {
	}
	type args struct {
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "BuiltInFunctionDelete",
			args:    args{},
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

			BuiltInFunctionDelete()

		})
	}
}
