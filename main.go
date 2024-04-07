package main

import "caseGenerator/parse"

func main() {
	err := parse.Extract("./core")
	if err != nil {
		return
	}
}
