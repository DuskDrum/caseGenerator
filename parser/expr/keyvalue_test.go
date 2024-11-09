package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestKeyValueCase(t *testing.T) {
	src := `
package main

func main() {
    // 定义包含键值对的映射和结构体字面量
    numbers := map[string]int{"one": 1, "two": 2}
    person := struct{
        Name string
        Age  int
    }{
        Name: "Alice",
        Age:  30,
    }
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
		// 检查是否为 *ast.CompositeLit 节点
		if compositeLit, ok := n.(*ast.CompositeLit); ok {
			// 遍历复合字面量中的元素
			for _, elt := range compositeLit.Elts {
				// 检查元素是否为 *ast.KeyValueExpr
				if kvExpr, ok := elt.(*ast.KeyValueExpr); ok {
					fmt.Println("找到键值对表达式:")
					fmt.Printf("键: %v, 值: %v\n", kvExpr.Key, kvExpr.Value)
				}
			}
		}
		return true
	})
}
