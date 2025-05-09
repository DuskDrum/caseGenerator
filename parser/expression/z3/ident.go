package z3

import (
	"caseGenerator/common/enum"
	"caseGenerator/common/utils"
	"caseGenerator/go-z3"
	"caseGenerator/parser/bo"
	"caseGenerator/parser/expr"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
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
			paramType, err = getGlobalVariableType(context, param.IdentName, "")
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

// getGlobalVariableType 获取全局变量类型，找
func getGlobalVariableType(context bo.ExpressionContext, relPackagePath, variableName string) (enum.BasicParameterType, error) {
	// 解析此包里的内容,找是否存在此全局变量，如果有的话取其类型或者推导出其类型
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
								return false
							case *ast.Ident:
								if vaType.Name == "true" || vaType.Name == "false" {
									paramType = enum.BASIC_PARAMETER_TYPE_BOOL
									return false
								}
								// 如果是其他的ident类型，那么要去找到引用的变量类型

							case *ast.SelectorExpr:
								// 如果是selector类型，那么要去找对应的引用变量类型

							case *ast.CallExpr:
								//如果是call类型，那么解析出call的返回类型。一般只有一个返回值
							}
						}
					}

				}
			}

			return true
		})
	}
	return paramType, nil
}

func getIdentVariableType(context bo.ExpressionContext, ident *ast.Ident) (enum.BasicParameterType, error) {
	paramType := enum.BASIC_PARAMETER_TYPE_UNKNOWN

	// 1. 使用os.ReadDir遍历相对package路径
	entries, err := os.ReadDir(context.RealPackPath)
	if err != nil {
		return paramType, err
	}
	// 1. 找到本包内的ident，并判断它对应的值的类型
	// 2. 如果他的值的类型是基础类型，那么直接返回
	// 3. 如果他的值的类型是ident，那么递归次方法
	// 4. 如果他的值的类型是selector， 那么调用selector方法
	for _, entry := range entries {
		fullPath := filepath.Join(context.RealPackPath, entry.Name())
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
						if v.Name == ident.Name {
							va := valueSpec.Values[i]
							switch vaType := va.(type) {
							case *ast.BasicLit:
								paramType = enum.GetParamTypeByBasicLit(vaType)
								return false
							case *ast.Ident:
								if vaType.Name == "true" || vaType.Name == "false" {
									paramType = enum.BASIC_PARAMETER_TYPE_BOOL
									return false
								}
								// 如果是其他的ident类型，那么要去找到引用的变量类型
								variableType, err2 := getIdentVariableType(context, vaType)
								if err2 != nil {
									paramType = enum.BASIC_PARAMETER_TYPE_UNKNOWN
									return false
								}
								paramType = variableType
								return false
							case *ast.SelectorExpr:
								// 如果是selector类型，那么要去找对应的引用变量类型
							}
						}
					}

				}
			}

			return true
		})
	}
	return paramType, nil
}

// GetSelectorVariableType  找到selector对应的变量类型
func GetSelectorVariableType(context bo.ExpressionContext, selector *ast.SelectorExpr, realPackPath string) (enum.BasicParameterType, error) {
	paramType := enum.BASIC_PARAMETER_TYPE_UNKNOWN
	// 当前变量的绝对路径目录
	entries, err := os.ReadDir(realPackPath)
	if err != nil {
		return paramType, err
	}
	// 1. 根据selector的前缀，找到其对应值：a. 本包内的变量里的属性值 b. 其他包的变量 c. 其他包的变量的属性值
	xName := selector.X.(*ast.Ident).Name

	// 1.1 先在本包中找对应的变量

	// 2.2 再在import中找其他的变量

	// 2. 如果他值的类型是基础类型，那么直接返回
	// 3. 如果值的类型是ident，那么递归ident方法
	// 4. 如果值的类型是selector， 那么调用此方法
	for _, entry := range entries {
		fullPath := filepath.Join(realPackPath, entry.Name())
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
						if v.Name == xName {
							va := valueSpec.Values[i]
							switch vaType := va.(type) {
							case *ast.BasicLit:
								paramType = enum.GetParamTypeByBasicLit(vaType)
								return false
							case *ast.Ident:
								if vaType.Name == "true" || vaType.Name == "false" {
									paramType = enum.BASIC_PARAMETER_TYPE_BOOL
									return false
								}
								// 如果是其他的ident类型，那么要去找到引用的变量类型
								variableType, err2 := getIdentVariableType(context, vaType)
								if err2 != nil {
									paramType = enum.BASIC_PARAMETER_TYPE_UNKNOWN
									return false
								}
								paramType = variableType
								return false
							case *ast.SelectorExpr:
								// 如果是selector类型，那么要去找对应的引用变量类型
							}
						}
					}

				}
			}

			return true
		})
	}
	return paramType, nil
}

func getOuterPackageSelector(context bo.ExpressionContext, importName, realPackPath, variableName string) (enum.BasicParameterType, error) {
	importPath := getImportPath(context, importName)

	// 找不到import信息，说明是本包内的变量属性， 那么去找本包内的属性值
	if importName == "" {

	} else { // 如果找得到import信息，说明是跨包的调用
		path := strings.Trim(importPath, "\"")
		importCtx := build.Default
		pkg, _ := importCtx.Import(path, "", build.FindOnly)
		return getGlobalVariableType(context, pkg.Dir, variableName)

	}

	return enum.BASIC_PARAMETER_TYPE_UNKNOWN, nil
}

func getImportPath(context bo.ExpressionContext, importName string) string {
	af := context.AstFile
	importPath := ""
	for _, importSpec := range af.Imports {
		if importSpec.Name == nil {
			suffixAfterDot := utils.GetSuffixAfterDot(importSpec.Path.Value)
			if importName == suffixAfterDot {
				importPath = importSpec.Path.Value
			}
		} else {
			if importName == importSpec.Name.Name {
				importPath = importSpec.Path.Value
			}
		}
	}
	return importPath
}
