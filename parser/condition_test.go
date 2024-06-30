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
					ConditionWalk := ConditionWalk{}
					ast.Walk(&ConditionWalk, decl)
				}
				genDecl, ok := cg.(*ast.GenDecl)
				if ok {
					si := SourceInfo{}
					assignment := si.ParseAssignment(genDecl)
					if assignment != nil {
						marshal, err := json.Marshal(assignment)
						if err != nil {
							fmt.Print("ConditionWalk walk err: " + err.Error() + "\n")
						} else {
							fmt.Print("ConditionWalk walk: ", string(marshal)+"\n")
						}
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

type ConditionWalk struct {
}

func (v *ConditionWalk) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	var assignment []*Assignment
	switch node := n.(type) {
	case *ast.AssignStmt:
		si := SourceInfo{}
		assignment = si.ParseAssignment(node)
	case *ast.DeclStmt:
		si := SourceInfo{}
		assignment = si.ParseAssignment(node)
	case *ast.GenDecl:
		si := SourceInfo{}
		assignment = si.ParseAssignment(node)
	// 这种是没有响应值的function
	case *ast.ExprStmt:
		si := SourceInfo{}
		assignment = si.ParseAssignment(node)
	}
	if assignment == nil {
		return v
	}
	marshal, err := json.Marshal(assignment)
	if err != nil {
		fmt.Print("AssignmentWalk walk err: " + err.Error() + "\n")
	} else {
		fmt.Print("AssignmentWalk walk: ", string(marshal)+"\n")
	}

	return v
}
