package assignment

import (
	"fmt"
	"strconv"
)

// 此类用于转化数据类型

// ConvFunctionTest1 调整
func ConvFunctionTest1() {
	str := "111"
	atoi, err := strconv.Atoi(str)
	if err != nil {
		fmt.Print("get error")
	}
	if atoi != 111 {
		fmt.Print("convert atoi result is not 111")
	}
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Print("get error")
	}
	if i != 111 {
		fmt.Print("convert result is not 111")
	}
	fmt.Print("convert result is: 111")
}

// ConvFunctionTest2 调整
func ConvFunctionTest2() {
	intInfo := 111
	itoaInt := strconv.Itoa(intInfo)
	if itoaInt != "111" {
		fmt.Print("convert result is not 111")
	}
	formatInt := strconv.FormatInt(int64(intInfo), 10)
	if formatInt != "111" {
		fmt.Print("convert result is not 111")
	}
	fmt.Print("convert result is: 111")
}
