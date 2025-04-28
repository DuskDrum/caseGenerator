package expr

import (
	"caseGenerator/common/utils"
	"caseGenerator/parser/bo"
	_struct "caseGenerator/parser/struct"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ParseParameter 同时处理几种可能存在值的类型，如BasicLit、FuncLit、CompositeLit、CallExpr
// 得放在这个包中，不然会导致循环依赖
func ParseParameter(expr ast.Expr, context bo.ExprContext) _struct.Parameter {
	if expr == nil {
		return nil
	}
	switch exprType := expr.(type) {
	case *ast.SelectorExpr:
		return ParseSelector(exprType, context)
	case *ast.Ident:
		return ParseIdent(exprType, context)
		// 指针类型
	case *ast.StarExpr:
		return ParseStar(exprType, context)
	case *ast.FuncType:
		return ParseFuncType(exprType, context)
	case *ast.InterfaceType:
		return ParseInterface(exprType, context)
	case *ast.ArrayType:
		return ParseArray(exprType, context)
	case *ast.MapType:
		return ParseMap(exprType, context)
	case *ast.Ellipsis:
		return ParseEllipsis(exprType, context)
	case *ast.ChanType:
		return ParseChan(exprType, context)
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return ParseIndex(exprType, context)
	case *ast.BasicLit:
		return ParseBasicLit(exprType, context)
		// FuncLit 等待解析出内容值
	case *ast.FuncLit:
		return ParseFuncLit(exprType, context)
	case *ast.CompositeLit:
		return ParseCompositeLit(exprType, context)
		// CallExpr 等待解析出内容值
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		return ParseCall(exprType, context)
	case *ast.KeyValueExpr:
		return ParseKeyValue(exprType, context)
		// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		return ParseUnary(exprType, context)
	case *ast.BinaryExpr:
		return ParseBinary(exprType, context)
	case *ast.ParenExpr:
		return ParseParent(exprType, context)
	case *ast.StructType:
		return ParseStruct(exprType, context)
	case *ast.SliceExpr:
		return ParseSlice(exprType, context)
	case *ast.TypeAssertExpr:
		return ParseTypeAssert(exprType, context)
	//case *ast.FuncDecl:
	//	return ParseFuncDecl(exprType)
	default:
		panic("未知类型...")
	}
}

// ParseRecursionValue 解析递归的 value
func ParseRecursionValue(expr ast.Expr, context bo.ExprContext) *_struct.RecursionParam {
	parameter := ParseParameter(expr, context)
	ap := &_struct.RecursionParam{
		Parameter: parameter,
	}
	switch exprType := expr.(type) {
	case *ast.ArrayType:
		ap.Child = ParseRecursionValue(exprType.Elt, context)
	case *ast.MapType:
		ap.Child = ParseRecursionValue(exprType.Value, context)
	}
	return ap
}

// ParseCallExprResponse 解析Call类型的响应列表类型
// 本包的调用: func或者Func
// 其他包的调用: xx.Func
// rece的调用: xx.Func
// 本包内的直接调用: Func、func
func ParseCallExprResponse(importName, funcName string, context bo.ExprContext) ([]Field, error) {
	// 先判断如果没有包类型，那么在本包内找对应的方法
	if importName == "" {
		return ParseInnerCallResponse(funcName, context)
	}

	// 再找import去解析对应的方法
	result, err := ParsePackageCallResponse(importName, funcName, context)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	// 再找receiver的调用, 先找到receiver的类型，然后按照包名找到receiver的位置，最后找到方法
	response, err := ParseCallExprReceiverResponse(importName, funcName, context)
	if err != nil {
		return nil, err
	}
	if len(response) > 0 {
		return response, nil
	}

	return nil, nil
}

// ParseInnerCallResponse 解析包方法的调用响应列表
func ParseInnerCallResponse(funcName string, context bo.ExprContext) ([]Field, error) {
	relPackagePath := context.RealPackPath
	// 1. 使用os.ReadDir遍历相对package路径
	entries, err := os.ReadDir(relPackagePath)
	if err != nil {
		return nil, err
	}
	fields := make([]Field, 0, 10)
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
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				if funcDecl.Recv != nil {
					return true
				}
				if funcDecl.Name.Name == funcName {
					funcType := ParseFuncType(funcDecl.Type, context)
					fields = append(fields, funcType.Results...)
					return false
				}
			}
			return true
		})
	}
	return fields, nil
}

