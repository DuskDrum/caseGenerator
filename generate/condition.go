package generate

// Condition 判断条件
type Condition struct {
}

// 判断某个方法中的条件，找到最简单的路径  ConditionExprType
// 此方法主要用于解决一个条件语句，怎么能走到里面的每一个条件
// 条件需要考虑：1. 用方法的返回值做逻辑运算符的条件 2. 布尔值做逻辑运算符 3. 某个类的属性做逻辑运算符 4. 多个条件做嵌套
// 每个condition一般是一个值做逻辑运算符，比如等于，那么类似于a == b， 或者 a() == b
// 那么处理逻辑如下：
// 	方法 a() ==> Mock
//  参数 a.A ==> Assignment: 1. 赋值里面也是方法：Mock 2. 赋值里面是属性：Request/变量/常量  3. 处理不了
