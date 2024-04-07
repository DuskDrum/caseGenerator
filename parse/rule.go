package parse

type Rule interface {
	MockTrue() *MockInstruct
	MockFalse() *MockInstruct
}

//// EqlRule 处理Equal类型
//type EqlRule struct {
//	Ltr Unary
//	Rtr Unary
//}
//
//func (e *EqlRule) MockTrue() *MockInstruct {
//	lUnary, lOk := e.Ltr.(*BasicLitUnary)
//	rUnary, rOk := e.Rtr.(*BasicLitUnary)
//	// 此版本只支持 赋值类型的判断
//	if !lOk && !rOk {
//		return nil
//	}
//	// 如果两个都是常量，那么没必要处理
//	if lOk && rOk {
//		return nil
//	}
//	// 左边是常量
//	if lOk {
//		switch eRtr := e.Rtr.(type) {
//		// 常量的调用
//		case *InvocationUnary:
//			// 需要去找常量的映射，找到就发mock指令
//		case
//		}
//	} else {
//		// 右边是常量
//	}
//
//}
//
//func (e *EqlRule) MockFalse() *MockInstruct {
//	unary, lOk := e.Ltr.(*BasicLitUnary)
//	unary, rOk := e.Rtr.(*BasicLitUnary)
//	if !lOk && !rOk {
//		return nil
//	}
//}
//
//// LssRule 处理Less类型 <
//type LssRule struct {
//	Ltr Unary
//	Rtr Unary
//}
//
//// GtrRule 处理greater类型 >
//type GtrRule struct {
//	Ltr Unary
//	Rtr Unary
//}
//
//// NeqRule NotEqual类型
//type NeqRule struct {
//	Ltr Unary
//	Rtr Unary
//}
//
//// LeqRule LessEquals类型
//type LeqRule struct {
//	Ltr Unary
//	Rtr Unary
//}
//
//// GeqRule GreaterEquals类型
//type GeqRule struct {
//	Ltr Unary
//	Rtr Unary
//}
