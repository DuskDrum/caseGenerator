package goValuate

import (
	"caseGenerator/parser/expression"
	"caseGenerator/parser/expression/mockresult"
	"caseGenerator/parser/stmt"

	"github.com/Knetic/govaluate"
	"github.com/samber/lo"
)

// MockBasicFloatExpression mocker float basic 的表达式
func MockBasicFloatExpression(expression *expression.ExpressDetail, basicValueList []any, variablesMap map[string]any, seList []stmt.StatementExpression) []mockresult.MockResult {
	// 如果类型是 float
	var inList []float64
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 float
		if value, ok := v.(float64); ok {
			inList = append(inList, value)
		} else if fValue, fok := v.(float32); fok {
			inList = append(inList, float64(fValue))
		} else {
			// 如果转换失败，返回错误
			panic("element is not an float")
		}
	}
	// todo 取列表里的最大最小值，判断怎么快速得到取值范围
	minValue := lo.Min(inList)
	maxValue := lo.Max(inList)
	if minValue == maxValue {
		minValue = minValue - 100
		maxValue = maxValue + 100
	}
	// 参数名称
	params := make([]string, 0, len(expression.IdentMap))
	resultList := make([]mockresult.MockResult, 0, 10)
	// 参数名称, 取 ident
	for key := range expression.IdentMap {
		params = append(params, key)
	}
	// 参数名称，取call
	for key := range expression.CallMap {
		params = append(params, key)
	}
	// 参数名称，取SelectorMap
	for key := range expression.SelectorMap {
		params = append(params, key)
	}

	current := make([]float64, len(params))
	result := ComposeFloat(params, current, 0, minValue, maxValue, expression.Expr, seList, variablesMap)
	if result != nil {
		// 参数一一对应的值
		for i, param := range params {
			if ident, ok := expression.IdentMap[param]; ok {
				imr := &mockresult.IdentMockResult{
					Ident:     *ident,
					MockValue: result[i],
				}
				resultList = append(resultList, imr)
			}
			if call, ok := expression.CallMap[param]; ok {
				cmr := &mockresult.CallMockResult{
					Call:      *call,
					MockValue: result[i],
				}
				resultList = append(resultList, cmr)
			}
			if selector, ok := expression.SelectorMap[param]; ok {
				smr := &mockresult.SelectorMockResult{
					Selector:  *selector,
					MockValue: result[i],
				}
				resultList = append(resultList, smr)
			}
		}
	}
	return resultList
}

// ComposeFloat 组合float
func ComposeFloat(params []string, current []float64, index int, min, max float64, expr string, calculateExprList []stmt.StatementExpression, variablesMap map[string]any) []float64 {
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
				panic("Error calculating ")
			}
			result, err := cExpr.Evaluate(variablesMap)
			if err != nil {
				panic("Error calculating ")
			}
			parameters[v.Name] = result
		}
		exp, err := govaluate.NewEvaluableExpression(expr)
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
		composeInt := ComposeFloat(params, current, index+1, min, max, expr, calculateExprList, variablesMap)
		if composeInt != nil {
			return composeInt
		}
	}
	return nil
}
