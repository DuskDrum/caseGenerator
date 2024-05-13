package vistitor

import (
	"caseGenerator/parse/bo"
	"go/ast"
	"go/token"
	"log"
	"sync"
)

// Condition 应该是树形状
type Condition interface {
	MockTrue()
	MockFalse()
}

// BinaryCondition 使用等式左右判断的
type BinaryCondition struct {
	// 子树列表
	subtreeList []Condition
	// key 是result属性的"A.B.C"， value是param属性的"D.E.F"
	relateMap map[string]string
	// 等式的类型
	Op token.Token
	// 等式的左边
	Ltr any
	// 等式的右边
	Rtr any
	// 位置信息
	PointSite int
}

// MockTrue 要一层一层的打开。返回打开每条判断的mock+打开后对result的影响
func (b *BinaryCondition) MockTrue() {
}

// MockFalse 不需要管子树列表，对这一条进行关闭即可。返回关闭这条判断的mock
func (b *BinaryCondition) MockFalse() {

}

type FuncCondition struct {
	// 子树列表
	subtreeList []Condition
	// key 是result属性的"A.B.C"， value是param属性的"D.E.F"
	relateMap map[string]string
	// 位置信息
	PointSite int
}

func (f *FuncCondition) MockTrue() {
}

// MockFalse 返回关闭这条判断的mock
func (f *FuncCondition) MockFalse() {
}

// ElseCondition 考虑做成属性放在If里面
type ElseCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
}

// ElseIfCondition 条件需要加上前面if 取反
type ElseIfCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
}

// InitCondition 会先初始化， 再判断是否是ok
type InitCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
}

// MixCondition 混合的判断条件
type MixCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
}

type SwitchCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
	// 条件字段
	Factor bo.Param
	// Case值
	CaseValue bo.Param
}

// SwitchFallThroughCondition switch 的fallthrough
type SwitchFallThroughCondition struct {
	// 子树列表
	subtreeList []Condition
	// 位置信息
	PointSite int
	// 条件字段
	Factor bo.Param
	// Case值
	CaseValue bo.Param
}

// CondBody 条件语句下的Body：做了哪些条件相关的处理；返回值了等
type CondBody interface {
}

// ReturnCondBody  条件语句里有return
type ReturnCondBody struct {
	Return Result
}

type ConditionVisitor struct {
	addMu          sync.Mutex
	CondBinaryList []bo.CondBinary
	// 直接调用方法获取 true或者false的逻辑
	CallCondList []bo.CallUnary
}

func (v *ConditionVisitor) AddDetail(cb bo.CondBinary) {
	v.addMu.Lock()
	defer v.addMu.Unlock()
	if v.CondBinaryList == nil {
		v.CondBinaryList = make([]bo.CondBinary, 0, 10)
	}
	v.CondBinaryList = append(v.CondBinaryList, cb)
}

func (v *ConditionVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	if fn, ok := n.(*ast.IfStmt); ok {
		switch expr := fn.Cond.(type) {
		case *ast.BinaryExpr:
			var cb bo.CondBinary
			cb.Op = expr.Op
			cb.PointerSite = expr.OpPos
			x, y := v.BinaryParse(expr)
			if x == nil && y == nil {
				return v
			}
			cb.X = x
			cb.Y = y
			v.AddDetail(cb)
		case *ast.Ident:
			var cb bo.CondBinary
			cb.X = &bo.InvocationUnary{InvocationName: expr.Name}
			cb.Op = token.EQL
			cb.Y = &bo.BasicLitUnary{BasicLitValue: "true"}
			v.AddDetail(cb)
		}
	}
	// SwitchStmt 都是有 左右项的，并且是Equals
	if ss, ok := n.(*ast.SwitchStmt); ok {
		var ssTagUnary bo.Unary
		// 先确定等号的左边
		ssTagUnary = v.BinaryUnaryParse(ss.Tag, ss.Switch)

		for _, ssBody := range ss.Body.List {
			var cb bo.CondBinary
			// default 不处理了
			if ssBody.(*ast.CaseClause).List == nil {
				continue
			}
			expr := ssBody.(*ast.CaseClause)
			y := v.BinaryUnaryParse(expr.List[0], expr.Case)

			cb.Y = y
			cb.Op = token.EQL
			cb.X = ssTagUnary
			cb.PointerSite = expr.Case

			v.AddDetail(cb)
		}
	}
	return v
}

func (v *ConditionVisitor) BinaryParse(expr *ast.BinaryExpr) (X bo.Unary, Y bo.Unary) {
	X = v.BinaryUnaryParse(expr.X, expr.OpPos)
	Y = v.BinaryUnaryParse(expr.Y, expr.OpPos)
	return
}

func (v *ConditionVisitor) BinaryUnaryParse(expr ast.Expr, opPos token.Pos) bo.Unary {
	switch valueSe := expr.(type) {
	case *ast.ParenExpr:
		return v.BinaryUnaryParse(valueSe.X, opPos)
	case *ast.UnaryExpr:
		return v.BinaryUnaryParse(valueSe.X, opPos)
	case *ast.BinaryExpr:
		var cb bo.CondBinary
		X := v.BinaryUnaryParse(valueSe.X, opPos)
		Y := v.BinaryUnaryParse(valueSe.Y, opPos)
		cb.X = X
		cb.Y = Y
		cb.PointerSite = opPos
		cb.Op = valueSe.Op
		v.AddDetail(cb)
	case *ast.BasicLit:
		return &bo.BasicLitUnary{BasicLitValue: valueSe.Value}
	case *ast.Ident:
		return &bo.InvocationUnary{InvocationName: valueSe.Name}
	case *ast.SelectorExpr:
		return &bo.ParamUnary{ParamValue: GetRelationFromSelectorExpr(valueSe)}
	case *ast.CallExpr:
		if se, ok := valueSe.Fun.(*ast.SelectorExpr); ok {
			return &bo.CallUnary{CallValue: GetRelationFromSelectorExpr(se)}
		}
		if id, ok := valueSe.Fun.(*ast.Ident); ok {
			return &bo.CallUnary{CallValue: id.Name}
		}
	default:
		log.Fatalf("未知类型...")
		return nil
	}
	return nil
}

// FilterCondition 过滤直接调用方法获得的true或者false
func (v *ConditionVisitor) FilterCondition() {
	cuMap := make(map[string]*bo.CallUnary, 10)
	cuList := make([]bo.CallUnary, 0, 10)
	cbList := make([]bo.CondBinary, 0, 10)
	// 不增加mock
	for _, cb := range v.CondBinaryList {
		if cb.Op == token.LAND || cb.Op == token.LOR {
			// 方法调用，去重
			if xcu, ok := cb.X.(*bo.CallUnary); ok {
				cuMap[xcu.CallValue] = xcu
			}
			if xcu, ok := cb.Y.(*bo.CallUnary); ok {
				cuMap[xcu.CallValue] = xcu
			}
		} else {
			cbList = append(cbList, cb)
		}
	}
	v.CondBinaryList = cbList
	for _, v := range cuMap {
		cuList = append(cuList, *v)
	}
	v.CallCondList = cuList
}
