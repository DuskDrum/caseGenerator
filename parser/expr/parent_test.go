package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestParentCase(t *testing.T) {
	src := `
package main

func main() {
    a := 5
    b := 10
    result := (a + b) * 2
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
		// 检查是否为 *ast.ParenExpr 节点
		if parenExpr, ok := n.(*ast.ParenExpr); ok {
			fmt.Println("找到括号表达式:")

			// 打印括号中的表达式
			fmt.Printf("括号内的表达式: %T\n", parenExpr.X)

			// 若括号内表达式是二元表达式，可以进一步解析
			if binaryExpr, ok := parenExpr.X.(*ast.BinaryExpr); ok {
				fmt.Printf("左操作数: %v\n", binaryExpr.X)
				fmt.Printf("操作符: %v\n", binaryExpr.Op)
				fmt.Printf("右操作数: %v\n", binaryExpr.Y)
			}
		}
		return true
	})
}
