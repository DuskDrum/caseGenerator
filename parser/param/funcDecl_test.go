package param

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFuncDeclCase(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的index表达式
	// 使用 parser.ParseExpr 来解析表达式
	// Go 源代码
	src := `
package main

import "fmt"

func add(a int, b int) int {
    return a + b
}

func main() {
    fmt.Println(add(3, 4))
}
`
	// 创建文件集
	fset := token.NewFileSet()

	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 遍历 AST，查找 *ast.FuncDecl 节点
	ast.Inspect(file, func(n ast.Node) bool {
		// 类型断言，判断是否为 *ast.FuncDecl 类型
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
			fmt.Printf("Function type: %v\n", funcDecl.Type)
			if funcDecl.Body != nil {
				fmt.Println("Function body:")
				for _, stmt := range funcDecl.Body.List {
					fmt.Printf("  %T: %v\n", stmt, stmt)
				}
			}
			fmt.Println()
		}
		return true
	})
}
