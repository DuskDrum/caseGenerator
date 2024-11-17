package expr

import (
	_struct "caseGenerator/parser/struct"
	"go/ast"
)

// ParseParameter 同时处理几种可能存在值的类型，如BasicLit、FuncLit、CompositeLit、CallExpr
// 得放在这个包中，不然会导致循环依赖
func ParseParameter(expr ast.Expr) _struct.Parameter {
	switch exprType := expr.(type) {
	case *ast.SelectorExpr:
		return ParseSelector(exprType)
	case *ast.Ident:
		return ParseIdent(exprType)
		// 指针类型
	case *ast.StarExpr:
		return ParseStar(exprType)
	case *ast.FuncType:
		return ParseFuncType(exprType)
	case *ast.InterfaceType:
		return ParseInterface(exprType)
	case *ast.ArrayType:
		return ParseArray(exprType)
	case *ast.MapType:
		return ParseMap(exprType)
	case *ast.Ellipsis:
		return ParseEllipsis(exprType)
	case *ast.ChanType:
		return ParseChan(exprType)
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return ParseIndex(exprType)
	case *ast.BasicLit:
		return ParseBasicLit(exprType)
		// FuncLit 等待解析出内容值
	case *ast.FuncLit:
		return ParseFuncLit(exprType)
	case *ast.CompositeLit:
		return ParseCompositeLit(exprType)
		// CallExpr 等待解析出内容值
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		return ParseCall(exprType)
	case *ast.KeyValueExpr:
		return ParseKeyValue(exprType)
		// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		return ParseUnary(exprType)
	case *ast.BinaryExpr:
		return ParseBinary(exprType)
	case *ast.ParenExpr:
		return ParseParent(exprType)
	case *ast.StructType:
		return ParseStruct(exprType)
	case *ast.SliceExpr:
		return ParseSlice(exprType)
	//case *ast.FuncDecl:
	//	return ParseFuncDecl(exprType)
	default:
		panic("未知类型...")
	}
}

// ParseRecursionValue 解析递归的 value
func ParseRecursionValue(expr ast.Expr) *_struct.RecursionParam {
	parameter := ParseParameter(expr)
	ap := &_struct.RecursionParam{
		Parameter: parameter,
	}
	switch exprType := expr.(type) {
	case *ast.ArrayType:
		ap.Child = ParseRecursionValue(exprType.Elt)
	case *ast.MapType:
		ap.Child = ParseRecursionValue(exprType.Value)
	}
	return ap
}
