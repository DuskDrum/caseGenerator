package z3

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIdentZ3Case 测试ident变为Z3的逻辑
func TestIdentZ3Case(t *testing.T) {
	// 测试内部Receiver
	path := "test/ident.go"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expression/z3/test"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					// 查找函数声明
					ast.Inspect(f, func(n ast.Node) bool {
						// 检查是否为 *ast.CallExpr 节点
						if callExpr, ok := n.(*ast.Ident); ok {
							fmt.Println("找到函数调用表达式:")
							context := bo.ExprContext{
								AstFile:      f,
								AstFuncDecl:  funcDecl,
								RealPackPath: relPackagePath,
							}
							// 解析出了call
							call := expr.ParseIdent(callExpr, context)
							// 执行z3处理器
							eContext := bo.ExpressionContext{
								VariableParamMap: make(map[string]enum.BasicParameterType, 10),
								RequestParamMap:  make(map[string]enum.BasicParameterType, 10),
								TemporaryVariable: bo.TemporaryVariable{
									VariableName: "localVariable",
								},
								ExprContext: context,
							}
							expressCall, _ := ExpressIdent(call, eContext)

							fmt.Print(expressCall)

							fmt.Print(call)
						}
						return true
					})

				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}

}

// TestIdentRequestZeroInt 测试方法入参
func TestIdentRequestZeroInt(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "RequestZeroInt"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿方法入参的值，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalBool 测试公共bool变量
func TestIdentGlobalBool(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalBool"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalNumber 测试公共int变量
func TestIdentGlobalNumber(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalNumber"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalReferenceNumber 测试公共引用变量，引用其他局部变量
func TestIdentGlobalReferenceNumber(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalReferenceNumber"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalString 测试公共string变量
func TestIdentGlobalString(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalString"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalReceiver 测试公共struct变量
func TestIdentGlobalReceiver(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalReceiver"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalReceiverName 测试公共struct->string变量
func TestIdentGlobalReceiverName(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalReceiverName"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalReceiverAge 测试公共struct->int变量
func TestIdentGlobalReceiverAge(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalReceiverAge"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalReceiverIsPass 测试公共struct->bool变量
func TestIdentGlobalReceiverIsPass(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalReceiverIsPass"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalCallReceiverString 测试公共struct->string变量
func TestIdentGlobalCallReceiverString(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalCallReceiverString"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalCallInt 测试公共call->int变量
func TestIdentGlobalCallInt(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalCallInt"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalSelectorInt 测试公共selector->int变量
func TestIdentGlobalSelectorInt(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalSelectorInt"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

// TestIdentGlobalSelectorStruct 测试公共selector->struct变量
func TestIdentGlobalSelectorStruct(t *testing.T) {
	// 测试内部Receiver
	path := "test/test_ident.go"
	funcName := "GlobalSelectorStruct"
	//relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
	variableName := "localVariable"
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}

			// Parse file and create the AST
			var fset = token.NewFileSet()
			var f *ast.File
			f, err = parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return nil
			}

			// 查找函数声明
			ast.Inspect(f, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
					if funcDecl.Name.Name == funcName {
						fmt.Printf("Function name: %s\n", funcDecl.Name.Name)
						// 查找函数声明
						for _, b := range funcDecl.Body.List {
							switch bType := b.(type) {
							case *ast.AssignStmt:
								if bType.Lhs[0].(*ast.Ident).Name == variableName {
									// 解析对应Rhs的内容
									identName := bType.Rhs[0].(*ast.Ident).Name
									// todo 这个方法是测试拿包下的局部变量，所以要额外处理下
									fmt.Print(identName)
								}
							}
						}
					}
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}
