package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSwitchCase(t *testing.T) {
	sourceCode := `
       package main
    func main() {
        x := 2
    // 使用初始化表达式
    switch y := x * 2; y {
    case 2:
        fmt.Println("y is 2")
    case 4:
        fmt.Println("y is 4")
    case 6:
        fmt.Println("y is 6")
    default:
        fmt.Println("y is some other value")
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
		if stmt, ok := n.(*ast.SwitchStmt); ok {
			fmt.Println("找到 SwitchStmt")
			parseSwitch := ParseSwitch(stmt)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", parseSwitch)
		}
		return true
	})
}
