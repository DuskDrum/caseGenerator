package decl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestArrayCase(t *testing.T) {
	// 需要解析的源代码
	src := `
	package main
	func add(x, y int) int { 
		return x + y 
	}
`

	// 创建一个文件集
	fset := token.NewFileSet()

	// 解析源代码
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历 AST，查找数组类型
	ast.Inspect(file, func(n ast.Node) bool {
		// 判断节点是否为 *ast.ArrayType
		if arrayType, ok := n.(*ast.ArrayType); ok {
			if arrayType.Len != nil {
				fmt.Printf("Array with length: %v\n", arrayType.Len)
			} else {
				fmt.Println("Slice type detected")
			}
			fmt.Printf("Element type: %v\n", arrayType.Elt)
		}
		return true
	})
}
