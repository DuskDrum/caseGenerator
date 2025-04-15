package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestTypeAssertCase(t *testing.T) {
	src := `
package main

func main() {
	var x interface{} = "hello"
	// 类型断言
	y := x.(string)
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
		// 查找 *ast.TypeAssertExpr 节点
		if typeAssert, ok := n.(*ast.TypeAssertExpr); ok {
			fmt.Println("找到类型断言表达式:")
			assert := ParseTypeAssert(typeAssert, file)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", assert)
		}
		return true
	})
}
