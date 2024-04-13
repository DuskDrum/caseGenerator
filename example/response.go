package example

import (
	"caseGenerator/example/dict"
	"context"
	"time"
)

// ResponseSimple 简单的响应类型
func ResponseSimple() (req1 string, req2 int, req3 int32, req4 float32, req5 float64, req6 bool, req7 any, req8 error) {
	return req1, 0, 0, 0, 0, false, nil, req8
}

// ResponseStruct struct的响应
func ResponseStruct() (req1 Example) {
	return Example{}
}

// ResponsePackStruct 包名引用的struct
func ResponsePackStruct() (req1 dict.ExampleDict, ctx context.Context, orderTime time.Time) {
	return dict.ExampleDict{}, nil, time.Time{}
}

// ResponseFunc 方法的响应
func ResponseFunc() (req1 func(string, Example, context.Context) (dict.ExampleDict, error)) {
	return func(s string, example Example, ctx context.Context) (dict.ExampleDict, error) {
		return dict.ExampleDict{}, nil
	}
}

// ResponseBlankFunc 方法的响应
func ResponseBlankFunc() (req1 func()) {
	return func() {

	}
}

// ResponseArray 切片的请求
func ResponseArray() (req1 []string, req2 []int, req3 []bool, req4 []Example, req5 []dict.ExampleDict, req6 [][]string, req7 [][][][]Example, req8 [][]*dict.ExampleDict, req9 [][][]map[string]string, req10 [][][][][][]map[*Example][][][][]*dict.ExampleDict) {
	return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil
}

// ResponseVariableParam 多个参数
func ResponseVariableParam() (req1, req2, req3, req4 string) {
	return "", "", "", ""
}

// ResponseMap map的请求
func ResponseMap() (map[string]string, map[Example]dict.ExampleDict, map[*Example]*dict.ExampleDict, map[context.Context][]string, map[string][][][][][]*Example, map[string]map[*Example]map[context.Context]map[time.Time]bool) {
	return nil, nil, nil, nil, nil, nil
}

// ResponseStar 指针
func ResponseStar() (req1 *Example, req2 *dict.ExampleDict, req3 []*Example, req4 *[]Example, req5 *[][][][]*Example, req6 *map[string]string, req7 *map[*Example]map[*dict.ExampleDict][][][][]*Example) {
	return nil, nil, nil, nil, nil, nil, nil
}

// ResponseChan chan
func ResponseChan() (<-chan string, chan<- string, <-chan Example, <-chan dict.ExampleDict, chan<- *dict.ExampleDict, chan<- [][][][][][][]*Example, chan<- [][][]map[Example][][][][]*dict.ExampleDict) {
	return nil, nil, nil, nil, nil, nil, nil
}

// ResponseGeneric 泛型
func ResponseGeneric[T, R any]() (list []T, process func([]T) []R, batchSize int) {
	return nil, func(ts []T) []R {
		return nil
	}, 0
}

func ResponseGenericValue[T int | uint | int8 | int16 | int32 | int64 | float32 | float64 | string | bool](p *T) T {
	if p == nil {
		//基础类型初始化0值
		var t T
		return t
	}
	return *p
}
