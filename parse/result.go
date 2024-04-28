package parse

import (
	"go/ast"
	"go/token"
	"log"
	"sync"
)

type Result interface {
	GetResultInfo() any
}

// NilResult nil响应
type NilResult struct {
}

func (n *NilResult) GetResultInfo() any {
	return nil
}

// ErrorResult error响应
type ErrorResult struct {
	Error error
}

func (e *ErrorResult) GetResultInfo() any {
	return e.Error
}

// VoidResult 没有响应
type VoidResult struct {
}

func (v *VoidResult) GetResultInfo() any {
	return nil
}

// ParamResult 参数型响应，需要将每处改变这里的逻辑都记下
type ParamResult struct {
	// key是a.b.c， value是d.e.f
	// 位置信息
	PointSite token.Pos
	// string 响应信息 (a.b.c)
	ResultMark string
	// string 响应信息的赋值 (d.e.f)
	//ResultValueMark string
}

//type ParamResultDetail struct {
//	// 位置信息
//	PointSite int
//	// string 响应信息 (a.b.c)
//	ResultMark string
//	// string 响应信息的赋值 (d.e.f)
//	//ResultValueMark string
//}

func (p *ParamResult) GetResultInfo() any {

	attributeMap := make(map[string]string, 10)

	return attributeMap
}

// FuncResult 方法型响应 Package.Func
type FuncResult struct {
	// 调用方法的信息。import执行时再匹配
	//FuncImport      string
	FuncPackageName string
	FuncName        string
}

// GetResultInfo 获取响应信息，可以直接使用mock决定返回的是什么
func (f *FuncResult) GetResultInfo() any {
	return f
}

// CompositeResult 构造函数响应
type CompositeResult struct {
	// 默认是此包中的
	ResultPackageName string
	ResultStructName  string
	Relations         []AssignmentBinary
}

func (c CompositeResult) GetUnaryValue() any {
	return c
}

func (c CompositeResult) GetResultInfo() any {
	return c
}

type ResultVisitor struct {
	mu                      sync.Mutex
	addMu                   sync.Mutex
	ResultVisitorDetailList []ResultVisitorDetail
}

func (v *ResultVisitor) AddDetail(detail ResultVisitorDetail) {
	v.addMu.Lock()
	defer v.addMu.Unlock()
	if v.ResultVisitorDetailList == nil {
		v.ResultVisitorDetailList = make([]ResultVisitorDetail, 0, 10)
	}
	v.ResultVisitorDetailList = append(v.ResultVisitorDetailList, detail)
}

type ResultVisitorDetail struct {
	// 返回值的位置
	PointSite token.Pos
	// 返回值的索引，从0开始
	Index int
	// 解析出的返回值
	Result Result
}

func (v *ResultVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	// 暂时只支持传入funcType
	//funcType, ok := n.(*ast.FuncType)
	//if !ok {
	//	return v
	//}
	//list := funcType.Results.List
	//for i, result := range list {
	//	names := result.Names
	//	if names != nil && result.Type != nil {
	//		ParamParse(result.Type, uuid.NewUUID(), )
	//
	//	}
	//}

	//if fn, ok := n.(*ast.ReturnStmt); ok {
	//	v.mu.Lock()
	//	for i, res := range fn.Results {
	//		var detail ResultVisitorDetail
	//		detail.PointSite = fn.Return
	//		detail.Index = i
	//		// 返回某个参数
	//		if ident, ok := res.(*ast.Ident); ok {
	//			// 返回nil就不需要处理了
	//			if ident.Obj == nil {
	//				continue
	//			}
	//			pr := ParamResult{
	//				PointSite: ident.NamePos,
	//				// 某个变量的名字
	//				ResultMark: ident.Name,
	//			}
	//			detail.Result = &pr
	//		}
	//		// 返回某个方法：Package.Func
	//		if callExpr, ok := res.(*ast.CallExpr); ok {
	//			if se, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
	//				fr := FuncResult{
	//					FuncPackageName: GetRelationFromSelectorExpr(se),
	//					FuncName:        se.Sel.Name,
	//				}
	//				detail.Result = &fr
	//			}
	//		}
	//		// 直接返回类的构造函数，比如说返回 XXX{A:"1",B:"2"}， 那么可以拿到关联，比如XXX.A="1"；XXX.B=C.D
	//		if unaryExpr, ok := res.(*ast.UnaryExpr); ok {
	//			if se, ok := unaryExpr.X.(*ast.CompositeLit); ok {
	//				cr := CompositeLitParse(se)
	//				detail.Result = &cr
	//			}
	//			if ident, ok := unaryExpr.X.(*ast.Ident); ok {
	//				pr := ParamResult{
	//					PointSite: ident.NamePos,
	//					// 某个变量的名字
	//					ResultMark: ident.Name,
	//				}
	//				detail.Result = &pr
	//			}
	//		}
	//		v.AddDetail(detail)
	//	}
	//	v.mu.Unlock()
	//}
	return v
}

func GetRelationFromSelectorExpr(se *ast.SelectorExpr) string {
	if si, ok := se.X.(*ast.Ident); ok {
		return si.Name + "." + se.Sel.Name
	}
	if sse, ok := se.X.(*ast.SelectorExpr); ok {
		return GetRelationFromSelectorExpr(sse) + "." + se.Sel.Name
	}
	return se.Sel.Name
}

func CompositeLitParse(se *ast.CompositeLit) CompositeResult {
	cr := CompositeResult{}
	cr.ResultPackageName = se.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
	cr.ResultStructName = se.Type.(*ast.SelectorExpr).Sel.Name
	rprs := make([]AssignmentBinary, 0, 10)
	for _, el := range se.Elts {
		var rpr AssignmentBinary
		if kve, ok := el.(*ast.KeyValueExpr); ok {
			rpr.X = ParamUnary{ParamValue: kve.Key.(*ast.Ident).Name}

			switch valueSe := kve.Value.(type) {
			case *ast.BasicLit:
				rpr.Y = &BasicLitUnary{BasicLitValue: valueSe.Value}
			case *ast.Ident:
				rpr.Y = &InvocationUnary{InvocationName: valueSe.Name}
			case *ast.SelectorExpr:
				rpr.Y = &ParamUnary{ParamValue: GetRelationFromSelectorExpr(valueSe)}
			case *ast.CallExpr:
				if se, ok := valueSe.Fun.(*ast.SelectorExpr); ok {
					rpr.Y = &CallUnary{CallValue: GetRelationFromSelectorExpr(se)}
				}
				if id, ok := valueSe.Fun.(*ast.Ident); ok {
					rpr.Y = &CallUnary{CallValue: id.Name}
				}
			default:
				log.Fatalf("未知类型...")
			}
		} else {
			continue
		}
		rprs = append(rprs, rpr)
	}
	cr.Relations = rprs

	return cr
}
