package param

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestValueSpecCase(t *testing.T) {
	// 要解析的代码字符串，表示一个简单的index表达式
	// 使用 parser.ParseExpr 来解析表达式
	src := []byte(`
    	package test
		var myVariable, you int
		var x, y = 10, 20
		var a, v = "", s.GetInstance()
	`)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	// *ast.IndexExpr
	// *ast.ValueSpec
	for _, decl := range file.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && gendecl.Tok == token.VAR {
			for _, spec := range gendecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				for _, name := range valueSpec.Names {
					println("标识符名称:", name)
				}
			}
		}
	}
}
