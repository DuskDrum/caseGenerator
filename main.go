package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./example/type_assertion.go")
	if err != nil {
		return
	}
}
