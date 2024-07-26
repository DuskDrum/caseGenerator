package main

import (
	"caseGenerator/parser"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-l", "/var/log/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
	// 1. 解析全部ast字节
	sourceInfo := parser.SourceInfo{}
	sourceInfo.ParseSource("")
	// 2. 根据sourceInfo

	// 3. 解析模板

}
