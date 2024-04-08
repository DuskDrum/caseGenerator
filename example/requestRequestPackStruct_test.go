// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.
// Code generated by codeParser. DO NOT EDIT.

package example

import (
	"caseGenerator/example/dict"
	"context"
	"slp/reconcile/core/common/utils"
	"testing"
	"time"
)

func Test_requestRequestPackStruct(t *testing.T) {
	type fields struct {
	}
	type args struct {
		req1 dict.ExampleDict

		ctx context.Context

		orderTime time.Time
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "1542480a-f59d-11ee-af52-7af6acbff8ed",
			args: args{

				req1: utils.Empty[dict.ExampleDict](),

				ctx: context.Background(),

				orderTime: time.Now(),
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

			RequestPackStruct(tt.args.req1, tt.args.ctx, tt.args.orderTime)

		})
	}
}
