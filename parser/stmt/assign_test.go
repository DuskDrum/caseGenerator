package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestAssignCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
        x, y := 1, 2
        z := 3
        x = y + z
    }
    `
	fset := token.NewFileSet()
	// 解析代码
	file, err := parser.ParseFile(fset, "", sourceCode, 0)
	if err != nil {
		fmt.Println("解析出错:", err)
		return
	}
	// 遍历函数体中的语句
	// 遍历 AST
	ast.Inspect(file, func(n ast.Node) bool {
		if stmt, ok := n.(*ast.AssignStmt); ok {
			fmt.Println("找到 AssignStmt")

			// 被断言的对象
			fmt.Printf("接口对象: %v\n", stmt)
		}
		return true
	})
}
