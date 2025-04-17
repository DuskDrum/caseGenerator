package expr

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

// TestParseCallExprReceiverResponseCase1
// receiver是直接对应其他包的方法
func TestParseCallExprReceiverResponseCase1(t *testing.T) {
	parseFile("test/test_receiver_function_call.go", "add")

}

// TestParseCallExprReceiverResponseCase2
// receiver是本包内定义的，new出来的变量(在本包中)
func TestParseCallExprReceiverResponseCase2(t *testing.T) {
	// 需要解析的代码字符串
	src := `package expr
	
	
	
	
	func add(a int, b int) int { return a + b }
	`

	// 创建文件集
	fset := token.NewFileSet()

	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 查找函数声明
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			fmt.Printf("Function name: %s\n", funcDecl.Name.Name)

			// 遍历参数列表
			for _, param := range funcDecl.Type.Params.List {
				// 获取每个参数的名称
				for _, name := range param.Names {
					fmt.Printf("Param name: %s\n", name.Name)
				}
				// 获取参数的类型
				fmt.Printf("Param type: %v\n", param.Type)
			}
		}
		return true
	})
}

// TestParseCallExprReceiverResponseCase3
// receiver是本包内定义的，new出来的变量(不在本包中，跨包)
func TestParseCallExprReceiverResponseCase3(t *testing.T) {
	// 需要解析的代码字符串
	src := `package main

	
	
	func add(a int, b int) int { return a + b }
	`

	// 创建文件集
	fset := token.NewFileSet()

	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 查找函数声明
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			fmt.Printf("Function name: %s\n", funcDecl.Name.Name)

			// 遍历参数列表
			for _, param := range funcDecl.Type.Params.List {
				// 获取每个参数的名称
				for _, name := range param.Names {
					fmt.Printf("Param name: %s\n", name.Name)
				}
				// 获取参数的类型
				fmt.Printf("Param type: %v\n", param.Type)
			}
		}
		return true
	})
}

// TestParseCallExprReceiverResponseCase4
// receiver是此方法的receiver调用
func TestParseCallExprReceiverResponseCase4(t *testing.T) {
	// 需要解析的代码字符串
	src := `package main

	
	
	func add(a int, b int) int { return a + b }
	`

	// 创建文件集
	fset := token.NewFileSet()

	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 查找函数声明
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			fmt.Printf("Function name: %s\n", funcDecl.Name.Name)

			// 遍历参数列表
			for _, param := range funcDecl.Type.Params.List {
				// 获取每个参数的名称
				for _, name := range param.Names {
					fmt.Printf("Param name: %s\n", name.Name)
				}
				// 获取参数的类型
				fmt.Printf("Param type: %v\n", param.Type)
			}
		}
		return true
	})
}

// 解析文件
func parseFile(path, funcName string) {
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
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}
