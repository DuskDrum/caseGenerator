package mocker

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expression"
	"caseGenerator/parser/expression/mockresult"
	"caseGenerator/parser/goValuate"
	"strings"

	"github.com/samber/lo"
)

// IMock 规范mock执行的逻辑
type IMock interface {
}

// MockExpression  mocker 表达式
//  1. 如果basicLit有值，那么有靶子了，给其他变量赋值代入表达式中
//     如果是 int、float，可能是==、!=、>、<。那么找到 basicLit 的最大或者最小值，计算得到需要
//     如果是 string，那么可能是==或者!=，也有可能是>、<
//     如果是 nil， 可能是==或者!=。 nil是属于 Ident 里的
//
// 2. 如果两边都没有靶子，那么将其中一边设置为零值，再继续用第一步的流程(ident 不知道变量类型，所以没办法处理)
func MockExpression(expression *expression.ExpressDetail, seList []bo.StatementAssignment) []mockresult.MockResult {
	// 1. 如果有 basicLit，那么按照 govalue进行解析试算得到最终结果
	if len(expression.BasicList) > 0 {
		return goValuate.MockBasicExpression(expression, seList)
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
						Ident: *v,
						// todo 需要获取零值
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
						Call: *v,
						// todo 需要获取零值
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
						Selector: *v,
						// todo 需要获取零值
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
