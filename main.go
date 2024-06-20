package main

import (
	"caseGenerator/common/utils"
	"caseGenerator/parse"
	"fmt"
)

func main() {
	err := parse.Extract("./example/assignment/")
	if err != nil {
		return
	}
	modulePath := utils.GetModulePath()
	fmt.Print(modulePath)

}
