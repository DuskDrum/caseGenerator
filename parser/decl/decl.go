package decl

import (
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// ParseDecl 源代码语法结构的包
func ParseDecl(expr ast.Decl) _struct.Parameter {
	switch exprType := expr.(type) {
	case *ast.FuncDecl:
		ParseFunc(exprType)
	//case *ast.GenDecl:
	//	return ParseIdent(exprType)

	//case *ast.FuncDecl:
	//	return ParseFuncDecl(exprType)
	default:
		panic("未知类型...")
	}

	return nil
}
