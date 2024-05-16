package assignment

import (
	"caseGenerator/example/dict"
	"fmt"
)

var (
	varConfigKey = "config.key"
	varRd        = dict.ReceiverDict{}
	// var定义的变量要考虑线程安全问题，如果发现在方法中出现赋值行为，可以使用todo告警或者使用多线程进行测试(使用todo告警)
	varMap  = make(map[string]string, 10)
	varList = make([]string, 0, 10)
)

// VarAssignmentTest1 使用 var进行赋值
func VarAssignmentTest1() {
	varConfigKey = "new.config.key"
	if len(varConfigKey) > 10 {
		fmt.Print("len(configKey) > 10")
	}
	// 这句可以mock
	receiverResult := varRd.TestReceiverFunc(varConfigKey)
	if receiverResult == varConfigKey {
		fmt.Print("len(configKey) > 10")
	}
	varMap["1"] = "1"
	varMap["2"] = "2"
	varMap["3"] = "3"

	varList = append(varList, "111")
	varList = append(varList, "222")
	varList = append(varList, "333")
	varList = append(varList, "444")
}
