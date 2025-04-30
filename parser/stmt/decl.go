package stmt

import (
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"go/ast"

	"github.com/samber/lo"
)

// Decl 声明语句.
// 通常用于局部变量或常量的声明
type Decl struct {
	expr.ValueSpec
}

func ParseDecl(stmt *ast.DeclStmt, context bo.ExprContext) *Decl {
	decl := &Decl{}
	switch declType := stmt.Decl.(type) {
	case *ast.GenDecl:
		for _, s := range declType.Specs {
			switch specType := s.(type) {
			case *ast.ValueSpec:
				spec := expr.ParseValueSpec(specType, context)
				if spec != nil {
					decl.ValueSpec = lo.FromPtr(spec)
				}
				return decl
			}
		}
	}
	return decl
}
