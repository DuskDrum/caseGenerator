package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseAssignment 测试assignment的测试用例，方法调用类型
func TestParseAssignment_call(t *testing.T) {
	parseFile("../example/assignment/conv_function.go")
}

// TestParseAssignment 测试assignment的测试用例，构造类型
func TestParseAssignment_composite(t *testing.T) {
	parseFile("./example/assignment/conv_struct.go")
}

// 解析文件
func parseFile(path string) {
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
				if !ok {
					fmt.Print("不处理非 ast.FuncDecl 的内容")
					continue
				}
				AssignmentWalk := AssignmentWalk{}
				ast.Walk(&AssignmentWalk, decl)
			}
		}
		return nil
	})
	if err != nil {
		return
	}
}

type AssignmentWalk struct {
}

func (v *AssignmentWalk) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	switch node := n.(type) {
	case *ast.AssignStmt:
		si := SourceInfo{}
		si.ParseAssignment(node)
	case *ast.DeclStmt:
		si := SourceInfo{}
		si.ParseAssignment(node)
	// 这种是没有响应值的function
	case *ast.ExprStmt:
		si := SourceInfo{}
		si.ParseAssignment(node)
	}

	return v
}
