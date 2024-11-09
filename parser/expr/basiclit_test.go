package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestBasicLitCase(t *testing.T) {
	// 定义包含基本字面量的 Go 源代码
	src := `
package main

func main() {
    a := 42
    b := 3.14
    c := "hello"
    d := true
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
		// 检查是否为 *ast.BasicLit 节点
		if basicLit, ok := n.(*ast.BasicLit); ok {
			fmt.Println("找到基本字面量:")

			// 输出字面量的类型和值
			fmt.Printf("  类型: %s\n", basicLit.Kind)
			fmt.Printf("  值: %s\n", basicLit.Value)
		}
		return true
	})
}
