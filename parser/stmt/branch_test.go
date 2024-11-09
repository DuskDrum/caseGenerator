package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestBranchCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
        outerLoop:
        for i := 0; i < 3; i++ {
            for j := 0; j < 3; j++ {
                if i == 1 && j == 1 {
                    break outerLoop
                }
                fmt.Println(i, j)
            }
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
		if stmt, ok := n.(*ast.BranchStmt); ok {
			fmt.Println("找到 BranchStmt")

			// 被断言的对象
			fmt.Printf("接口对象: %v\n", stmt)
		}
		return true
	})
}
