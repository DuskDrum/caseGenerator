package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestStarCase(t *testing.T) {
	src := `
package main

func main() {
    var x *int
    y := *x
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
		// 检查是否为 *ast.StarExpr 节点
		if starExpr, ok := n.(*ast.StarExpr); ok {
			fmt.Println("找到指针或解引用表达式:")

			// 打印被指向的类型或解引用的表达式
			switch x := starExpr.X.(type) {
			case *ast.Ident:
				fmt.Printf("指向的类型或变量: %s\n", x.Name)
			default:
				fmt.Printf("表达式类型: %T\n", starExpr.X)
			}
		}
		return true
	})
}
