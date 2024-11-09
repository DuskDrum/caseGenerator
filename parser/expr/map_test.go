package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestMapCase(t *testing.T) {
	// 定义包含 map 的 Go 源代码
	src := `
package main

var m map[string]map[string][][][][][][]map[string][]string
var n map[int]string
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
		// 检查是否为 *ast.MapType 节点
		if mapType, ok := n.(*ast.MapType); ok {
			fmt.Println("找到 map 类型:")

			// 解析 key 和 value 类型
			fmt.Printf("  Key 类型: %T\n", mapType.Key)
			fmt.Printf("  Value 类型: %T\n", mapType.Value)

			// 如果需要进一步解析嵌套类型，可以递归解析 mapType.Value
			printMapType(mapType, 1)
		}
		return true
	})
}

// 递归解析嵌套的 map 类型
func printMapType(node ast.Expr, level int) {
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "\t"
	}

	switch t := node.(type) {
	case *ast.MapType:
		fmt.Printf("%sMap 类型:\n", prefix)
		fmt.Printf("%s  Key 类型: %T\n", prefix, t.Key)
		fmt.Printf("%s  Value 类型: %T\n", prefix, t.Value)
		printMapType(t.Value, level+1) // 递归解析嵌套的 map
	case *ast.ArrayType:
		fmt.Printf("%s数组类型:\n", prefix)
		printMapType(t.Elt, level+1) // 递归解析数组的元素类型
	case *ast.Ident:
		fmt.Printf("%s基础类型: %s\n", prefix, t.Name)
	}
}
