package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

var (
	newDr = new(dict.ReceiverDict)
)

// NewAssignmentTest1 使用 new进行赋值，new得到的是指针
func NewAssignmentTest1() {
	// 这句可以mock
	receiverResult := newDr.TestReceiverFunc(configKey)
	if receiverResult == configKey {
		fmt.Print("len(configKey) > 10")
	}
}
