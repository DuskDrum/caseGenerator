package stmt

import (
	"caseGenerator/parser/expr"
	"go/ast"

	"github.com/samber/lo"
)

// Decl 声明语句.
// 通常用于局部变量或常量的声明
type Decl struct {
	expr.ValueSpec
}

func (d *Decl) Express() []StatementExpression {
	//TODO 待补充

	return nil
}

func ParseDecl(stmt *ast.DeclStmt) *Decl {
	decl := &Decl{}
	switch declType := stmt.Decl.(type) {
	case *ast.GenDecl:
		for _, s := range declType.Specs {
			switch specType := s.(type) {
			case *ast.ValueSpec:
				spec := expr.ParseValueSpec(specType)
				if spec != nil {
					decl.ValueSpec = lo.FromPtr(spec)
				}
				return decl
			}
		}
	}
	return decl
}
