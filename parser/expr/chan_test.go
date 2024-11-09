package expr

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestChanCase(t *testing.T) {
	// 定义包含通道类型的 Go 源代码
	src := `
package main

var ch1 chan int
var ch2 <-chan string
var ch3 chan<- bool

func main() {
	var ch4 chan float64
	_ = ch4
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
		// 检查是否为 *ast.ChanType 节点
		if chanType, ok := n.(*ast.ChanType); ok {
			fmt.Println("找到通道类型:")

			// 输出通道的方向
			switch chanType.Dir {
			case ast.SEND:
				fmt.Println("  通道方向: 发送 (chan<-)")
			case ast.RECV:
				fmt.Println("  通道方向: 接收 (<-chan)")
			default:
				fmt.Println("  通道方向: 双向 (chan)")
			}

			// 输出通道值的类型
			if ident, ok := chanType.Value.(*ast.Ident); ok {
				fmt.Printf("  值类型: %s\n", ident.Name)
			}
		}
		return true
	})
}
