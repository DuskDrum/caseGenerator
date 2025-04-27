package expr

import (
	"caseGenerator/parser/bo"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseCallExprReceiverResponseCase1
// receiver是直接对应其他包的方法
func TestParseCallExprReceiverResponseCase1(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "add"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ExprStmt); eok {
								if c, cok := t.X.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 fmt， Sel 是 Print
										context := bo.ExprContext{
											AstFile:     f,
											AstFuncDecl: funcDecl,
										}
										response, err2 := ParsePackageCallResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase2
// receiver是本包内定义的，new出来的变量(在本包中)
func TestParseCallExprReceiverResponseCase2(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddReceiver1"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase3
// receiver是本包内定义的，new出来的变量(不在本包中，跨包)
func TestParseCallExprReceiverResponseCase3(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddPackageService"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase4
// receiver是本包内定义的，new出来的变量(在本包中)， 方法中用的指针定义
func TestParseCallExprReceiverResponseCase4(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddPtrReceiver1"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase5
// receiver是本包内定义的，new出来的变量(不在本包中，跨包)， 方法中用的指针定义
func TestParseCallExprReceiverResponseCase5(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddPtrPackageService"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase6
// receiver是此方法的receiver调用
func TestParseCallExprReceiverResponseCase6(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddReceiver"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase7
// receiver是此方法的receiver调用 (指针类型)
func TestParseCallExprReceiverResponseCase7(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "AddReceiverPtr"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"

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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.SelectorExpr); sok {
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										// X 是 receiver1， Sel 是 add
										response, err2 := ParseCallExprReceiverResponse(s.X.(*ast.Ident).Name, s.Sel.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase8
// 本包内方法直接调用
func TestInnerCallResponseCase8(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "innerCall"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.Ident); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprResponse("", s.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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

// TestParseCallExprReceiverResponseCase9
// 本包内方法直接调用
func TestInnerCallResponseCase9(t *testing.T) {
	path := "test/test_receiver_function_call.go"
	funcName := "InnerCall"
	relPackagePath := "/Users/wangyi/githubProject/caseGenerator/parser/expr/test"
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
						for _, stmt := range funcDecl.Body.List {
							// 判断是callExpr
							if t, eok := stmt.(*ast.ReturnStmt); eok {
								result := t.Results[0]
								if c, cok := result.(*ast.CallExpr); cok {
									// 找到Func
									if s, sok := c.Fun.(*ast.Ident); sok {
										// X 是 receiver1， Sel 是 add
										context := bo.ExprContext{
											AstFile:      f,
											AstFuncDecl:  funcDecl,
											RealPackPath: relPackagePath,
										}
										response, err2 := ParseCallExprResponse("", s.Name, context)
										if err2 != nil {
											fmt.Printf("Error: %s\n", err2)
										} else {
											fmt.Printf("Response: %v\n", response)
										}
									}
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
