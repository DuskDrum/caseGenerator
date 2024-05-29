package main

import (
	"caseGenerator/utils"
	"fmt"
)

func main() {
	//err := parse.Extract("./example/assignment/")
	//if err != nil {
	//	return
	//}
	modulePath := utils.GetModulePath()
	fmt.Print(modulePath)

}
