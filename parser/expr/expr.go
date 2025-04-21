package expr

import (
	"caseGenerator/common/utils"
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
func ParseParameter(expr ast.Expr, af *ast.File) _struct.Parameter {
	if expr == nil {
		return nil
	}
	switch exprType := expr.(type) {
	case *ast.SelectorExpr:
		return ParseSelector(exprType, af)
	case *ast.Ident:
		return ParseIdent(exprType, af)
		// 指针类型
	case *ast.StarExpr:
		return ParseStar(exprType, af)
	case *ast.FuncType:
		return ParseFuncType(exprType, af)
	case *ast.InterfaceType:
		return ParseInterface(exprType, af)
	case *ast.ArrayType:
		return ParseArray(exprType, af)
	case *ast.MapType:
		return ParseMap(exprType, af)
	case *ast.Ellipsis:
		return ParseEllipsis(exprType, af)
	case *ast.ChanType:
		return ParseChan(exprType, af)
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return ParseIndex(exprType, af)
	case *ast.BasicLit:
		return ParseBasicLit(exprType, af)
		// FuncLit 等待解析出内容值
	case *ast.FuncLit:
		return ParseFuncLit(exprType, af)
	case *ast.CompositeLit:
		return ParseCompositeLit(exprType, af)
		// CallExpr 等待解析出内容值
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		return ParseCall(exprType, af)
	case *ast.KeyValueExpr:
		return ParseKeyValue(exprType, af)
		// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		return ParseUnary(exprType, af)
	case *ast.BinaryExpr:
		return ParseBinary(exprType, af)
	case *ast.ParenExpr:
		return ParseParent(exprType, af)
	case *ast.StructType:
		return ParseStruct(exprType, af)
	case *ast.SliceExpr:
		return ParseSlice(exprType, af)
	case *ast.TypeAssertExpr:
		return ParseTypeAssert(exprType, af)
	//case *ast.FuncDecl:
	//	return ParseFuncDecl(exprType)
	default:
		panic("未知类型...")
	}
}

// ParseRecursionValue 解析递归的 value
func ParseRecursionValue(expr ast.Expr, af *ast.File) *_struct.RecursionParam {
	parameter := ParseParameter(expr, af)
	ap := &_struct.RecursionParam{
		Parameter: parameter,
	}
	switch exprType := expr.(type) {
	case *ast.ArrayType:
		ap.Child = ParseRecursionValue(exprType.Elt, af)
	case *ast.MapType:
		ap.Child = ParseRecursionValue(exprType.Value, af)
	}
	return ap
}

// ParseCallExprResponse 解析Call类型的响应列表类型
// 其他包的调用: xx.Func
// rece的调用: xx.Func
// 本包内的直接调用: Func、func
func ParseCallExprResponse(importName, funcName string, af *ast.File) []Field {
	// 先找import去解析对应的方法
	result, err := ParseCallExprImportResponse(importName, funcName, af)
	if err == nil {
		return result
	}
	// 再找receiver的调用, 先找到receiver的类型，然后按照包名找到receiver的位置，最后找到方法

	// 如果没有包类型，那么在本包内找对应的方法
	//af.Package
	return nil
}

// ParseCallExprImportResponse 解析import方法的响应列表
func ParseCallExprImportResponse(importName, funcName string, af *ast.File) ([]Field, error) {
	if importName == "" {
		return nil, errors.New("import name is blank")
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
		panic("can not find import path")
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
					funcType := ParseFuncType(funcDecl.Type, nil)
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
func ParseCallExprReceiverResponse(receiveName, funcName, relPackagePath string, af *ast.File) ([]Field, error) {
	if receiveName == "" {
		return nil, errors.New("import name is blank")
	}
	// 1. 使用os.ReadDir遍历相对package路径
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
														// 检查是否为 *ast.BasicLit 节点
														if funcDecl, ok := n.(*ast.FuncDecl); ok {
															if funcDecl.Name.Name == funcName {
																for _, v := range funcDecl.Recv.List {
																	if identV, typeOk := v.Type.(*ast.Ident); typeOk && identV.Name == receiverName {
																		funcType := ParseFuncType(funcDecl.Type, nil)
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
												for _, importSpec := range af.Imports {
													if importSpec.Name == nil {
														suffixAfterDot := utils.GetSuffixAfterDot(importSpec.Path.Value)
														if receiveName == suffixAfterDot {
															importPath = importSpec.Path.Value
														}
													} else {
														if receiveName == importSpec.Name.Name {
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
																funcType := ParseFuncType(funcDecl.Type, nil)
																fields = append(fields, funcType.Results...)
																return false
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
