package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSendCase(t *testing.T) {
	sourceCode := `
    package main
    import "fmt"
    func main() {
        ch := make(chan int)
        value := 10
        ch <- value
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
		if stmt, ok := n.(*ast.SendStmt); ok {
			fmt.Println("找到 SendStmt")

			// 被断言的对象
			fmt.Printf("接口对象: %v\n", stmt)
		}
		return true
	})
}
