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
			// 解析 condition
			conditionResult := parseSwitch.ParseSwitchCondition()
			fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
		}
		return true
	})
}

func TestDuplicateSwitchCase(t *testing.T) {
	sourceCode := `
       package main
    func main() {
        x := 2
    // 使用初始化表达式
    switch y := x * 2; y {
    case 2:
		switch y := y * 2; y {
		case 2:
        	fmt.Println("y is 2")
		case 4:
        	fmt.Println("y is 4")
    	case 6:
        	fmt.Println("y is 6")
    	default:
        	fmt.Println("y is some other value")
    	}
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
			// 解析 condition
			conditionResult := parseSwitch.ParseSwitchCondition()
			fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
		}
		return true
	})
}

func TestParseIfProblemCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		switch c {
	case 9:
		fmt.Println("c is 9")
	case 10:
		fmt.Println("c is 10")
	case 11:
		fmt.Println("c is 11")
	case 12:
		if b == 1 {
			fmt.Println("b is 1")
		} else if b > 2 {
			if b > 2 && b < 10 {
				fmt.Println("b > 2 && b < 10")
				return
			}
			if a == 4 {
				fmt.Println("a is 4")
			}
		}
		fmt.Println("c is 12")
	default:
		if b == 1 {
			fmt.Println("b is 1")
		} else if b > 2 {
			if b > 2 && b < 10 {
				fmt.Println("b > 2 && b < 10")
				return
			}
			if a == 4 {
				fmt.Println("a is 4")
			}
		}
		fmt.Println("c is default")
	}
	fmt.Println("a is greater than 3")
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
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.SwitchStmt); ok {
				fmt.Println("找到 SwitchStmt")
				parseSwitch := ParseSwitch(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseSwitch)
				// 解析 condition
				conditionResult := parseSwitch.ParseSwitchCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
			}
		}
	}
}

func TestParseIfExpressCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 1
		b := 2
		c := a * b + a
		switch c {
	case 9:
		fmt.Println("c is 9")
	case 10:
		fmt.Println("c is 10")
	case 11:
		fmt.Println("c is 11")
	case 12:
		b = a 
		if b == 1 {
			fmt.Println("b is 1")
		} else if b > 2 {
			b = c * 2
			if b > 2 && b < 10 {
				fmt.Println("b > 2 && b < 10")
				return
			}
			c = b + 4
			a = c * 3
			if a == 4 {
				fmt.Println("a is 4")
			}
		}
		fmt.Println("c is 12")
	default:
		a = b
		b = a + c + a	
		if b == 1 {
			a = 4
			fmt.Println("b is 1")
		} else if b > 2 {
			b = a + 2	
			if b > 2 && b < 10 {
				b = 20
				fmt.Println("b > 2 && b < 10")
				return
			}
			c = b + 4
			if a == 4 {
				a = 10
				fmt.Println("a is 4")
			}
		}
		fmt.Println("c is default")
	}
	fmt.Println("a is greater than 3")
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
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.SwitchStmt); ok {
				fmt.Println("找到 SwitchStmt")
				parseSwitch := ParseSwitch(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseSwitch)
				// 解析 condition
				conditionResult := parseSwitch.ParseSwitchCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
			}
		}
	}
}

func TestSimpleExpressCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 1
		b := 2
		c := a * b + a
		switch c {
	case 12:
		b = a
		fmt.Println("c is 12")
	default:
		a = b
		b = a + c + a	
		if b == 1 {
			a = 4
			b = a + 2	
			if b > 2 && b < 10 {
				b = 20
				fmt.Println("b > 2 && b < 10")
				return
			}
			c = b + 4
			if a == 4 {
				a = 10
				fmt.Println("a is 4")
			}
		}
		fmt.Println("c is default")
	}
	fmt.Println("a is greater than 3")
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
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.SwitchStmt); ok {
				fmt.Println("找到 SwitchStmt")
				parseSwitch := ParseSwitch(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseSwitch)
				// 解析 condition
				conditionResult := parseSwitch.ParseSwitchCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
			}
		}
	}
}
