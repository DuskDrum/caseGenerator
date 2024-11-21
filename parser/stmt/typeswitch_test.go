package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestTypeSwitchCase(t *testing.T) {
	sourceCode := `
   package main
    import "fmt"
    func main() {
        var a interface{} = 10
		switch i := a; v := i.(type) {
		case int:
			fmt.Println("It's an int:", v)
		case string:
			fmt.Println("It's a string:", v)
		default:
			fmt.Println("Unknown type")
		}
		return nil
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
		if stmt, ok := n.(*ast.TypeSwitchStmt); ok {
			fmt.Println("找到 TypeSwitchStmt")
			typeSwitch := ParseTypeSwitch(stmt)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", typeSwitch)
		}
		return true
	})
}

func TestTypeSwitchSimpleCase(t *testing.T) {
	sourceCode := `
   package main
    import "fmt"
    func main() {
        var a interface{} = 10
		switch a.(type) {
		case int:
			fmt.Println("It's an int:")
		case string:
			fmt.Println("It's a string:")
		default:
			fmt.Println("Unknown type")
		}
		return nil
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
		if stmt, ok := n.(*ast.TypeSwitchStmt); ok {
			fmt.Println("找到 TypeSwitchStmt")
			typeSwitch := ParseTypeSwitch(stmt)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", typeSwitch)
		}
		return true
	})
}
