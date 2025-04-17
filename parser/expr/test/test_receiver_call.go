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
