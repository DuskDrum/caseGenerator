package goValuate

import (
	govaluate2 "caseGenerator/parser/expression/govaluate"
	mockresult2 "caseGenerator/parser/expression/govaluate/mockresult"

	"github.com/Knetic/govaluate"
	"github.com/samber/lo"
)

// MockBasicIntExpression mocker int basic 的表达式
func MockBasicIntExpression(expression *govaluate2.ExpressDetail, basicValueList []any, variablesMap map[string]any, seList []govaluate2.StatementAssignment) []mockresult2.MockResult {
	// 1. 找代码中参数的所有变化，比如说
	// a,b,c := 1,2,3
	// b = c *3
	// a = b+1
	// 那么 a * b * c * d > 500 时需要计算 a、b、c
	// 如果类型是 int
	var inList []int
	inList = append(inList, 0)
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 int
		if value, ok := v.(int); ok {
			inList = append(inList, value)
		} else {
			// 如果转换失败，返回错误
			panic("element is not an int")
		}
	}
	// todo 取列表里的最大最小值，判断怎么快速得到取值范围
	minValue := lo.Min(inList)
	maxValue := lo.Max(inList)
	if minValue == maxValue {
		minValue = minValue - 100
		maxValue = maxValue + 100
	}

	params := make([]string, 0, len(expression.IdentMap))
	resultList := make([]mockresult2.MockResult, 0, 10)

	// 参数名称, 取 ident
	for key := range expression.IdentMap {
		_, ok := variablesMap[key]
		if !ok {
			params = append(params, key)
		}
	}
	// 参数名称，取call
	for key := range expression.CallMap {
		_, ok := variablesMap[key]
		if !ok {
			params = append(params, key)
		}
	}
	// 参数名称，取SelectorMap
	for key := range expression.SelectorMap {
		_, ok := variablesMap[key]
		if !ok {
			params = append(params, key)
		}
	}

	current := make([]int, len(params))
	result := ComposeInt(params, current, 0, minValue, maxValue, expression.Expr, seList, variablesMap)
	if result != nil {
		// 参数一一对应的值
		for i, param := range params {
			if ident, ok := expression.IdentMap[param]; ok {
				imr := &mockresult2.IdentMockResult{
					Ident:     *ident,
					MockValue: result[i],
				}
				resultList = append(resultList, imr)
			}
			if call, ok := expression.CallMap[param]; ok {
				cmr := &mockresult2.CallMockResult{
					Call:      *call,
					MockValue: result[i],
				}
				resultList = append(resultList, cmr)
			}
			if selector, ok := expression.SelectorMap[param]; ok {
				smr := &mockresult2.SelectorMockResult{
					Selector:  *selector,
					MockValue: result[i],
				}
				resultList = append(resultList, smr)
			}
		}
	}
	return resultList
}

// ComposeInt 组合int
func ComposeInt(params []string, current []int, index, min, max int, inequalityExpr string, calculateExprList []govaluate2.StatementAssignment, variablesMap map[string]any) []int {
	// 如果当前索引超出了参数范围，保存组合并返回
	if index == len(params) {
		//fmt.Printf("current is:%v \n", current)
		// 匹配参数是否满足不等式
		// 创建 map
		parameters := make(map[string]any)
		// 迭代并填充 map
		for i := 0; i < len(params); i++ {
			parameters[params[i]] = current[i]
		}
		for k, v := range variablesMap {
			parameters[k] = v
		}
		// 执行完所有前面的赋值
		for _, v := range calculateExprList {
			cExpr, err := govaluate.NewEvaluableExpression(v.Expr)
			if err != nil {
				//return 0, err
				panic("Error calculating")
			}
			result, err := cExpr.Evaluate(variablesMap)
			if err != nil {
				panic("Error calculating ")
			}
			parameters[v.Name] = result
		}

		exp, err := govaluate.NewEvaluableExpression(inequalityExpr)
		if err != nil {
			panic(err.Error())
		}

		result, err := exp.Evaluate(parameters)
		if err != nil {
			panic(err.Error())
		}
		// 尝试类型断言为 bool
		if v, ok := result.(bool); ok {
			if v {
				return current
			}
		}
		return nil
	}

	// 遍历当前参数从 0 到 max 的所有可能值
	for i := min; i <= max; i++ {
		current[index] = i
		composeInt := ComposeInt(params, current, index+1, min, max, inequalityExpr, calculateExprList, variablesMap)
		if composeInt != nil {
			return composeInt
		}
	}
	return nil
}
