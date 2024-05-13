package vistitor

import (
	"caseGenerator/parse/bo"
	"go/ast"
)

// ReceiverVisitor receive
type ReceiverVisitor struct {
}

func (v *ReceiverVisitor) Visit(n ast.Node) ast.Visitor {
	decl, ok := n.(*ast.FuncDecl)
	if !ok {
		return v
	}
	if decl.Recv == nil {
		return v
	}
	if len(decl.Recv.List) == 0 {
		return v
	}
	var recvName string
	switch recvType := decl.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		switch astStartExpr := recvType.X.(type) {
		case *ast.Ident:
			recvName = astStartExpr.Name
		// 下标类型，实际上是泛型。泛型先不处理
		case *ast.IndexExpr:
			return v
		}
	case *ast.Ident:
		recvName = recvType.Name
	default:
		recvName = decl.Recv.List[0].Names[0].Name
	}
	rec := bo.ReceiverInfo{
		ReceiverName: recvName,
	}
	switch typeType := decl.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		switch astStartExpr := typeType.X.(type) {
		case *ast.Ident:
			rec.ReceiverValue = bo.InvocationUnary{InvocationName: astStartExpr.Name}
		// 下标类型，实际上是泛型。泛型先不处理
		case *ast.IndexExpr:
			return v
		}
	case *ast.Ident:
		rec.ReceiverValue = bo.InvocationUnary{InvocationName: typeType.Name}
	}
	bo.AddParamNeedToMapDetail(recvName, &rec)
	bo.SetReceiverInfo(&rec)
	return v
}
