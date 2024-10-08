package parser

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseCondition_if 测试condition的测试用例，if
func TestParseCondition_if(t *testing.T) {
	parseConditionFile("../example/condition/if_condition.go")
}

// TestParseCondition_if 测试condition的测试用例，if
func TestParseStructureCondition_if(t *testing.T) {
	parseConditionFile("../example/condition/if_structure_condition.go")
}

// TestParseCondition_switch 测试condition的测试用例，wsitch
func TestParseCondition_switch(t *testing.T) {
	parseConditionFile("../example/condition/switch_condition.go")
}

// TestParseCondition_if_switch
func TestParseCondition_if_switch(t *testing.T) {
	parseConditionFile("../example/condition/if_switch_condition.go")
}

// 解析文件
func parseConditionFile(path string) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			if f, err = parser.ParseFile(fset, path, nil, parser.ParseComments); err != nil {
				return nil
			}

			// 组装所有方法
			for _, cg := range f.Decls {
				decl, ok := cg.(*ast.FuncDecl)
				if ok {
					list := decl.Body.List
					nodes := make([]*ConditionNode, 0, 10)
					for _, stmt := range list {
						switch node := stmt.(type) {
						case *ast.IfStmt:
							si := SourceInfo{}
							conditionNode := si.ParseCondition(node)
							nodes = append(nodes, conditionNode)
						case *ast.SwitchStmt:
							si := SourceInfo{}
							conditionNode := si.ParseCondition(node)
							nodes = append(nodes, conditionNode)
						}
					}
					if len(nodes) == 0 {
						continue
					}
					marshal, err := json.Marshal(nodes)
					if err != nil {
						fmt.Print("ConditionWalk walk err: " + err.Error() + "\n")
					} else {
						fmt.Print("ConditionWalk walk: ", string(marshal)+"\n")
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return
	}
}
