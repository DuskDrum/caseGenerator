package stmt

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestIfCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
        if a > 3 {
			fmt.Println("a is greater than 3")
		} else if a < 0 {
			fmt.Println("a is less than or equal to 0")
		} else if a< 1 {
			fmt.Println("a is less than or equal to 1")
		} else {
			fmt.Println("a is other probably")
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
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				// 解析 condition
				conditionResult := parseIf.ParseIfCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
				if len(conditionResult) != 3 {
					panic("result is error")
				}
			}
		}
	}
}

func TestParseIfCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
		b := 10
		c := 11
        if a > 3 {
			if b > 1 {
				fmt.Println("b is greater than 1, return")
				return
			}
			if b<0 {
				switch c {
				case 9:
					fmt.Println("c is 9")
				case 10:
					fmt.Println("c is 10")
				case 11:
					fmt.Println("c is 11")
				case 12:
					fmt.Println("c is 12")
				default:
					if b == 1{
						fmt.Println("b is 1")
					} else if b >2 {
						if b >2 && b <10 {
							fmt.Println("b > 2 && b < 10")	
							return
						} 
						if a == 4 {
							fmt.Println("a is 4")	
						}
					}
					fmt.Println("c is default")
				}
			}
			fmt.Println("a is greater than 3")
		} else if a < 0 {
			fmt.Println("a is less than or equal to 0")
		} else if a< 1 {
			fmt.Println("a is less than or equal to 1")
		} else {
			fmt.Println("a is other probably")
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
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseIf)
			}
		}
	}
}

// TestParseBlockCase 主要用来解析if内部的多个条件分支，是怎么找到条件关系的
// 这个例子最后要得到的结果是:
// 1. a>3 && b < -1 return
// 2. a>3 && b >0 && b ==1 && a ==5
// 3. a>3 && b >0 && b >2 && b >3 && b <10 return
// 4. a>3 && b >0 && b >2 && b >20 && a <10 && a==5
// 5. a>3 && b >0 && b >2 && b >20 && a >100 && b >100 && a==5
// 6. a>3 && b >0 && b >2 && b >20 && a >100 && b <=100 return
// 7. a <= 3 && b > 20
// 8. a <= 3 && b <= 20
func TestParseBlockCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
		b := 10
        if a > 3 {
			if b < -1 {
				fmt.Println("b is greater than 1, return")
				return
			}
			if b>0 {
				if b == 1{
					fmt.Println("b is 1")
				} else if b >2 {
					if b >3 && b <10 {
						fmt.Println("b > 3 && b < 10")	
						return
					} 
					if b > 20 {
						if a < 10 {
							fmt.Println("a < 10 && a > 3")
						} else if a >100 {
							if b > 100 {
								fmt.Println(" a > 100 && b > 100")
							} else {
								fmt.Println(" a > 100 && b < 100 && b>20")
								return
							}
						}
						fmt.Println("a is 4")	
					}
				}
				if a == 5 {
					fmt.Println("a ==5 && b > 0")
				}
			}
		} else {
			fmt.Println("a <= 3")
			if b > 20 {
				fmt.Println("b > 20")
			} else {
				fmt.Println("b <= 20")
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
	// 遍历 AST
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseIf)
				// 解析 condition
				conditionResult := parseIf.ParseIfCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)

			}
		}
	}
}

// TestParseBlockCase 主要用来解析if内部的多个条件分支，是怎么找到条件关系的
// 这个例子最后要得到的结果是:
// 1. a>3
// 2. !(a>3) && b > 20
// 3. !(a>3) && !(b > 20)
func TestSimpleElseParseBlockCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
		b := 10
        if a > 3 {
			fmt.Println("a > 3")
		} else {
			fmt.Println("a <= 3")
			if b > 20 {
				fmt.Println("b > 20")
			} else {
				fmt.Println("b <= 20")
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
	// 遍历 AST
	for _, d := range file.Decls {
		decl, ok := d.(*ast.FuncDecl)
		if !ok {
			fmt.Println("解析出错:")
			return
		}
		for _, b := range decl.Body.List {
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseIf)
				// 解析 condition
				conditionResult := parseIf.ParseIfCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)

			}
		}
	}
}

// TestParseBlockCase 主要用来解析if内部的多个条件分支，是怎么找到条件关系的
// 这个例子最后要得到的结果是:
// 1. b > 0 && a == 5 && b == 10
// 2. b > 0 && a != 5 && a > 5
// 3. b <= 0
func TestParseSimpleBlockCase(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
		b := 10

		if b>0 {
			if a == 5 {
				fmt.Println("a ==5 && b > 0")
	     	} else if a < 5 {
				fmt.Println("a < 5 && b > 0")
			} else if a > 10 {
				fmt.Println("a > 5 && b > 0")
				return
			}
			if b == 10 {
				fmt.Println("b == 10")
				return
		   }
			if b < 10 {
				fmt.Println("b < 10")
			}
			return
		}
    }
    `
	//
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
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseIf)
				// 解析 condition
				conditionResult := parseIf.ParseIfCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
				if len(conditionResult) != 6 {
					panic("result is error")
				}
			}
		}
	}
}

// TestParseZ3Case 主要用来解析if内部的多个条件分支，是怎么找到条件关系的
func TestParseZ3Case(t *testing.T) {
	sourceCode := `
    package main
    func main() {
		a := 5
		b := 10

		if b>0 && a < 1 && a+b > 10 {
			fmt.Println("b>0 && a < 1 && a+b > 10")
		}
    }
    `
	//
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
			if stmt, ok := b.(*ast.IfStmt); ok {
				fmt.Println("找到 IfStmt")
				parseIf := ParseIf(stmt)
				// 被断言的对象
				fmt.Printf("接口对象: %v\n", parseIf)
				// 解析 condition
				conditionResult := parseIf.ParseIfCondition()
				fmt.Printf(" 解析condition后的结果: %v\n", conditionResult)
			}
		}
	}
}
