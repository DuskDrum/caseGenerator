package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestUnaryCase(t *testing.T) {
	src := `
package main

func main() {
    a := -10
    b := !true
    c := &a
    d := *c
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
		// 检查是否为 *ast.UnaryExpr 节点
		if unaryExpr, ok := n.(*ast.UnaryExpr); ok {
			fmt.Println("找到一元表达式:")

			// 打印操作符
			fmt.Printf("操作符: %v\n", unaryExpr.Op)

			// 打印操作数
			switch x := unaryExpr.X.(type) {
			case *ast.BasicLit:
				fmt.Printf("操作数: %v (基础字面量)\n", x.Value)
			case *ast.Ident:
				fmt.Printf("操作数: %v (标识符)\n", x.Name)
			default:
				fmt.Printf("操作数类型: %T\n", x)
			}
		}
		return true
	})
}
