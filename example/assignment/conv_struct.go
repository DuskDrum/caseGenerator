package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

// ConvStructTest1 struct定义
func ConvStructTest1() {
	str := dict.ExampleDict{
		Name:     "testName",
		Age:      0,
		IsDelete: false,
	}
	fmt.Print("str: ", str)
}

// ConvStructTest2 struct定义
func ConvStructTest2() {
	type args struct {
		merchantId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "检查错误参数",
			args: args{
				merchantId: "",
			},
			wantErr: true,
		},
	}

	fmt.Print("tests: ", tests)
}
