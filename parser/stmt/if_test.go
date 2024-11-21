package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIfCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
        if a := 5; a > 3 {
			fmt.Println("a is greater than 3")
		} else if a < 0 {
			fmt.Println("a is less than or equal to 0")
		} else {
			fmt.Println("a is other probably")
		}
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
		if stmt, ok := n.(*ast.IfStmt); ok {
			fmt.Println("找到 IfStmt")
			parseIf := ParseIf(stmt)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", parseIf)
		}
		return true
	})
}
