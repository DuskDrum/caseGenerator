package condition

import (
	"fmt"
	"math/rand"
	"strings"
)

// IfComplexLogicComplexConditionTest 用复杂逻辑运算符连接多个条件
func IfComplexLogicComplexConditionTest() {
	a := rand.Int()
	b := rand.Int()
	c := rand.Int()
	d := rand.Int() > 0
	if a > 0 && (b > 0 || c > 0) && !d {
		fmt.Println("a > 0 && (b > 0 || c >0) && !d")
	}
}

// IfComplexCallConditionTest if中进行接口调用
func IfComplexCallConditionTest() {
	if rand.Int() > 0 {
		fmt.Println("a > 0")
	}
	if strings.Contains("51618617", "1") {
		fmt.Println("strings.Contains(\"51618617\", \"1\")")
	}
}

// IfComplexBoolConditionTest 直接判断布尔类型
func IfComplexBoolConditionTest() {
	isSuccess := true
	if isSuccess {
		fmt.Println("isSuccess is true")
	}
}
