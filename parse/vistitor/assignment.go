package vistitor

import (
	"caseGenerator/common/enum"
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"fmt"
	"go/ast"
	"sync"
)

type AssignmentVisitor struct {
	addMu         sync.Mutex
	typeAssertMap map[string]generate.CaseRequest
}

func (v *AssignmentVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	var adi bo.AssignmentDetailInfo
	switch node := n.(type) {
	case *ast.AssignStmt:
		// 左边：变量，层级调用
		// 右边：类型断言，赋值，方法
		adi.LeftName = make([]string, 0, 10)
		for _, nodeLhs := range node.Lhs {
			// 有一些赋值这一句无法判断出等式右边变量的类型，只能靠变量名上下文联系
			parseResult := ParseParamWithoutInit(nodeLhs, "")
			adi.LeftName = append(adi.LeftName, parseResult.ParamType)
		}

		switch nRhsType := node.Rhs[0].(type) {
		case *ast.CallExpr:
			paramRequests := make([]bo.ParamParseRequest, 0, 10)
			leftResults := make([]string, 0, 10)
			// 解析call的方法入参参数
			for _, v := range nRhsType.Args {
				parseResult := ParseParamRequest(v)
				for _, s := range parseResult {
					paramRequests = append(paramRequests, *s)
				}
			}
			adi.RightFunctionParam = paramRequests
			// 解析call的方法的出参类型
			switch callFunType := nRhsType.Fun.(type) {
			case *ast.Ident:
				adi.RightType = enum.ASSIGNMENT_TYPE_CALL
				adi.RightFormula = callFunType.Name
				// 解析返回的类型
				switch callDecl := callFunType.Obj.Decl.(type) {
				case *ast.Field:
					switch callDeclType := callDecl.Type.(type) {
					case *ast.FuncType:
						_, result := ParseFuncTypeParamParseResult(callDeclType)
						for _, v := range result {
							leftResults = append(leftResults, v.ParamType)
						}
						adi.LeftResultType = leftResults
					}
				}
			case *ast.SelectorExpr:
				adi.RightType = enum.ASSIGNMENT_TYPE_CALL
				adi.RightFormula = GetRelationFromSelectorExpr(callFunType)
				// SelectorExpr类型解析不出类型，只有靠后面的逻辑猜测
			default:
				panic("不支持此类型")
			}
			// 类型断言可以是 a.(type) 也可以是A.B.C.(type)
		case *ast.TypeAssertExpr:
			// 类型断言已在上面处理了
		case *ast.UnaryExpr:
			paramRequests := make([]bo.ParamParseRequest, 0, 10)
			parseResult := ParseParamRequest(nRhsType.X)
			for _, s := range parseResult {
				paramRequests = append(paramRequests, *s)
			}
			adi.RightFunctionParam = paramRequests

			if se, ok := nRhsType.X.(*ast.CompositeLit); ok {
				adi.RightType = enum.ASSIGNMENT_TYPE_COMPOSITE
				adi.RightFormula = CompositeLitParse(se).ResultStructName
				adi.LeftResultType = []string{"nil"}
			}
			if ident, ok := nRhsType.X.(*ast.Ident); ok {
				adi.RightType = enum.ASSIGNMENT_TYPE_COMPOSITE
				adi.RightFormula = ident.Name
				adi.LeftResultType = []string{"nil"}
			}
		// 构造类型
		case *ast.CompositeLit:
			paramRequests := make([]bo.ParamParseRequest, 0, 10)
			parseResult := ParseParamRequest(nRhsType)
			for _, s := range parseResult {
				paramRequests = append(paramRequests, *s)
			}
			adi.RightFunctionParam = paramRequests

			adi.RightFormula = CompositeLitParse(nRhsType).ResultStructName
			adi.LeftResultType = []string{"nil"}
		case *ast.BasicLit:
		// 基本字面值,数字或者字符串。跳过不解析
		case *ast.FuncLit:
			paramRequests := make([]bo.ParamParseRequest, 0, 10)
			parseResult := ParseParamRequest(nRhsType)
			for _, s := range parseResult {
				paramRequests = append(paramRequests, *s)
			}
			adi.RightFunctionParam = paramRequests

			funcType := nRhsType.Type
			paramType := parseFuncType(funcType)
			adi.RightFormula = paramType
			// 匿名函数
			adi.RightType = enum.ASSIGNMENT_TYPE_FUNCTION
			adi.LeftResultType = []string{enum.ASSIGNMENT_TYPE_FUNCTION.Code}
		default:
			panic("不支持此类型")
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
							panic("类型不支持")
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
			panic("不支持此类型")
		}
		bo.AppendAssignmentDetailInfoToList(adi)
	// 这种是没有响应值的function
	case *ast.ExprStmt:
		switch nd := node.X.(type) {
		case *ast.CallExpr:
			// 没有响应值的function，没有响应信息
			adi.LeftResultType = []string{}

			paramRequests := make([]bo.ParamParseRequest, 0, 10)
			for _, v := range nd.Args {
				parseResult := ParseParamRequest(v)
				for _, s := range parseResult {
					paramRequests = append(paramRequests, *s)
				}
			}
			adi.RightFunctionParam = paramRequests
			// 解析右边的方法
			init := ParseParamWithoutInit(nd.Fun, "")
			adi.RightType = enum.ASSIGNMENT_TYPE_CALL
			adi.RightFormula = init.ParamType
		default:
			panic("不支持此类型")
		}
		bo.AppendAssignmentDetailInfoToList(adi)
	}

	return v
}
