package assignment

import "fmt"

//  闭包（Closure）是匿名函数的一个特例。当一个匿名函数所访问的变量定义在函数体的外部时，就称这样的匿名函数为闭包。每一个闭包都会绑定一个它自己的外围变量

// ClosureFunctionTest1 闭包
func ClosureFunctionTest1() func() int {
	counter := func() func() int {
		count := 0          // 计数器变量
		return func() int { // 返回一个闭包
			count++ // 闭包可以访问和修改外部函数中的变量
			return count
		}
	}
	next := counter()   // 获取闭包函数
	fmt.Println(next()) // 输出: 1
	fmt.Println(next()) // 输出: 2
	fmt.Println(next()) // 输出: 3
	return next
}

// ClosureFunctionTest2 闭包
func ClosureFunctionTest2() {
	multiplier := func(x int) func(int) int {
		return func(y int) int {
			return x * y
		}
	}

	double := multiplier(2)
	treble := multiplier(3)
	fmt.Println(double(5))
	fmt.Println(treble(7))
}

// ClosureFunctionTest3 闭包
func ClosureFunctionTest3() {
	makePrinter := func(prefix string) func(string) {
		return func(msg string) {
			fmt.Println(prefix, msg)
		}
	}

	printInfo := makePrinter("[INFO] ")
	printWarning := makePrinter("[WARNING] ")

	printInfo("This is an information message.") // 输出: [INFO] This is an information message.
	printWarning("This is a warning message.")   // 输出: [WARNING] This is a warning message.
}
