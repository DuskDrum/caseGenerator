package generate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser"
	"go/token"
	"math"
	"reflect"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/samber/lo"
)

// Mock 记录mock的信息
type Mock struct {
	// request相关
	Request Request
	// response 返回值相关
	Response Response
	// go:linkname关联的方法
	// go:linkname awxCommonConvertSettlementReportAlgorithm slp/reconcile/core/message/standard.commonConvertSettlementReportAlgorithm
	// func awxCommonConvertSettlementReportAlgorithm(transactionType enums.TransactionType, createdAt time.Time, ctx context.Context, dataBody dto.AwxSettlementReportDataBody) (result []service.OrderAlgorithmResult, err error)
	GoLinkedName string
}

type Response struct {
	// assert 返回值
	AssertList []Assert
	// 参数
	ParamList []parser.ParamValue
}

type Assert struct {
	// assert的参数
	Param parser.ParamValue
	// assert的值
	Value string
}

// 找到方法receiver、值的列表、request列表、go:linkedname关联的方法
// 	EQL    // ==
//	NEQ      // !=
//	LEQ      // <=
//	GEQ      // >=
//	LSS    // <
//	GTR    // >

// MockInstruct 根据符号左边和右边，组装mock指令。
// 约束求解器， SMT（Satisfiability Modulo Theories，可满足性模理论）
// 2-CNF可满足性问题 和 3-CNF可满足性问题
// NP完全性
// 电路可满足性问题
func MockInstruct(xParam *parser.ParamValue, yParam *parser.ParamValue, op *token.Token, funcInstructList []MockFuncInstruct, paramInstructList []MockParamInstruct) {
	// 校验参数
	if xParam == nil || yParam == nil || op == nil {
		return
	}
	opPt := lo.FromPtr(op)
	if opPt != token.EQL && opPt != token.NEQ && opPt != token.LEQ && opPt != token.GEQ && opPt != token.LSS && opPt != token.GTR {
		panic("MockInstruct don't support this Op: " + op.String())
	}
	// 1. 第一场景， xParam是定量BasicLit， yParam是变量Ident(x、y交换同理)
	if xParam.AstType == enum.PARAM_AST_TYPE_BasicLit && yParam.AstType == enum.PARAM_AST_TYPE_Ident {

	}

	if yParam.AstType == enum.PARAM_AST_TYPE_BasicLit && xParam.AstType == enum.PARAM_AST_TYPE_Ident {

	}
}

// MockMatchConditionValue 根据指定的value 和 比较运算符token 得到符合条件对应的值
func MockMatchConditionValue(op *token.Token, value any) any {
	switch *op {
	case token.EQL:
		return value
	case token.NEQ:
		fallthrough
	case token.LEQ:
	case token.GEQ:
	case token.LSS:
	case token.GTR:

	}
	return nil
}

func ProcessLessValue(value any) any {
	switch a := reflect.TypeOf(value); a.Kind() {
	case reflect.Bool:
		return value
	case reflect.Int:
		i, ok := value.(int)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(math.MinInt, i)
	case reflect.Int8:
		i, ok := value.(int8)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(math.MinInt, int(i))
	case reflect.Int16:
		i, ok := value.(int16)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(math.MinInt, int(i))
	case reflect.Int32:
		i, ok := value.(int32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(math.MinInt, int(i))
	case reflect.Int64:
		i, ok := value.(int64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(math.MinInt, int(i))
	case reflect.Uint:
		i, ok := value.(uint)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(0, i)
	case reflect.Uint8:
		i, ok := value.(uint8)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(0, uint(i))
	case reflect.Uint16:
		i, ok := value.(uint16)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(0, uint(i))
	case reflect.Uint32:
		i, ok := value.(uint32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(0, uint(i))
	case reflect.Uint64:
		i, ok := value.(uint64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(0, uint(i))
	//uintptr是一个整数类型，用于存放一个指针值。
	// 实现与指针相关的算术运算：可以对uintptr类型的值进行算术运算，例如加上或减去一个偏移量，以便在内存中进行指针操作。这在一些底层编程场景中很有用，例如与内存映射、不安全的指针操作等结合使用。
	// 辅助垃圾回收器：在某些情况下，uintptr类型的值可以帮助垃圾回收器确定对象是否可达。但需要小心使用，因为不正确的使用可能导致错误的垃圾回收行为。
	// 需要注意的是，使用uintptr类型需要谨慎，因为它涉及到底层的指针操作，可能会导致不安全的行为和错误。在一般的 Go 编程中，应该尽量避免直接使用uintptr类型，除非你非常清楚自己在做什么并且有特定的底层编程需求。
	// uintptr 是可以使用>、<进行比较的
	// 在 Go 语言中，一旦将整数值存储在uintptr类型的变量中，无法直接获取回原始的整数值形式，因为uintptr类型主要用于指针相关的操作，并不直接对应一个整数值的表示形式。
	//如果想要获取原始的整数值，可以考虑在将整数值转换为uintptr之前，将整数值存储在另一个变量中，或者在转换时记录下原始整数值以便后续使用。
	// 综上所述，uintptr先不考虑处理
	case reflect.Uintptr:
		panic("still not support uintptr type")
	case reflect.Float32:
		i, ok := value.(float32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Float32Range(0, i)
	case reflect.Float64:
		i, ok := value.(float64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Float64Range(0, i)
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Array:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Pointer:
	case reflect.Slice:
	case reflect.String:
	case reflect.Struct:
	case reflect.UnsafePointer:
	default:
		return value
	}

}
