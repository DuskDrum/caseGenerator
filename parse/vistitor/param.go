package vistitor

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"caseGenerator/utils"
	"encoding/json"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"sync"
)

// Constant 常量
type Constant struct {
	ConstantName  string
	ConstantValue any
}

func (c *Constant) GetParamName() string {
	return c.ConstantName
}

// UnmarshalerInfo 不用处理
func (c *Constant) UnmarshalerInfo(_ string) {
}

type Request struct {
	RequestName  string
	RequestValue any
	// 使用typeOf
	RequestType reflect.Type
}

func (r *Request) GetParamName() string {
	return r.RequestName
}

func (r *Request) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		r.RequestValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

// Variable 变量
type Variable struct {
	VariableName  string
	VariableValue any
	// 使用typeOf
	VariableType reflect.Type
}

func (v *Variable) GetParamName() string {
	return v.VariableName
}

func (v *Variable) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		v.VariableValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

// FuncReturn 方法返回的参数
type FuncReturn struct {
	FuncReturnName  string
	FuncReturnValue any
	FuncReturnIndex int

	// 调用方法的信息
	FuncImport      string
	FuncPackageName string
	FuncName        string
}

func (f *FuncReturn) GetParamName() string {
	return f.FuncReturnName
}

// UnmarshalerInfo 方法返回的结果，是没有类型信息的，需要想办法获取到
func (f *FuncReturn) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[any]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		f.FuncReturnValue = dat
	} else {
		fmt.Println(jsonString)
	}
}

type BinaryParam struct {
	ParamName   string
	BinaryParam bo.Binary
}

func (b *BinaryParam) GetParamName() string {
	return b.ParamName
}

func (b *BinaryParam) UnmarshalerInfo(_ string) {

}

type ParamVisitor struct {
	addMu         sync.Mutex
	typeAssertMap map[string]generate.RequestDetail
}

func (v *ParamVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	switch node := n.(type) {
	case *ast.AssignStmt:
		for _, nodeLhs := range node.Lhs {
			// 左边：变量，层级调用
			// 右边：类型断言，赋值，方法
			var ab bo.AssignmentBinary
			switch nLhsType := nodeLhs.(type) {
			case *ast.Ident:
				name := nLhsType.Name
				if name == "_" {
					continue
				}
				ab.X = bo.ParamUnary{nLhsType.Name}
			case *ast.SelectorExpr:
				ab.X = bo.ParamUnary{GetRelationFromSelectorExpr(nLhsType)}
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
				var tau bo.TypeAssertUnary
				switch tae := nRhsType.X.(type) {
				case *ast.Ident:
					tau.ParamValue = tae.Name
				case *ast.SelectorExpr:
					tau.ParamValue = GetRelationFromSelectorExpr(tae)
				default:
					log.Fatalf("不支持此类型")
				}
				switch nr := nRhsType.Type.(type) {
				case *ast.Ident:
					tau.AssertType = nr.Name
				case *ast.SelectorExpr:
					tau.AssertType = GetRelationFromSelectorExpr(nr)
				default:
					log.Fatalf("不支持此类型")
				}
				ab.Y = &tau
			case *ast.UnaryExpr:
				if se, ok := nRhsType.X.(*ast.CompositeLit); ok {
					ab.Y = CompositeLitParse(se)
				}
				if ident, ok := nRhsType.X.(*ast.Ident); ok {
					ab.Y = &bo.ParamUnary{ident.Name}
				}
			case *ast.CompositeLit:
				ab.Y = CompositeLitParse(nRhsType)
			default:
				log.Fatalf("不支持此类型")
			}
			bo.AddParamNeedToMapDetail(ab.X.ParamValue, &BinaryParam{
				ParamName:   ab.X.ParamValue,
				BinaryParam: &ab,
			})
		}
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
