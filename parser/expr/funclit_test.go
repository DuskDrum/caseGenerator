package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFuncLitCase(t *testing.T) {
	src := `
package main

func main() {
    // 定义一个匿名函数并调用
    add := func(a int, b int) int {
        return a + b
    }
    fmt.Println(add(1, 2))
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
		// 检查是否为 *ast.FuncLit 节点
		if funcLit, ok := n.(*ast.FuncLit); ok {
			fmt.Println("找到匿名函数:")
			// 打印参数类型
			if funcLit.Type.Params != nil {
				for _, param := range funcLit.Type.Params.List {
					fmt.Printf("参数类型: %s\n", param.Type)
				}
			}
			// 打印返回类型
			if funcLit.Type.Results != nil {
				for _, result := range funcLit.Type.Results.List {
					fmt.Printf("返回类型: %s\n", result.Type)
				}
			}
			// 打印函数体
			fmt.Printf("函数体: %v\n", funcLit.Body)
		}
		return true
	})
}
