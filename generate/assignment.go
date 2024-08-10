package generate

// Assignment 记录变量的信息
type Assignment struct {
}

// 判断是赋值--> 去找Mock、request、其他赋值
// 判断是方法调用--> 去找Mock
// 判断是参数--> 去找Request
// 判断是局部变量--> 往前面找自己
// 判断是全局变量/常量--> 处理不了
