package assignment

import "fmt"

// AnonymousFunctionTest1 匿名函数1
func AnonymousFunctionTest1() {
	// 定义一个接受回调函数的函数
	performAction := func(callback func(int) int) (int, error) {
		result := callback(10)
		return result * 2, nil
	}

	// 使用匿名函数作为回调函数
	result, err := performAction(func(x int) int {
		fmt.Println("Callback function called with", x)
		return x + 5
	})

	fmt.Println("Result:", result) // 输出: Callback function called with 10, Result: 30
	fmt.Println("Err:", err)       // 输出: Callback function called with 10, Result: 30
}

// AnonymousFunctionTest2 匿名函数2
func AnonymousFunctionTest2() {
	printers := make([]func(), 5)
	for i := 0; i < 5; i++ {
		// 创建闭包，捕获循环变量i
		printers[i] = func() {
			fmt.Println(i)
		}
	}
	i := printers[0:1]
	fmt.Print(i)
	// 因为i在循环结束后变为5，所有闭包捕获的都是同一个变量i的最终值
	for _, printer := range printers {
		printer() // 输出: 5 5 5 5 5
	}
}
