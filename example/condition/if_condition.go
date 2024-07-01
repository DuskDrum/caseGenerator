package condition

import (
	"fmt"
	"math/rand"
	"strings"
)

// IfConditionTest if条件测试
func IfConditionTest() {
	a := rand.Int()
	if a > 0 {
		fmt.Println("a > 0")
	}
}

// IfElseConditionTest if else条件测试
func IfElseConditionTest() {
	a := rand.Int()
	if a > 0 {
		fmt.Println("a > 0")
	} else {
		fmt.Println("a <= 0")
	}
}

// IfElseIfConditionTest if else if条件测试
func IfElseIfConditionTest() {
	a := rand.Int()
	if a > 0 {
		fmt.Println("a > 0")
	} else if a < 0 {
		fmt.Println("a < 0")
	} else {
		fmt.Println("a = 0")
	}
}

// IfElseIfNestConditionTest if else 嵌套条件测试
func IfElseIfNestConditionTest() {
	a := rand.Int()
	if a > 0 {
		fmt.Println("a > 0")
	} else {
		if a < 0 {
			fmt.Println("a < 0")
		} else {
			fmt.Println("a = 0")
		}
	}
}

// IfCompareConditionTest 比较
func IfCompareConditionTest() {
	a := rand.Int()
	b := rand.Int()
	if a > b {
		fmt.Println("a > b")
	}
}

// IfLogicConditionTest 是用逻辑运算符连接多个条件
func IfLogicConditionTest() {
	a := rand.Int()
	b := rand.Int()
	if a > 0 && b > 0 {
		fmt.Println("a > 0 && b > 0")
	}
}

// IfLogicComplexConditionTest 用复杂逻辑运算符连接多个条件
func IfLogicComplexConditionTest() {
	a := rand.Int()
	b := rand.Int()
	c := rand.Int()
	d := rand.Int() > 0
	if a > 0 && (b > 0 || c > 0) && !d {
		fmt.Println("a > 0 && (b > 0 || c >0) && !d")
	}
}

// IfCallConditionTest if中进行接口调用
func IfCallConditionTest() {
	if rand.Int() > 0 {
		fmt.Println("a > 0")
	}
	if strings.Contains("51618617", "1") {
		fmt.Println("strings.Contains(\"51618617\", \"1\")")
	}
}

// IfTypeAssertConditionTest if中进行类型断言
func IfTypeAssertConditionTest() {
	var i interface{}
	i = "hello"
	if s, ok := i.(string); ok {
		fmt.Println("s is string" + s)
	}
}

// IfZeroValueCheckConditionTest 零值检查
func IfZeroValueCheckConditionTest() {
	var err error
	if err != nil {
		fmt.Println("err is not nil")
	}
}

// IfBoolConditionTest 直接判断布尔类型
func IfBoolConditionTest() {
	isSuccess := true
	if isSuccess {
		fmt.Println("isSuccess is true")
	}
}

// IfChanConditionTest 使用通道的操作结果进行判断
func IfChanConditionTest() {
	ch := make(chan int)
	if v, ok := <-ch; ok {
		fmt.Println("v is " + string(rune(v)))
	}
}

// IfAssignmentConditionTest 赋值并立即检查
func IfAssignmentConditionTest() {
	if x := rand.Int(); x > 0 {
		fmt.Println("x > 0")
	}
}
