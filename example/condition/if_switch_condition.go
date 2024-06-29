package condition

import "math/rand"

// SwitchIfCaseNestConditionTest switch case条件测试
func SwitchIfCaseNestConditionTest() {
	rand.Seed(123)
	a := rand.Intn(3)
	switch a {
	case 1:
		b := rand.Intn(2)
		if b == 1 {
			println("a = 1, b = 1")
		} else if b == 2 {
			println("a = 1, b = 2")
		} else {
			println("a = 1, b = 0")
		}
		println("a = 1")
	case 2:
		println("a = 2")
	case 3:
		println("a = 3")
	default:
		println("a = 0")
	}
}

// IfSwitchCaseNestConditionTest switch case条件测试
func IfSwitchCaseNestConditionTest() {
	rand.Seed(123)
	a := rand.Intn(3)
	if a == 1 {
		b := rand.Intn(2)
		switch b {
		case 1:
			println("a = 1, b = 1")
		case 2:
			println("a = 1, b = 2")
		default:
			println("a = 1, b = 0")
		}
		println("a = 1")
	} else if a == 2 {
		println("a = 2")
	} else if a == 3 {
		println("a = 3")
	}

}
