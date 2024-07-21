package condition

import "math/rand"

// SwitchCaseFallthroughConditionTest switch case fallthrough条件测试
func SwitchCaseFallthroughConditionTest() {
	rand.Seed(123)
	a := rand.Intn(10)
	switch a {
	case 1:
		println("a = 1")
		fallthrough
	case 2:
		println("a = 2")
		fallthrough
	case 3:
		println("a = 3")
		fallthrough
	case 4:
		println("a = 4")
		fallthrough
	case 5:
		println("a = 5")
		fallthrough
	case 6:
		println("a = 6")
		fallthrough
	case 7:
		println("a = 7")
		fallthrough
	case 8:
		println("a = 8")
		fallthrough
	case 9:
		println("a = 9")
		fallthrough
	case 10:
		println("a = 10")
		fallthrough
	default:
		println("a = 0")
	}
}

// SwitchCaseNestConditionTest switch case条件测试
func SwitchCaseNestConditionTest() {
	rand.Seed(123)
	a := rand.Intn(3)
	switch a {
	case 1, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20:
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
	case 2:
		println("a = 2")
	case 3:
		println("a = 3")
	default:
		println("a = 0")
	}
}
