package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

const CONST_CONFIG_KEY = "config.key"

// ConstAssignmentTest1 使用 const进行赋值
func ConstAssignmentTest1() {
	var testStr = "config"
	rd := dict.ReceiverDict{}
	receiverResult := rd.TestReceiverFunc(testStr)
	if CONST_CONFIG_KEY == receiverResult {
		fmt.Print("CONST_CONFIF_KEY == config.key")
	}
}

// ConstAssignmentTest2 使用 const进行赋值
func ConstAssignmentTest2(str string) string {
	rd := dict.ReceiverDict{}
	receiverResult := rd.TestReceiverFunc(str)
	if CONST_CONFIG_KEY == receiverResult {
		fmt.Print("CONST_CONFIF_KEY == config.key")
	}
	return str
}
