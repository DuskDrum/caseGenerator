package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./example/request.go")
	if err != nil {
		return
	}
}
