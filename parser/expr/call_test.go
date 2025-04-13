package expr

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCallCase(t *testing.T) {
	src := `
package main

import "fmt"

func main() {
	var a = 0
    fmt.Println("Hello, World!")
    a = add(1, 2)
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
		// 检查是否为 *ast.CallExpr 节点
		if callExpr, ok := n.(*ast.CallExpr); ok {
			fmt.Println("找到函数调用表达式:")

			// 打印被调用的函数名
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident:
				fmt.Printf("函数名: %s\n", fun.Name)
			case *ast.SelectorExpr:
				// 处理选择器表达式，例如 fmt.Println
				if x, ok := fun.X.(*ast.Ident); ok {
					fmt.Printf("包名: %s, 函数名: %s\n", x.Name, fun.Sel.Name)
				}
			}

			// 打印参数
			for i, arg := range callExpr.Args {
				fmt.Printf("参数 %d: %v\n", i+1, arg)
			}
		}
		return true
	})
}

func TestCallImportCase(t *testing.T) {
	src := `
package main

import "fmt"

func main() {
	var a = 0
    fmt.Println("Hello, World!")
    a = add(1, 2)
}
`

	// 创建文件集
	fset := token.NewFileSet()

	// 解析源文件
	node, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}

	// 解析导入声明的别名和路径
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		importCtx := build.Default
		pkg, _ := importCtx.Import(path, "", build.FindOnly)
		fmt.Println(pkg.Dir) // 输出包源码目录

		pkgs, _ := parser.ParseDir(fset, pkg.Dir, nil, parser.ParseComments)
		for _, pkgAst := range pkgs {
			ast.Inspect(pkgAst, func(n ast.Node) bool {
				// 检查是否为 *ast.BasicLit 节点
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					if funcDecl.Name.Name == "Println" {
						fmt.Println("找到基本字面量:")
						// 输出字面量的类型和值
						fmt.Printf("  类型: %v\n", funcDecl.Type.Results)
						fmt.Printf("  值: %v\n", funcDecl.Type.Func)
					}
				}
				return true
			})
		}
	}
}
