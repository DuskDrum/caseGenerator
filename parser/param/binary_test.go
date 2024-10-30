package param

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"testing"
)

func TestBinaryCase(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的二元表达式
	src := "5 + 4 > 3 + a.b < d"

	// 使用 parser.ParseExpr 来解析表达式
	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 类型断言，判断是否为 ast.BinaryExpr 类型
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok {
		binary := ParseBinary(binaryExpr, "")
		marshal, err := json.Marshal(binary)
		if err != nil {
			panic("Errors: " + err.Error())
		}
		fmt.Printf("Parsed Binary Expression:%s", string(marshal))
	} else {
		fmt.Println("The expression is not a binary expression.")
	}
}

func TestBinaryFuncCall(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的二元表达式
	src := "5 + 3 + aa()"

	// 使用 parser.ParseExpr 来解析表达式
	// parser.ParseExpr函数用于解析单个表达式。适合解析像"myVariable + 5"、"3 * 4"或者"funcCall(arg1, arg2)"这样的表达式
	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 类型断言，判断是否为 ast.BinaryExpr 类型
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok {
		fmt.Printf("Parsed Binary Expression: %v %v %v\n", binaryExpr.X, binaryExpr.Op, binaryExpr.Y)
	} else {
		fmt.Println("The expression is not a binary expression.")
	}
}

func TestBinaryFunc(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的二元表达式
	src := "5 + 3 + aa(1,2)"

	// 使用 parser.ParseExpr 来解析表达式
	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 类型断言，判断是否为 ast.BinaryExpr 类型
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok {
		fmt.Printf("Parsed Binary Expression: %v %v %v\n", binaryExpr.X, binaryExpr.Op, binaryExpr.Y)
	} else {
		fmt.Println("The expression is not a binary expression.")
	}
}

func TestBinarySymbol(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的二元表达式
	src := "a() && 2>3 || !aa(1,2) && c() && d()"

	// 使用 parser.ParseExpr 来解析表达式
	expr, err := parser.ParseExpr(src)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 类型断言，判断是否为 ast.BinaryExpr 类型
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok {
		fmt.Printf("Parsed Binary Expression: %v %v %v\n", binaryExpr.X, binaryExpr.Op, binaryExpr.Y)
	} else {
		fmt.Println("The expression is not a binary expression.")
	}
}
