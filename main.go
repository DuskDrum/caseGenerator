package main

import "caseGenerator/parse"

func main() {
	//err := parse.Extract("./example/type_assertion.go")
	//err := parse.Extract("./example/assignment/var.go")
	err := parse.Extract("./example/assignment/new.go")
	if err != nil {
		return
	}
}
