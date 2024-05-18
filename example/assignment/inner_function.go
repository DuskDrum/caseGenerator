package assignment

import (
	"caseGenerator/example"
	"caseGenerator/example/dict"
	"fmt"
)

// InnerFunctionTest1 内部方法test1
func InnerFunctionTest1() {
	innerEmptyFunction()
	fmt.Print("完成")
}

// InnerFunctionTest2 内部方法test2
func InnerFunctionTest2() {
	withReturn := innerEmptyFunctionWithReturn()
	fmt.Print("完成: " + withReturn)
}

// InnerFunctionTest3 内部方法test3
func InnerFunctionTest3() {
	withReturn := innerEmptyFunctionWithReturn2()
	fmt.Print("完成: " + withReturn)
}

// InnerFunctionTest4 内部方法test2
func InnerFunctionTest4() {
	err, i := innerFunctionWithReturn(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return
	}
	fmt.Println(i)
}

func innerEmptyFunction() {

}

func innerEmptyFunctionWithReturn() string {
	return ""
}

func innerEmptyFunctionWithReturn2() (str string) {
	return
}

func innerFunctionWithReturn(req1 []string, req2 []int, req3 []bool, req4 []example.Example, req5 []dict.ExampleDict, req6 [][]string, req7 [][][][]example.Example, req8 [][]*dict.ExampleDict, req9 [][][]map[string]string, req10 [][][][][][]map[*example.Example][][][][]*dict.ExampleDict) (error, [][][][][][]map[*example.Example][][][][]*dict.ExampleDict) {
	return nil, nil

}
