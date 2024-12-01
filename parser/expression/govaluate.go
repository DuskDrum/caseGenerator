package expression

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/expression/mockresult"
	"caseGenerator/parser/stmt"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/samber/lo"
)

// MockExpression  mock 表达式
//  1. 如果basicLit有值，那么有靶子了，给其他变量赋值代入表达式中
//     如果是 int、float，可能是==、!=、>、<。那么找到 basicLit 的最大或者最小值，计算得到需要
//     如果是 string，那么可能是==或者!=，也有可能是>、<
//     如果是 nil， 可能是==或者!=。 nil是属于 Ident 里的
//
// 2. 如果两边都没有靶子，那么将其中一边设置为零值，再继续用第一步的流程(ident 不知道变量类型，所以没办法处理)
func MockExpression(expression *Expression, seList []stmt.StatementExpression) []mockresult.MockResult {
	// 1. 如果有 basicLit，那么按照 govalue进行解析试算得到最终结果
	if len(expression.BasicList) > 0 {
		return MockBasicExpression(expression, seList)
	}
	// 2. 如果有nil，那么将其他属性都变成 nil
	resultList := make([]mockresult.MockResult, 0, 10)
	if len(expression.IdentMap) > 0 {
		var nilTag bool
		for _, v := range expression.IdentMap {
			if strings.EqualFold(v.IdentName, "nil") {
				nilTag = true
			}
		}
		if nilTag {
			for _, v := range expression.IdentMap {
				if v.IdentName == "nil" {
					continue
				}
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult.IdentMockResult{
						Ident:     *v,
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult.IdentMockResult{
						Ident:     *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
			for _, v := range expression.CallMap {
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult.CallMockResult{
						Call:      *v,
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult.CallMockResult{
						Call:      *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
			for _, v := range expression.SelectorMap {
				if lo.Contains(expression.ElementList, "!=") {
					result := &mockresult.SelectorMockResult{
						Selector:  *v,
						MockValue: "ZERO",
					}
					resultList = append(resultList, result)
				} else {
					result := &mockresult.SelectorMockResult{
						Selector:  *v,
						MockValue: nil,
					}
					resultList = append(resultList, result)
				}
			}
		}
		return resultList
	}
	// 3. 如果没有 basic 也没有 nil，那么先不处理

	return nil
}

// MockBasicExpression mock 有 basic 的表达式
func MockBasicExpression(expression *Expression, seList []stmt.StatementExpression) []mockresult.MockResult {
	// todo 这种多个basicLit 的类型一般是一样的，不一样就告警出去
	var specificType *enum.SpecificType
	basicValueList := make([]any, 0, 10)
	// 找到不等式中的基本字面量
	for _, v := range expression.BasicList {
		if specificType == nil {
			specificType = &v.SpecificType
		} else if &v.SpecificType != specificType {
			panic("多个 basicLit的类型不一样，请检查")
		}
		basicValueList = append(basicValueList, v.Value)
	}
	// 找到已初始化的变量
	variablesMap := make(map[string]any, 10)
	for _, v := range seList {
		if v.InitParam == nil {
			continue
		}
		_, ok := variablesMap[v.Name]
		if !ok {
			variablesMap[v.Name] = v.InitParam.GetValue()
		}
	}

	// 如果类型是 int、float
	if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_INT {
		return MockBasicIntExpression(expression, basicValueList, variablesMap, seList)
	} else if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_FLOAT64 || lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_FLOAT32 {
		return MockBasicFloatExpression(expression, basicValueList, variablesMap, seList)
	} else if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_STRING {
		return MockBasicStringExpression(expression, basicValueList, variablesMap, seList)
	}
	return nil
}

type StatementExpressionValue struct {
	ExpressionList []stmt.StatementExpression
	InitValue      any
}

// MockBasicIntExpression mock int basic 的表达式
func MockBasicIntExpression(expression *Expression, basicValueList []any, variablesMap map[string]any, seList []stmt.StatementExpression) []mockresult.MockResult {
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
	resultList := make([]mockresult.MockResult, 0, 10)

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

// MockBasicFloatExpression mock float basic 的表达式
func MockBasicFloatExpression(expression *Expression, basicValueList []any, variablesMap map[string]any, seList []stmt.StatementExpression) []mockresult.MockResult {
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

// MockBasicStringExpression mock string basic 的表达式
func MockBasicStringExpression(expression *Expression, basicValueList []any, variablesMap map[string]any, seList []stmt.StatementExpression) []mockresult.MockResult {
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
	resultList := make([]mockresult.MockResult, 0, 10)

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

// ComposeInt 组合int
func ComposeInt(params []string, current []int, index, min, max int, inequalityExpr string, calculateExprList []stmt.StatementExpression, variablesMap map[string]any) []int {
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

// ComposeString 组合string
func ComposeString(params []string, values []string, current []string, index int, expr string, calculateExprList []stmt.StatementExpression, variablesMap map[string]any) []string {
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

// MockByResult 根据 mockResult来组装 mock，比如说 IdentMockResult，那么要去找到对应的 ident 是怎么来的，是怎么被赋值的，对应的去做mock
// 变量可以进行 mock， 指定方法的请求参数
// 不变量不能 mock，这个不变量需要想办法解析出来
// 方法调用需要mock，mock 时要考虑私有方法
// selector 一般都是变量，可以想办法用 json 去组装和类型转化
func MockByResult(mockResultList []mockresult.MockResult) {
	for _, v := range mockResultList {
		switch x := v.(type) {
		case *mockresult.IdentMockResult:
			// 1. 根据 ident 的 name 找到他对应的类型：是变量还是不变量，是从方法入参进来，还是调用其他方法得到
			fmt.Print(x)
			// 2. 如果是赋值而来，就去找其对应的关系

			// 3. 如果是方法调用得到，那么就去mock 那个方法

		case *mockresult.CallMockResult:
			// 1. 根据 call 的 信息去mock：直接mock公共方法，或者 go:linkname去mock私有方法
		case *mockresult.SelectorMockResult:
			// 1. 根据 selector 的 name 找到他对应的类型：是变量还是不变量，是从方法入参进来，还是调用其他方法得到

			// 2. 根据要selector的结构得到要 mock 的值

			// 3. 如果是赋值而来，就去找其对应的关系

			// 4. 如果是方法调用得到，那么就去mock 那个方法
		}

	}
}
