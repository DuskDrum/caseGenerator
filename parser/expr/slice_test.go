package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSliceCase(t *testing.T) {
	src := `
package main

func main() {
    a := []int{1, 2, 3, 4, 5}
    b := a[1:3]       // 二参数切片
    c := a[:4]        // 从开头到索引4
    d := a[2:]        // 从索引2到末尾
    e := a[1:3:5]     // 三参数切片
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
		// 查找 *ast.SliceExpr 节点
		if sliceExpr, ok := n.(*ast.SliceExpr); ok {
			fmt.Println("找到切片表达式:")

			// 被切片的对象
			fmt.Printf("被切片的对象: %v\n", sliceExpr.X)

			// 起始索引
			if sliceExpr.Low != nil {
				fmt.Printf("起始索引: %v\n", sliceExpr.Low)
			} else {
				fmt.Println("起始索引: 无（从头开始）")
			}

			// 结束索引
			if sliceExpr.High != nil {
				fmt.Printf("结束索引: %v\n", sliceExpr.High)
			} else {
				fmt.Println("结束索引: 无（到末尾）")
			}

			// 最大索引（仅在三参数切片时有值）
			if sliceExpr.Max != nil {
				fmt.Printf("最大索引: %v\n", sliceExpr.Max)
			} else {
				fmt.Println("最大索引: 无（非三参数切片）")
			}
			fmt.Println()
		}
		return true
	})
}

func TestSliceIdentCase(t *testing.T) {
	src := `
package main

func main() {
    a := []int{1, 2, 3, 4, 5}
    c := 10
	b := a[1:c]
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
		// 查找 *ast.SliceExpr 节点
		if sliceExpr, ok := n.(*ast.SliceExpr); ok {
			fmt.Println("找到切片表达式:")

			// 被切片的对象
			fmt.Printf("被切片的对象: %v\n", sliceExpr.X)

			// 起始索引
			if sliceExpr.Low != nil {
				fmt.Printf("起始索引: %v\n", sliceExpr.Low)
			} else {
				fmt.Println("起始索引: 无（从头开始）")
			}

			// 结束索引
			if sliceExpr.High != nil {
				fmt.Printf("结束索引: %v\n", sliceExpr.High)
			} else {
				fmt.Println("结束索引: 无（到末尾）")
			}

			// 最大索引（仅在三参数切片时有值）
			if sliceExpr.Max != nil {
				fmt.Printf("最大索引: %v\n", sliceExpr.Max)
			} else {
				fmt.Println("最大索引: 无（非三参数切片）")
			}
			fmt.Println()
		}
		return true
	})
}
