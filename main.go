package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./example/request_problem.go")
	if err != nil {
		return
	}
}
