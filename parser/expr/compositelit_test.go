package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCompositeLitCase(t *testing.T) {
	src := `
package main

func main() {
    // 定义复合字面量
    numbers := []int{1, 2, 3}
	numbers2 := [10]int{1, 2, 3} 
    names := map[string]int{"Alice": 10, "Bob": 20}
    person := struct{Name string}{Name: "Charlie"}
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
			fmt.Println("找到复合字面量:")

			// 打印类型
			fmt.Printf("类型: %T\n", compositeLit.Type)

			// 打印复合字面量的元素
			for _, elt := range compositeLit.Elts {
				fmt.Printf("元素: %s, %s\n", fset.Position(elt.Pos()), elt)
			}
		}
		return true
	})
}
