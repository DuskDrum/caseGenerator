package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./example/receiver.go")
	if err != nil {
		return
	}
}
