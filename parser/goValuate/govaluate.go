package goValuate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/expression/govaluate"
	mockresult2 "caseGenerator/parser/expression/govaluate/mockresult"
	"fmt"

	"github.com/samber/lo"
)

// MockBasicExpression mocker 有 basic 的表达式
// 1. 将此段代码前的所有条件都执行，不同的参数都按照公式去计算
// 2. 如果变量定义初始值了，那么每条都执行
// 3. 如果变量没有定义初始值，为请求中的变量或者方法得到，那么需要根据后面倒推得到需要mock的值
func MockBasicExpression(expression *govaluate.ExpressDetail, seList []govaluate.StatementAssignment) []mockresult2.MockResult {
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
	ExpressionList []govaluate.StatementAssignment
	InitValue      any
}

// MockByResult 根据 mockResult来组装 mocker，比如说 IdentMockResult，那么要去找到对应的 ident 是怎么来的，是怎么被赋值的，对应的去做mock
// 变量可以进行 mocker， 指定方法的请求参数
// 不变量不能 mocker，这个不变量需要想办法解析出来
// 方法调用需要mock，mocker 时要考虑私有方法
// selector 一般都是变量，可以想办法用 json 去组装和类型转化
func MockByResult(mockResultList []mockresult2.MockResult) {
	for _, v := range mockResultList {
		switch x := v.(type) {
		case *mockresult2.IdentMockResult:
			// 1. 根据 ident 的 name 找到他对应的类型：是变量还是不变量，是从方法入参进来，还是调用其他方法得到
			fmt.Print(x)
			// 2. 如果是赋值而来，就去找其对应的关系

			// 3. 如果是方法调用得到，那么就去mock 那个方法

		case *mockresult2.CallMockResult:
			// 1. 根据 call 的 信息去mock：直接mock公共方法，或者 go:linkname去mock私有方法
		case *mockresult2.SelectorMockResult:
			// 1. 根据 selector 的 name 找到他对应的类型：是变量还是不变量，是从方法入参进来，还是调用其他方法得到

			// 2. 根据要selector的结构得到要 mocker 的值

			// 3. 如果是赋值而来，就去找其对应的关系

			// 4. 如果是方法调用得到，那么就去mock 那个方法
		}

	}
}
