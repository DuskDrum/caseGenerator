package goValuate

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

type keyFormulas struct {
	Key      string
	Formulas string
}

func TestMultipleLinesArithmeticExpression(t *testing.T) {
	// 定义初始变量
	variables := map[string]any{}

	// 定义公式
	formulas := []keyFormulas{
		{Key: "a", Formulas: "1"},
		{Key: "b", Formulas: "2"},
		{Key: "c", Formulas: "3"},
		{Key: "a", Formulas: "b+a"},
		{Key: "a", Formulas: "b+a"},
		{Key: "b", Formulas: "c * b + a"},
	}

	// 计算所有公式
	for _, formula := range formulas {
		value, err := calculate(formula.Formulas, variables)
		if err != nil {
			fmt.Printf("Error calculating %s: %v\n", formula.Key, err)
			return
		}
		variables[formula.Key] = value
	}

	// 最终计算 a * b * c
	result, err := calculate("a * b * c", variables)
	if err != nil {
		fmt.Printf("Error calculating result: %v\n", err)
		return
	}

	fmt.Printf("Final result: %v\n", result)
}

// 使用 govaluate 计算公式
func calculate(expression string, variables map[string]interface{}) (any, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return 0, err
	}

	result, err := expr.Evaluate(variables)
	if err != nil {
		return 0, err
	}

	return result, nil
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
	result := ComposeInt(params, current, 0, minValue, maxValue, expr, nil, nil)
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
	result := ComposeFloat(params, current, 0, minValue, maxValue, expr, nil, nil)
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
	result := ComposeString(params, values, current, 0, expr, nil, nil)
	fmt.Printf("get result, detail is:%+v", result)
	// Output:
	// get result, detail is:[  xxa]
}
