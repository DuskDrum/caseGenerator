package test

import (
	"fmt"
)

type Receiver1 struct {
}

func (r Receiver1) add(a int, b int) int {
	fmt.Print("a +b")
	return a + b
}

func (r *Receiver1) ptrAdd(a int, b int) int {
	fmt.Print("a +b")
	return a + b
}

func InnerAdd(a int, b int) int {
	fmt.Print("a +b")
	return a + b
}
