package condition

import (
	"fmt"
	"math/rand"
)

// IfStructureElseIfNestConditionTest if else 嵌套条件测试
func IfStructureElseIfNestConditionTest() {
	a := rand.Int()
	if a > 0 {
		fmt.Println("a > 0")
	} else {
		if a < 0 {
			fmt.Println("a < 0")
			if a < -1 {
				fmt.Println("a < -1")
			} else {
				fmt.Println("0> a > -1")
			}
		} else {
			fmt.Println("a = 0")
		}
	}
}
