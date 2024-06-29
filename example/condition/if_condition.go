package condition

import (
	"fmt"
	"math/rand"
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
