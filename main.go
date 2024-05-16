package main

import "caseGenerator/parse"

func main() {
	//err := parse.Extract("./example/type_assertion.go")
	err := parse.Extract("./example/request.go")
	if err != nil {
		return
	}
}