// ParsePackageCallResponse 解析包方法的调用响应列表
func ParsePackageCallResponse(importName, funcName string, context bo.ExprContext) ([]Field, error) {
	af := context.AstFile
	if importName == "" {
		return nil, errors.New("import name is blank")
	}
	if af == nil {
		return nil, errors.New("ast file is nil")
	}
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
	if importPath == "" {
		// 这是receiver调用的场景
		return nil, nil
	}

	path := strings.Trim(importPath, "\"")
	importCtx := build.Default
	pkg, _ := importCtx.Import(path, "", build.FindOnly)
	fmt.Println(pkg.Dir) // 输出包源码目录

	fields := make([]Field, 0, 10)

	fset := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fset, pkg.Dir, nil, parser.ParseComments)
	for _, pkgAst := range pkgs {
		ast.Inspect(pkgAst, func(n ast.Node) bool {
			// 检查是否为 *ast.BasicLit 节点
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == funcName {
					funcType := ParseFuncType(funcDecl.Type, context)
					fields = append(fields, funcType.Results...)
					return false
				}
			}
			return true
		})
	}
	return fields, nil
}

// ParseCallExprReceiverResponse 解析receiver方法的响应列表
// receiveName: 调用方名称
// funcName: 方法名
// relativeFilePath: 文件对应的相对路径
func ParseCallExprReceiverResponse(receiveName, funcName string, context bo.ExprContext) ([]Field, error) {
	// 1. 参数校验
	if receiveName == "" {
		return nil, errors.New("import name is blank")
	}
	// 2. 判断fd的receiver是否也是receiveName
	response, err := ParseInnerReceiverResponse(receiveName, funcName, context)
	if err != nil {
		return nil, err
	}
	if len(response) > 0 {
		return response, nil
	}
	// 3. 判断receiver是否是var定义的
	receiverResponse, err := ParseVariableReceiverResponse(receiveName, funcName, context)
	if err != nil {
		return nil, err
	}
	return receiverResponse, nil
}

// ParseInnerReceiverResponse 解析内部的receiver互相调用方法
func ParseInnerReceiverResponse(receiveName, funcName string, context bo.ExprContext) ([]Field, error) {
	relPackagePath := context.RealPackPath
	fd := context.AstFuncDecl
	// 1. 使用os.ReadDir遍历相对package路径
	entries, err := os.ReadDir(relPackagePath)
	if err != nil {
		return nil, err
	}
	fields := make([]Field, 0, 10)
	// 1. 如果是本receiver的方法，那么遍历本包
	if fd.Recv != nil {
		for _, v := range fd.Recv.List {
			if v.Names[0].Name == receiveName {
				var receiverType string
				switch exprType := v.Type.(type) {
				case *ast.Ident:
					receiverType = exprType.Name
				case *ast.StarExpr:
					receiverType = exprType.X.(*ast.Ident).Name
				}
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
						if funcDecl, ok := n.(*ast.FuncDecl); ok {
							if funcDecl.Name.Name == funcName {
								if funcDecl.Recv == nil {
									return true
								}
								for _, v := range funcDecl.Recv.List {
									if identV, typeOk := v.Type.(*ast.Ident); typeOk && identV.Name == receiverType {
										funcType := ParseFuncType(funcDecl.Type, context)
										fields = append(fields, funcType.Results...)
										return false
									}
									if identV, typeOk := v.Type.(*ast.StarExpr); typeOk && identV.X.(*ast.Ident).Name == receiverType {
										funcType := ParseFuncType(funcDecl.Type, context)
										fields = append(fields, funcType.Results...)
										return false
									}
								}
							}
						}
						return true
					})
				}
			}
		}
	}
	return fields, nil
}

