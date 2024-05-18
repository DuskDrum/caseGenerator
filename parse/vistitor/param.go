package vistitor

import (
	"caseGenerator/parse/bo"
	"caseGenerator/utils"
	"encoding/json"
	"fmt"
	"reflect"
)

// Constant 常量
type Constant struct {
	ConstantName  string
	ConstantValue any
}

func (c *Constant) GetParamName() string {
	return c.ConstantName
}

// UnmarshalerInfo 不用处理
func (c *Constant) UnmarshalerInfo(_ string) {
}

type Request struct {
	RequestName  string
	RequestValue any
	// 使用typeOf
	RequestType reflect.Type
}

func (r *Request) GetParamName() string {
	return r.RequestName
}

func (r *Request) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		r.RequestValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

// Variable 变量
type Variable struct {
	VariableName  string
	VariableValue any
	// 使用typeOf
	VariableType reflect.Type
}

func (v *Variable) GetParamName() string {
	return v.VariableName
}

func (v *Variable) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		v.VariableValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

// FuncReturn 方法返回的参数
type FuncReturn struct {
	FuncReturnName  string
	FuncReturnValue any
	FuncReturnIndex int

	// 调用方法的信息
	FuncImport      string
	FuncPackageName string
	FuncName        string
}

func (f *FuncReturn) GetParamName() string {
	return f.FuncReturnName
}

// UnmarshalerInfo 方法返回的结果，是没有类型信息的，需要想办法获取到
func (f *FuncReturn) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		f.FuncReturnValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

type BinaryParam struct {
	ParamName   string
	BinaryParam bo.Binary
}

func (b *BinaryParam) GetParamName() string {
	return b.ParamName
}

func (b *BinaryParam) UnmarshalerInfo(_ string) {

}
