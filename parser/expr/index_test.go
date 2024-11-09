package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIndexCase(t *testing.T) {
	// 定义包含索引表达式的 Go 源代码
	src := `
package main

func main() {
	arr := []int{1, 2, 3}
	value := arr[1]
	m := map[string]int{"a": 1, "b": 2}
	val := m["a"]
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
		// 检查是否为 *ast.IndexExpr 节点
		if indexExpr, ok := n.(*ast.IndexExpr); ok {
			fmt.Println("找到索引表达式:")

			// 输出被索引的表达式和索引值
			fmt.Printf("  被索引的表达式: %T\n", indexExpr.X)
			fmt.Printf("  索引值: %T\n", indexExpr.Index)

			// 输出被索引的表达式的类型
			if ident, ok := indexExpr.X.(*ast.Ident); ok {
				fmt.Printf("  表达式名称: %s\n", ident.Name)
			}

			// 输出索引值的类型
			if lit, ok := indexExpr.Index.(*ast.BasicLit); ok {
				fmt.Printf("  索引字面量: %s\n", lit.Value)
			}
		}
		return true
	})
}
