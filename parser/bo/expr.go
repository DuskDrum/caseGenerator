package bo

import "go/ast"

// ExprContext  expr上下文
type ExprContext struct {
	AstFile      *ast.File     // 文件的ast详情
	AstFuncDecl  *ast.FuncDecl // 本方法的ast详情
	RealPackPath string        // 包的绝对路径
}
