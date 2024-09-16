package main

import (
	"caseGenerator/generate"
	pparser "caseGenerator/parser"
	"fmt"
	"testing"
)

// TestParseCondition_if 测试condition的测试用例，if
func TestParseCondition_if(t *testing.T) {
	sie := pparser.SourceInfoExt{}
	sie.ParseSource("../example/condition/if_complex_condition.go")
	fmt.Printf("执行的结果为: %+v", sie)

	for _, source := range sie.SourceInfo {
		for _, condition := range source.ConditionList {
			generate.GenerateCondition(condition)
		}
	}

}
