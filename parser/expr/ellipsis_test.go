package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestEllipsisCase(t *testing.T) {
	// 定义包含省略号的 Go 源代码
	src := `
package main

func example(args ...int) {}

func main() {
	numbers := []int{1, 2, 3}
	example(numbers...)
}
`
	// 创建文件集
	fset := token.NewFileSet()

	// 解析源文件
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	// 遍历 AST
	ast.Inspect(file, func(n ast.Node) bool {
		// 查找 *ast.Ellipsis 类型的节点
		if ellipsis, ok := n.(*ast.Ellipsis); ok {
			fmt.Println("找到 Ellipsis (省略号) 节点:")

			// 解析省略号的类型
			if ellipsis.Elt != nil {
				fmt.Printf("  参数类型: %T\n", ellipsis.Elt)
				if ident, ok := ellipsis.Elt.(*ast.Ident); ok {
					fmt.Printf("  类型名称: %s\n", ident.Name)
				}
			} else {
				fmt.Println("  没有指定类型（在调用中使用 `...`）")
			}
		}
		return true
	})
}
