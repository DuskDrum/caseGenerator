package factory

import (
	"caseGenerator/generate/faker"
	"sync"
)

var (
	mu sync.RWMutex
)

// GenericMap 定义一个泛型类型，包含 map[string]T
type GenericMap[T any] struct {
	data map[string]faker.Faker[T]
}

// NewGenericMap 定义一个泛型函数，用于创建一个泛型 map
// Go 的泛型参数只能出现在函数、方法或类型定义中，不能直接用于变量定义。
// 提供一个初始化方法
func NewGenericMap[T any]() *GenericMap[T] {
	return &GenericMap[T]{data: make(map[string]faker.Faker[T])}
}

// Get 获取 map 的值
func (gm *GenericMap[T]) Get(key string) faker.Faker[T] {
	return gm.data[key]
}

// Set 设置 map 的值
func (gm *GenericMap[T]) Set(key string, value faker.Faker[T]) {
	gm.data[key] = value
}
