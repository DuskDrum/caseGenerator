package test

import (
	"caseGenerator/parser/expr/test/testpackage"
	"fmt"
)

var testService = new(testpackage.TestService)
var receiver1 = new(Receiver1)

type Receiver struct {
}

// add 直接调用fmt包里的方法，不是receiver
func add(a int, b int) int {
	fmt.Print("a +b")
	return a + b
}

// AddPackageService 跨包new的receiver
func AddPackageService(a int, b int) int {
	return testService.Add(a, b)
}

// AddReceiver 直接调用fmt包里的方法，不是receiver
func (r Receiver) AddReceiver(a int, b int) int {
	return r.addReceiver(a, b)
}

func (r Receiver) addReceiver(a int, b int) int {
	fmt.Print("a +b")
	return a + b
}

func AddReceiver1(a int, b int) int {
	return receiver1.add(a, b)
}
