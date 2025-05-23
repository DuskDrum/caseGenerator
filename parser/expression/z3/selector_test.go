package z3

import (
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
func TestSelectorZ3Case(t *testing.T) {
	// 测试内部Receiver
	path := "test/selector.go"
	_ = "/Users/wangyi/githubProject/caseGenerator/parser/expression/z3/test"
	funcName := "Selector"
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
					if funcDecl.Name.Name == funcName {
						for _, b := range funcDecl.Body.List {
							switch stmt := b.(type) {
							case *ast.AssignStmt:
								stmt.Pos()
							}

						}
					}
					//// 查找函数声明
					//ast.Inspect(f, func(n ast.Node) bool {
					//	// 检查是否为 *ast.CallExpr 节点
					//	if callExpr, ok := n.(*ast.SelectorExpr); ok {
					//		fmt.Println("找到函数调用表达式:")
					//		context := bo.ExprContext{
					//			AstFile:      f,
					//			AstFuncDecl:  funcDecl,
					//			RealPackPath: relPackagePath,
					//		}
					//		// 解析出了call
					//		call := expr.ParseSelector(callExpr, context)
					//		// 执行z3处理器
					//		eContext := bo.ExpressionContext{
					//			VariableParamMap: make(map[string]enum.BasicParameterType, 10),
					//			RequestParamMap:  make(map[string]enum.BasicParameterType, 10),
					//			TemporaryVariable: bo.TemporaryVariable{
					//				VariableName: "localVariable",
					//			},
					//			ExprContext: context,
					//		}
					//		expressCall, _ := GetSelectorVariableType(eContext, callExpr, relPackagePath)
					//
					//		fmt.Print(expressCall)
					//
					//		fmt.Print(call)
					//	}
					//	return true
					//})

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
