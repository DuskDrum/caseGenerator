package assignment

import (
	"caseGenerator/example"
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
	intInfo, strngInfo, strngInfo2, strngInfo3 := 111, "222", "3333", "4444"
	inad := strngInfo
	inadx := &strngInfo
	result := inad + strngInfo2 + strngInfo3
	fmt.Print("convert result is not 111" + result)

	fmt.Print("convert result is not 111" + inad + strngInfo2)
	fmt.Print(inadx)

	e := example.Example{}
	fmt.Print(e)
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

// ConvFunctionTest3 调整
func ConvFunctionTest3() {
	str, str2, str3 := "111", "2222", 333333
	_, err := strconv.Atoi(str)
	if err != nil {
		fmt.Print("get error")
	}

	fmt.Print("convert result is: 111")
	fmt.Print("convert str2 is: " + str2)
	fmt.Print("convert str3 is: " + strconv.Itoa(str3))
}
