package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestEmptyCase(t *testing.T) {
	src := `
package main

func main() {
	for ; ; {
		// 空语句在 for 循环中充当占位符
	}
	;
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
		// 查找 *ast.EmptyStmt 节点
		if emptyStmt, ok := n.(*ast.EmptyStmt); ok {
			fmt.Println("找到空语句:")
			fmt.Printf("  位置: %v\n", fset.Position(emptyStmt.Semicolon))
			fmt.Printf("  隐式: %v\n", emptyStmt.Implicit)
		}
		return true
	})
}
