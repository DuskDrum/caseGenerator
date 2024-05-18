package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

// go常见的内置函数，len、make、new、append、delete

// BuiltInFunctionLen len
func BuiltInFunctionLen() {
	strings := make([]string, 0, 10)
	i := len(strings)
	if len(strings) > 100 {
		fmt.Print("长度超过 100")
	}
	fmt.Println(i)
}

func BuiltInFunctionNew() {
	d := new(dict.ExampleDict)
	if d == nil {
		fmt.Println("new result is nil")
	}
	fmt.Println("new result is not nil")
}

func BuiltInFunctionAppend() {
	d := make([]string, 0, 10)
	d = append(d, "1")
	d = append(d, "2")
	d = append(d, "3")
	fmt.Println("完成调用")
}

func BuiltInFunctionDelete() {
	targetMap := make(map[string]string, 10)
	targetMap["1"] = "1"
	targetMap["2"] = "2"
	targetMap["3"] = "3"

	delete(targetMap, "1")
	delete(targetMap, "2")
	fmt.Println("完成调用")
}
