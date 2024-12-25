package goValuate

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expression"
	"caseGenerator/parser/expression/mockresult"

	"github.com/Knetic/govaluate"
)

// MockBasicStringExpression mocker string basic 的表达式
func MockBasicStringExpression(expression *expression.ExpressDetail, basicValueList []any, variablesMap map[string]any, seList []bo.StatementAssignment) []mockresult.MockResult {
	// 如果类型是 string
	var strList []string
	strList = append(strList, "")
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 string
		if value, ok := v.(string); ok {
			strList = append(strList, value)
			strList = append(strList, value+"a")
		} else {
			// 如果转换失败，返回错误
			panic("element is not an string")
		}
	}

	// 参数名称
	params := make([]string, 0, len(expression.IdentMap))

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

	current := make([]string, len(params))
	result := ComposeString(params, strList, current, 0, expression.Expr, seList, variablesMap)
	resultList := convertToAnyResultList(expression, result, params)
	return resultList
}

func convertToAnyResultList(expression *expression.ExpressDetail, result, params []string) []mockresult.MockResult {
	resultList := make([]mockresult.MockResult, 0, 10)
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

// ComposeString 组合string

func ComposeString(params []string, values []string, current []string, index int, expr string, calculateExprList []bo.StatementAssignment, variablesMap map[string]any) []string {
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

	// 遍历当前参数的所有可能值
	for _, value := range values {
		current[index] = value
		result := ComposeString(params, values, current, index+1, expr, calculateExprList, variablesMap)
		if result != nil {
			return result
		}
	}
	return nil
}
