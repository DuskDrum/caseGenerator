package example

import (
	"caseGenerator/example/dict"
	"context"
	"time"
)

// RequestSimple 简单的请求类型
func RequestSimple(req1 string, req2 int, req3 int32, req4 float32, req5 float64, req6 bool, req7 any, req8 error) {

}

// RequestStruct struct的请求
func RequestStruct(req1 Example) {

}

// RequestPackStruct 包名引用的struct
func RequestPackStruct(req1 dict.ExampleDict, ctx context.Context, orderTime time.Time) {

}

// RequestFunc 方法的请求
func RequestFunc(req1 func(string, Example, context.Context) (dict.ExampleDict, error)) {

}

// RequestArray 切片的请求
func RequestArray(req1 []string, req2 []int, req3 []bool, req4 []Example, req5 []dict.ExampleDict, req6 [][]string, req7 [][][][]Example, req8 [][]*dict.ExampleDict, req9 [][][]map[string]string, req10 [][][][][][]map[*Example][][][][]*dict.ExampleDict) {

}

// RequestVariable 可变长度参数，一定是最后一位
func RequestVariable(req2 string, req1 ...string) {

}

// RequestVariableParam 多个参数
func RequestVariableParam(req1, req2, req3, req4 string) {

}

// RequestMap map的请求
func RequestMap(req1 map[string]string, req2 map[Example]dict.ExampleDict, req3 map[*Example]*dict.ExampleDict, req4 map[context.Context][]string, req5 map[string][][][][][]*Example, req6 map[string]map[*Example]map[context.Context]map[time.Time]bool) {

}

// RequestStar 指针
func RequestStar(req1 *Example, req2 *dict.ExampleDict, req3 []*Example, req4 *[]Example, req5 *[][][][]*Example, req6 *map[string]string, req7 *map[*Example]map[*dict.ExampleDict][][][][]*Example) {

}

// RequestChan chan
func RequestChan(req1 <-chan string, req2 chan<- string, req3 <-chan Example, req4 <-chan dict.ExampleDict, req5 chan<- *dict.ExampleDict, req6 chan<- [][][][][][][]*Example, req7 chan<- [][][]map[Example][][][][]*dict.ExampleDict) {

}

// RequestEllipsis 泛型
func RequestEllipsis[T, R any](list []T, process func([]T) []R, batchSize int) {

}
