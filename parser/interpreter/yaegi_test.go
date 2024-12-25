package interpreter

import (
	"fmt"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// 解释器
func ExampleInterpreter() {
	// 参数定义
	// 定义需要执行的 Go 代码
	code := `
import "fmt"
func main() {
	fmt.Println("Hello from yaegi!")
}
`
	// 初始化解释器
	i := interp.New(interp.Options{})
	err := i.Use(stdlib.Symbols)
	if err != nil {
		return
	}
	// 运行代码
	_, err = i.Eval(code)
	if err != nil {
		fmt.Println("Error executing code:", err)
		return
	}
	// Output:
	// Hello from yaegi!

}
