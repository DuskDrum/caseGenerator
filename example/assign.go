package example

import "fmt"

func AssignExample() {
	//
	s := test()
	if len(s) > 0 {
		fmt.Print("s>0")
	} else {
		fmt.Println("s<=0")
	}
}

func test() string {
	return ""
}
