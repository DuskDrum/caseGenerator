package param

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFieldCase(t *testing.T) {
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
