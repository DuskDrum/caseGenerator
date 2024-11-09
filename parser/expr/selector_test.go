package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSelectorCase(t *testing.T) {
	// 需要解析的源代码
	src := `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    myStruct.Field.ZZ.aa.BB
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
		// 检查是否为 *ast.SelectorExpr 节点
		if selectorExpr, ok := n.(*ast.SelectorExpr); ok {
			fmt.Println("找到选择器表达式:")

			// 获取前缀部分
			if xIdent, ok := selectorExpr.X.(*ast.Ident); ok {
				fmt.Printf("前缀: %s\n", xIdent.Name)
			}

			// 获取选择的标识符
			fmt.Printf("选择的成员: %s\n", selectorExpr.Sel.Name)
		}
		return true
	})
}
