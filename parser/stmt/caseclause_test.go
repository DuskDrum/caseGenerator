package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestCaseClauseCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
        num := 2
        switch num {
        case 1:
            fmt.Println("The number is 1")
        case 2:
            fmt.Println("The number is 2")
        default:
            fmt.Println("The number is something else")
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
		if stmt, ok := n.(*ast.CaseClause); ok {
			fmt.Println("找到 CaseClause")
			clause := ParseCaseClause(stmt)

			// 被断言的对象
			fmt.Printf("接口对象: %v\n", clause)
		}
		return true
	})
}
