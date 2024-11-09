package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFuncTypeCase(t *testing.T) {
	src := `
package main

func add(a int, b int) int {
    return a + b
}

var subtract func(int, int) int
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
		// 检查是否为 *ast.FuncType 节点
		if funcType, ok := n.(*ast.FuncType); ok {
			fmt.Println("找到函数类型:")
			// 打印参数类型
			if funcType.Params != nil {
				for _, param := range funcType.Params.List {
					fmt.Printf("参数类型: %s\n", param.Type)
				}
			}
			// 打印返回类型
			if funcType.Results != nil {
				for _, result := range funcType.Results.List {
					fmt.Printf("返回类型: %s\n", result.Type)
				}
			}
		}
		return true
	})
}
