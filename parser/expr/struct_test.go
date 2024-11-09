package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestStructCase(t *testing.T) {
	// 需要解析的源代码
	src := `
package main

type Person struct {
    Name string
    Age  int
}

func main() {
    var p Person
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
		// 检查是否为 *ast.TypeSpec 节点
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			// 检查类型是否为 *ast.StructType
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				fmt.Printf("找到结构体类型: %s\n", typeSpec.Name.Name)
				// 遍历结构体的字段
				for _, field := range structType.Fields.List {
					// 获取字段名
					for _, name := range field.Names {
						fmt.Printf("字段名: %s, 类型: %s\n", name.Name, field.Type)
					}
				}
			}
		}
		return true
	})
}