// ParseVariableReceiverResponse 解析receiver为变量的场景
func ParseVariableReceiverResponse(receiveName, funcName string, context bo.ExprContext) ([]Field, error) {
	relPackagePath := context.RealPackPath
	af := context.AstFile
	entries, err := os.ReadDir(relPackagePath)
	if err != nil {
		return nil, err
	}
	fields := make([]Field, 0, 10)

	for _, entry := range entries {
		fullPath := filepath.Join(relPackagePath, entry.Name())
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		var fset = token.NewFileSet()
		f, _ := parser.ParseFile(fset, fullPath, nil, parser.ParseComments)
		ast.Inspect(f, func(n ast.Node) bool {
			if decl, ok := n.(*ast.GenDecl); ok && decl.Tok == token.VAR {
				for _, spec := range decl.Specs {
					if vspec, ok := spec.(*ast.ValueSpec); ok {
						for i, name := range vspec.Names {
							if name.Name == receiveName {
								// 一般是用new出来的局部变量，先处理new
								typeValue := vspec.Values[i]
								switch tv := typeValue.(type) {
								case *ast.CallExpr:
									if tvname, tvnameOk := tv.Fun.(*ast.Ident); tvnameOk {
										if tvname.Name == "new" {
											switch tvArg := tv.Args[0].(type) {
											case *ast.Ident: // receiver是本包
												receiverName := tvArg.Name // 在本包里找
												for _, receiverEntry := range entries {
													receiverFullPath := filepath.Join(relPackagePath, receiverEntry.Name())
													if receiverEntry.IsDir() || filepath.Ext(receiverEntry.Name()) != ".go" {
														continue
													}
													var receiverFset = token.NewFileSet()
													receiverF, _ := parser.ParseFile(receiverFset, receiverFullPath, nil, parser.ParseComments)
													ast.Inspect(receiverF, func(n ast.Node) bool {
														if n == nil {
															return true
														}
														// 检查是否为 *ast.BasicLit 节点
														if funcDecl, ok := n.(*ast.FuncDecl); ok {
															if funcDecl.Name.Name == funcName {
																if funcDecl.Recv == nil {
																	return true
																}
																for _, v := range funcDecl.Recv.List {
																	if identV, typeOk := v.Type.(*ast.Ident); typeOk && identV.Name == receiverName {
																		funcType := ParseFuncType(funcDecl.Type, context)
																		fields = append(fields, funcType.Results...)
																		return false
																	}
																	if identV, typeOk := v.Type.(*ast.StarExpr); typeOk && identV.X.(*ast.Ident).Name == receiverName {
																		funcType := ParseFuncType(funcDecl.Type, context)
																		fields = append(fields, funcType.Results...)
																		return false
																	}
																}
															}
														}
														return true
													})
												}
											case *ast.SelectorExpr: // receiver是别的包里的
												importPath := ""
												tvArgReceiver, tvArgReceiverOk := tvArg.X.(*ast.Ident)
												if !tvArgReceiverOk {
													return true
												}
												tvArgFuncName := tvArg.Sel.Name
												for _, importSpec := range af.Imports {
													if importSpec.Name == nil {
														suffixAfterDot := utils.GetSuffixAfterDot(importSpec.Path.Value)
														if tvArgReceiver.Name == suffixAfterDot {
															importPath = importSpec.Path.Value
														}
													} else {
														if tvArgReceiver.Name == importSpec.Name.Name {
															importPath = importSpec.Path.Value
														}
													}
												}
												if importPath == "" {
													panic("can not find import path")
												}

												path := strings.Trim(importPath, "\"")
												importCtx := build.Default
												pkg, _ := importCtx.Import(path, "", build.FindOnly)
												fmt.Println(pkg.Dir) // 输出包源码目录

												tfset := token.NewFileSet()
												pkgs, _ := parser.ParseDir(tfset, pkg.Dir, nil, parser.ParseComments)
												for _, pkgAst := range pkgs {
													ast.Inspect(pkgAst, func(n ast.Node) bool {
														// 检查是否为 *ast.BasicLit 节点
														if funcDecl, ok := n.(*ast.FuncDecl); ok {
															if funcDecl.Name.Name == funcName {
																if funcDecl.Recv == nil {
																	return true
																}
																for _, v := range funcDecl.Recv.List {
																	if identV, typeOk := v.Type.(*ast.Ident); typeOk && identV.Name == tvArgFuncName {
																		funcType := ParseFuncType(funcDecl.Type, context)
																		fields = append(fields, funcType.Results...)
																		return false
																	}
																	if identV, typeOk := v.Type.(*ast.StarExpr); typeOk && identV.X.(*ast.Ident).Name == tvArgFuncName {
																		funcType := ParseFuncType(funcDecl.Type, context)
																		fields = append(fields, funcType.Results...)
																		return false
																	}
																}
															}
														}
														return true
													})
												}
											}
										}

									}
								}
							}

						}
					}
				}
			}
			return true
		})
	}
	return fields, nil
}
