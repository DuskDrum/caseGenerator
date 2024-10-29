package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// Parameter 参数接口类型，将参数需要的方法定义出来
type Parameter interface {
	// GetType 获取类型(每个子类维护一个类型)
	GetType() enum.ParameterType
	// GetName 获取name
	GetName() string
	// GetInstance 有些场景这个接口不能满足，直接返回自身去特殊处理
	GetInstance() Parameter
	// GetZeroValue 返回对应的零值
	GetZeroValue() Parameter
	// GetFormula 获取公式，用于展示
	GetFormula() string
}

// ValueAble 可以直接取值的类型
type ValueAble interface {
	GetValue() any
}

type BasicParam struct {
	ParameterType enum.ParameterType
	Name          string
}

// RecursionParam 递归的参数，比如切片
type RecursionParam struct {
	Parameter
	Child *RecursionParam
}

func (b BasicParam) GetName() string {
	return b.Name
}

// BasicValue 具体的值
type BasicValue struct {
	// 这是具体的值
	Value any
}

// SpecificType 参数的具体类型，比如说 int、uint。这个需要组装为另一个枚举
type SpecificType struct {
	enum.SpecificType
}

func (s SpecificType) GetZeroValue() any {
	return s.ZeroValue
}

// ParseParameter 同时处理几种可能存在值的类型，如BasicLit、FuncLit、CompositeLit、CallExpr
// 得放在这个包中，不然会导致循环依赖
func ParseParameter(expr ast.Expr) Parameter {
	switch exprType := expr.(type) {
	case *ast.SelectorExpr:
		return ParseSelector(exprType, "")
	case *ast.Ident:
		return ParseIdent(exprType, "")
		// 指针类型
	case *ast.StarExpr:
		return ParseStar(exprType, "")
	case *ast.FuncType:
		return ParseFuncType(exprType, "")
	case *ast.InterfaceType:
		return ParseInterface(exprType, "")
	case *ast.ArrayType:
		return ParseArray(exprType, "")
	case *ast.MapType:
		return ParseMap(exprType, "")
	case *ast.Ellipsis:
		return ParseEllipsis(exprType, "")
	case *ast.ChanType:
		return ParseChan(exprType, "")
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return ParseIndex(exprType, "")
	case *ast.BasicLit:
		return ParseBasicLit(exprType, "")
		// FuncLit 等待解析出内容值
	case *ast.FuncLit:
		return ParseFuncLit(exprType, "")
	case *ast.CompositeLit:
		return ParseCompositeLit(exprType, "")
		// CallExpr 等待解析出内容值
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		return ParseCall(exprType, "")
	case *ast.KeyValueExpr:
		return ParseKeyValue(exprType, "")
		// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		return ParseUnary(exprType, "")
	case *ast.BinaryExpr:
		return ParseBinary(exprType, "")
	case *ast.ParenExpr:
		return ParseParent(exprType, "")
	default:
		panic("未知类型...")
	}
}
