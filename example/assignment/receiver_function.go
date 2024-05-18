package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

func ReceiverFunctionTest1() {
	rd := dict.ReceiverDict{}
	receiverFunc := rd.TestReceiverFunc("")
	fmt.Print(receiverFunc)

}
