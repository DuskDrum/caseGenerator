package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestDeclCase(t *testing.T) {
	src := `
package main

func main() {
	var x int = 10
	const y = "hello"
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
		// 查找 *ast.DeclStmt 节点
		if declStmt, ok := n.(*ast.DeclStmt); ok {
			fmt.Println("找到声明语句:")

			// 打印声明的具体信息
			switch decl := declStmt.Decl.(type) {
			case *ast.GenDecl:
				fmt.Printf("通用声明 (token: %s):\n", decl.Tok)
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.ValueSpec:
						fmt.Printf("  名称: %v\n", spec.Names)
						fmt.Printf("  类型: %v\n", spec.Type)
						fmt.Printf("  值: %v\n", spec.Values)
					}
				}
			default:
				fmt.Println("未知的声明类型")
			}
		}
		return true
	})
}
