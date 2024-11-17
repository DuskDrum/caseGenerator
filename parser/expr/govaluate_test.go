package expr

import (
	"fmt"
	"testing"

	"github.com/Knetic/govaluate"
)

func TestEvaluateArithmeticExpression(t *testing.T) {
	// 测试基本的算术表达式
	expr := "a + b * c"
	parameters := map[string]interface{}{
		"a": 5,
		"b": 3,
		"c": 2,
	}

	exp, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		t.Fatalf("表达式解析失败: %v", err)
	}

	result, err := exp.Evaluate(parameters)
	if err != nil {
		t.Fatalf("表达式计算失败: %v", err)
	}

	expected := 11.0
	if result != expected {
		t.Errorf("预期结果为 %v, 但得到 %v", expected, result)
	}
}

func TestEvaluateConditionExpression(t *testing.T) {
	// 测试基本的算术表达式
	expr := "( ( a + b - c ) * def >= 1 || g < 10 ) && j == 1"
	parameters := map[string]interface{}{
		"a":   5,
		"b":   1,
		"c":   2,
		"def": 3,
		"g":   3,
		"j":   1,
	}

	exp, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		t.Fatalf("表达式解析失败: %v", err)
	}

	result, err := exp.Evaluate(parameters)
	if err != nil {
		t.Fatalf("表达式计算失败: %v", err)
	}

	expected := 11.0
	if result != expected {
		t.Errorf("预期结果为 %v, 但得到 %v", expected, result)
	}
}

func TestEvaluateBooleanExpression(t *testing.T) {
	// 测试布尔表达式
	expr := "x > 10 && y < 20"
	parameters := map[string]interface{}{
		"x": 15,
		"y": 18,
	}

	exp, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		t.Fatalf("表达式解析失败: %v", err)
	}

	result, err := exp.Evaluate(parameters)
	if err != nil {
		t.Fatalf("表达式计算失败: %v", err)
	}

	if result != true {
		t.Errorf("预期结果为 %v, 但得到 %v", true, result)
	}
}

func TestEvaluateFunction(t *testing.T) {
	// 测试内置函数 sqrt
	expr := "sqrt(a + b)"
	parameters := map[string]interface{}{
		"a": 16,
		"b": 9,
	}

	exp, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		t.Fatalf("表达式解析失败: %v", err)
	}

	result, err := exp.Evaluate(parameters)
	if err != nil {
		t.Fatalf("表达式计算失败: %v", err)
	}

	expected := 5.0
	if result != expected {
		t.Errorf("预期结果为 %v, 但得到 %v", expected, result)
	}
}

func TestEvaluateInvalidExpression(t *testing.T) {
	// 测试无效的表达式
	expr := "a + (b *"
	//parameters := map[string]interface{}{
	//	"a": 5,
	//	"b": 3,
	//}

	_, err := govaluate.NewEvaluableExpression(expr)
	if err == nil {
		t.Fatalf("预期表达式解析错误, 但没有发生")
	}
}

func ExampleComposeInt() {
	// 参数定义
	params := []string{"a", "b", "c"} // 参数名称
	minValue := 99                    // 参数的最大值
	maxValue := 101                   // 参数的最大值

	// 初始化结果和中间变量
	current := make([]int, len(params))
	expr := "a+b+c>299"

	// 生成组合
	result := ComposeInt(params, current, 0, minValue, maxValue, expr)
	fmt.Printf("get result, detail is:%+v", result)
	// Output:
	// get result, detail is:[99 100 101]

}

func ExampleComposeFloat() {
	// 参数定义
	params := []string{"a", "b", "c"} // 参数名称
	minValue := 99.1                  // 参数的最大值
	maxValue := 101.2                 // 参数的最大值

	// 初始化结果和中间变量
	current := make([]float64, len(params))
	expr := "a+b+c>299"

	// 生成组合
	result := ComposeFloat(params, current, 0, minValue, maxValue, expr)
	fmt.Printf("get result, detail is:%+v", result)
	// Output:
	// get result, detail is:[99.1 99.1 101.1]
}

func ExampleComposeString() {
	// 参数定义
	params := []string{"a", "b", "c"}
	// 每个参数的可能取值
	values := []string{"", "xx", "xxa"}

	// 初始化结果和中间变量
	current := make([]string, len(params))
	expr := "a+b+c>\"xx\""

	// 生成组合
	result := ComposeString(params, values, current, 0, expr)
	fmt.Printf("get result, detail is:%+v", result)
	// Output:
	// get result, detail is:[  xxa]
}
