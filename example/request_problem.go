package example

import (
	"caseGenerator/example/dict"
	"context"
	"time"
)

// RequestResponseSimpleFuncProblem 简单的响应类型
func RequestResponseSimpleFuncProblem(func() (req1 string, req2 int, req3 int32, req4 float32, req5 float64, req6 bool, req7 any, req8 error)) {
}

// RequestResponseStructProblem struct的响应
func RequestResponseStructProblem(func() (req1 Example)) {
}

// RequestResponsePackStructProblem 包名引用的struct
func RequestResponsePackStructProblem(func() (req1 dict.ExampleDict, ctx context.Context, orderTime time.Time)) {
}

// RequestResponseFuncProblem 方法的响应
func RequestResponseFuncProblem(func() (req1 func(string, Example, context.Context) (dict.ExampleDict, error))) {

}

// RequestResponseBlankFuncProblem 方法的响应
func RequestResponseBlankFuncProblem(func() (req1 func())) {
}

// RequestResponseArrayProblem 切片的请求
func RequestResponseArrayProblem(func() (req1 []string, req2 []int, req3 []bool, req4 []Example, req5 []dict.ExampleDict, req6 [][]string, req7 [][][][]Example, req8 [][]*dict.ExampleDict, req9 [][][]map[string]string, req10 [][][][][][]map[*Example][][][][]*dict.ExampleDict)) {
}

// RequestResponseVariableParamProblem 多个参数
func RequestResponseVariableParamProblem(func() (req1, req2, req3, req4 string)) {
}

// RequestResponseMapProblem map的请求
func RequestResponseMapProblem(func() (map[string]string, map[Example]dict.ExampleDict, map[*Example]*dict.ExampleDict, map[context.Context][]string, map[string][][][][][]*Example, map[string]map[*Example]map[context.Context]map[time.Time]bool)) {
}

// RequestResponseStarProblem 指针
func RequestResponseStarProblem(func() (req1 *Example, req2 *dict.ExampleDict, req3 []*Example, req4 *[]Example, req5 *[][][][]*Example, req6 *map[string]string, req7 *map[*Example]map[*dict.ExampleDict][][][][]*Example)) {
}

// RequestResponseChanProblem chan
func RequestResponseChanProblem(func() (<-chan string, chan<- string, <-chan Example, <-chan dict.ExampleDict, chan<- *dict.ExampleDict, chan<- [][][][][][][]*Example, chan<- [][][]map[Example][][][][]*dict.ExampleDict)) {
}

// RequestResponseGenericProblem 泛型
func RequestResponseGenericProblem[T, R any](func() (list []T, process func([]T) []R, batchSize int)) {
}

func RequestResponseGenericValueProblem[T int | uint | int8 | int16 | int32 | int64 | float32 | float64 | string | bool](func() (p *T, s T)) {
}
