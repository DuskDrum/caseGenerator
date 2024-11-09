package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestLabeledCase(t *testing.T) {
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
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			for _, stmt := range funcDecl.Body.List {
				if labeledStmt, ok := stmt.(*ast.LabeledStmt); ok {
					fmt.Println("找到LabeledStmt")
					fmt.Println("标签名称:", labeledStmt.Label.Name)
					// 可以进一步处理语句部分（labeledStmt.Stmt）等
				}
			}
		}
	}
}
