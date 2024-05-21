package vistitor

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"fmt"
	"go/ast"
	"log"
	"sync"
)

type AssignmentVisitor struct {
	addMu         sync.Mutex
	typeAssertMap map[string]generate.RequestDetail
}

func (v *AssignmentVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	switch node := n.(type) {
	case *ast.AssignStmt:
		// 左边：变量，层级调用
		// 右边：类型断言，赋值，方法
		var ab bo.AssignmentBinary
		var adi bo.AssignmentDetailInfo
		adi.LeftName = make([]string, 0, 10)
		for _, nodeLhs := range node.Lhs {
			switch nLhsType := nodeLhs.(type) {
			case *ast.Ident:
				name := nLhsType.Name
				if name == "_" {
					continue
				}
				adi.LeftName = append(adi.LeftName, nLhsType.Name)
			case *ast.SelectorExpr:
				adi.LeftName = append(adi.LeftName, GetRelationFromSelectorExpr(nLhsType))
			}
		}

		switch nRhsType := node.Rhs[0].(type) {
		case *ast.CallExpr:
			switch callFunType := nRhsType.Fun.(type) {
			case *ast.Ident:
				ab.Y = &bo.ParamUnary{callFunType.Name}
			case *ast.SelectorExpr:
				ab.Y = &bo.ParamUnary{GetRelationFromSelectorExpr(callFunType)}
			default:
				log.Fatalf("不支持此类型")
			}
			// 类型断言可以是 a.(type) 也可以是A.B.C.(type)
		case *ast.TypeAssertExpr:
			// 类型断言已在上面处理了
		case *ast.UnaryExpr:
			if se, ok := nRhsType.X.(*ast.CompositeLit); ok {
				ab.Y = CompositeLitParse(se)
			}
			if ident, ok := nRhsType.X.(*ast.Ident); ok {
				ab.Y = &bo.ParamUnary{ident.Name}
			}
		case *ast.CompositeLit:
			ab.Y = CompositeLitParse(nRhsType)
		case *ast.BasicLit:
		// 基本字面值,数字或者字符串。跳过不解析
		case *ast.FuncLit:
			funcType := nRhsType.Type
			paramType := parseFuncType(funcType)
			adi.RightType = "function"
			adi.RightFormula = paramType
			adi.RightEmptyValue = "nil"
		// 内部函数，暂时不处理
		default:
			log.Fatalf("不支持此类型")
		}
		bo.AppendAssignmentDetailInfoToList(adi)
	case *ast.DeclStmt:
		switch nd := node.Decl.(type) {
		case *ast.GenDecl:
			for _, ndSpec := range nd.Specs {
				switch npVa := ndSpec.(type) {
				case *ast.ValueSpec:
					for _, npVaName := range npVa.Names {
						if npVaName.Name == "_" {
							continue
						}
						var ab bo.AssignmentBinary
						ab.X = bo.ParamUnary{npVaName.Name}
						if npVa.Type == nil {
							continue
						}
						switch vaType := npVa.Type.(type) {
						case *ast.Ident:
							ab.Y = &bo.ParamUnary{vaType.Name}
						case *ast.FuncType:
							// 空的
							ab.Y = &bo.FuncUnary{}
						case *ast.SelectorExpr:
							ab.Y = &bo.ParamUnary{GetRelationFromSelectorExpr(vaType)}
						default:
							log.Fatalf("类型不支持")
						}
						bo.AddParamNeedToMapDetail(ab.X.ParamValue, &BinaryParam{
							ParamName:   ab.X.ParamValue,
							BinaryParam: &ab,
						})
					}
				case *ast.TypeSpec:
					fmt.Println("进来了")
				}
			}
		default:
			log.Fatalf("不支持此类型")
		}
	}

	return v
}
