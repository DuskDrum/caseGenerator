package expr

import (
	stmt2 "caseGenerator/parser/expression"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
		binary := ParseBinary(binaryExpr)
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

func TestBinaryCondition(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		if ((a+b()-c())*d.e.f >= 1 || g<10) && j == 1 && j==k && i(z) == nil{
			fmt.Println(i, j)
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
		if stmt, ok := n.(*ast.BinaryExpr); ok {
			fmt.Println("找到 BinaryExpr")
			param := ParseBinary(stmt)
			// 被断言的对象
			fmt.Printf("接口对象: %v\n", param)
			// 解析出所有要对比的列表
			expressionList := stmt2.Express(param)
			fmt.Printf("expression: %v\n", expressionList)
			for _, v := range expressionList {
				mockList := stmt2.MockExpression(v)
				if len(mockList) > 0 {
					fmt.Printf("mock结果列表: %v\n", mockList)
				}
			}
		}
		return true
	})
}
