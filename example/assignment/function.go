package assignment

import (
	"caseGenerator/example/dict"
	"errors"
	"fmt"
)

// FunctionTest1 函数测试 本类内调用
func FunctionTest1() {
	targetResult := TargetFunction("", 0, dict.ExampleDict{}, nil)
	fmt.Print(targetResult)
}

// FunctionTest2 函数测试 本类内调用
func FunctionTest2() {
	TargetEmptyFunction()
	fmt.Print("完成调用")
}

// FunctionPackageTest1 函数测试 本包内调用
func FunctionPackageTest1() {
	result := ConstAssignmentTest2("")
	fmt.Print("完成调用: " + result)
}

// FunctionPackageTest2 函数测试 本包内调用
func FunctionPackageTest2() {
	ConstAssignmentTest1()
	fmt.Print("完成调用")
}

// FunctionOutsideTest1 函数测试 外包内调用
func FunctionOutsideTest1() {
	result := ConstAssignmentTest2("")
	fmt.Print("完成调用: " + result)
}

// FunctionOutsideTest2 函数测试 外包内调用
func FunctionOutsideTest2() error {
	ConstAssignmentTest1()
	fmt.Print("完成调用")
	return errors.New("内部错误")
}

func TargetFunction(str string, int2 int, dict dict.ExampleDict, exampleDict *dict.ExampleDict) string {
	return ""
}

func TargetEmptyFunction() {

}
