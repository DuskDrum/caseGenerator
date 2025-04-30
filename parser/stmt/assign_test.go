package stmt

import (
	"caseGenerator/parser/bo"
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
		c = []int{1, 2, 3}
        x, y := 1, 2
        z.b = 3
		i := j
        x = y + z
		x += 1
		x -=10
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
			context := bo.ExprContext{
				AstFile:      file,
				AstFuncDecl:  nil,
				RealPackPath: "",
			}
			assign := ParseAssign(stmt, context)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", assign)
		}
		return true
	})
}

func TestAssignCallCase(t *testing.T) {
	sourceCode := `
    package main
	
	import "fmt"
	
    func main() {
		a, b := fmt.Println("解析出错")
		fmt.Println(a)
		fmt.Println(b)
    }
    `
	fset := token.NewFileSet()
	// 解析代码
	file, err := parser.ParseFile(fset, "", sourceCode, 0)
	if err != nil {
		fmt.Println("解析出错:", err.Error())
		return
	}
	// 遍历函数体中的语句
	// 遍历 AST
	ast.Inspect(file, func(n ast.Node) bool {
		if stmt, ok := n.(*ast.AssignStmt); ok {
			fmt.Println("找到 AssignStmt")
			context := bo.ExprContext{
				AstFile:      file,
				AstFuncDecl:  nil,
				RealPackPath: "",
			}
			assign := ParseAssign(stmt, context)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", assign)
		}
		return true
	})
}
