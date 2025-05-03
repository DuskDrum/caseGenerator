package z3

import (
	"caseGenerator/common/enum"
	"caseGenerator/go-z3"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func ExpressIdent(param *expr.Ident, context bo.ExpressionContext) (*z3.AST, []*z3.AST) {
	var err error
	astList := make([]*z3.AST, 0, 10)
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	// 1. 根据上下文解析出变量名
	symbolName := context.TemporaryVariable.VariableName
	// 2. 判断这个变量的类型
	// 2.1 判断请求变量中有没有这个变量
	// 2.2 判断局部变量中有没有这个变量
	paramType := context.GetKnownParamType(symbolName)
	if paramType == enum.BASIC_PARAMETER_TYPE_UNKNOWN {
		// 2.3 判断ident对应的变量是否在请求变量中（ident 可能指向其他的变量）
		// 2.4 判断ident对应的变量是否在局部变量中（ident 可能指向其他的变量）
		paramType = context.GetKnownParamType(param.IdentName)
		if paramType == enum.BASIC_PARAMETER_TYPE_UNKNOWN {
			// 2.5 判断ident对应的变量是否在全局变量中 (ident可能是包内的其他变量)
			paramType, err = getGlobalVariableType(context)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	// 2.5 判断能否从这个ident的值推导出类型 (比如ident的值是: 8 那么就是int)
	// 这里需要注意2.5里ident的值有可能是全局变量，是个string

	// iType
	switch paramType {
	case enum.BASIC_PARAMETER_TYPE_INT: // rune就是int32
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.IntSort()))
	case enum.BASIC_PARAMETER_TYPE_FLOAT:
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.FloatSort()))
	case enum.BASIC_PARAMETER_TYPE_DOUBLE:
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.DoubleSort()))
	case enum.BASIC_PARAMETER_TYPE_BOOL:
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.BoolSort()))
	case enum.BASIC_PARAMETER_TYPE_STRING:
		astList = append(astList, ctx.Const(ctx.Symbol(symbolName), ctx.StringSort()))
	default:
		// 啥也不干
	}
	return nil, astList
}

func getGlobalVariableType(context bo.ExpressionContext) (enum.BasicParameterType, error) {
	// 解析此包里的内容,找是否存在此全局变量，如果有的话取其类型或者推导出其类型
	relPackagePath := context.RealPackPath
	variableName := context.TemporaryVariable.VariableName
	// 1. 使用os.ReadDir遍历相对package路径
	entries, err := os.ReadDir(relPackagePath)
	if err != nil {
		return enum.BASIC_PARAMETER_TYPE_UNKNOWN, err
	}
	paramType := enum.BASIC_PARAMETER_TYPE_UNKNOWN
	// 1. 如果是本receiver的方法，那么遍历本包
	for _, entry := range entries {
		fullPath := filepath.Join(relPackagePath, entry.Name())
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		var fset = token.NewFileSet()
		f, _ := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			// 检查是否为 *ast.BasicLit 节点
			if genDecl, ok := n.(*ast.GenDecl); ok {
				// 处理每个变量声明
				for _, spec := range genDecl.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for i, v := range valueSpec.Names {
						if v.Name == variableName {
							va := valueSpec.Values[i]
							switch vaType := va.(type) {
							case *ast.BasicLit:
								paramType = enum.GetParamTypeByBasicLit(vaType)
								return true
							case *ast.Ident:
								if vaType.Name == "true" || vaType.Name == "false" {
									paramType = enum.BASIC_PARAMETER_TYPE_BOOL
									return true
								}
							}
						}
					}
					//
					//// 提取类型和初始值
					//typeStr := exprToString(valueSpec.Type)
					//values := exprSliceToString(valueSpec.Values)
					//
					//// 处理每个变量名
					//for i, name := range valueSpec.Names {
					//	gVar := GlobalVar{
					//		Name: name.Name,
					//		Type: typeStr,
					//	}
					//	// 关联初始值（可能为多个）
					//	if i < len(values) {
					//		gVar.Value = values[i]
					//	}
					//	globalVars = append(globalVars, gVar)
					//}
				}
			}

			return true
		})
	}
	return paramType, nil
}
