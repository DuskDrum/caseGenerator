package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// 未得到 indexList
func TestInterfaceTypeCase(t *testing.T) {
	src := `
package main

type MyInterface interface {
    Method1(a int) int
    Method2(b string)
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
			// 检查类型是否为 *ast.InterfaceType
			if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
				fmt.Printf("找到接口类型: %s\n", typeSpec.Name.Name)
				// 遍历接口中的方法
				for _, method := range interfaceType.Methods.List {
					// 检查是否是方法
					if funcType, ok := method.Type.(*ast.FuncType); ok {
						fmt.Printf("方法名: %s\n", method.Names[0].Name)
						// 打印参数类型
						if funcType.Params != nil {
							for _, param := range funcType.Params.List {
								fmt.Printf("  参数类型: %s\n", param.Type)
							}
						}
						// 打印返回类型
						if funcType.Results != nil {
							for _, result := range funcType.Results.List {
								fmt.Printf("  返回类型: %s\n", result.Type)
							}
						}
					}
				}
			}
		}
		return true
	})
}
