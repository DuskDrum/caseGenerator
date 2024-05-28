package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./example/type_assertion.go")
	err1 := parse.Extract("./example/assignment/var.go")
	err2 := parse.Extract("./example/assignment/new.go")
	err3 := parse.Extract("./example/assignment/")
	if err != nil {
		return
	}
	if err1 != nil {
		return
	}
	if err2 != nil {
		return
	}
	if err3 != nil {
		return
	}
}
