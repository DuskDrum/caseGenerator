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
		fallthrough
	case token.LSS:
		return ProcessLessValue(value)
	case token.GEQ:
		fallthrough
	case token.GTR:
		return ProcessMoreValue(value)
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
	// todo 综上所述，uintptr先不考虑处理
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
	// complex64是一种复合类型，表示由两个 32 位浮点数组成的复数类型。
	//一个complex64类型的复数由实部和虚部组成，实部和虚部都是float32类型。可以使用内置函数complex来创建一个complex64类型的复数，例如：c := complex(1.2, 3.4)，这里创建了一个实部为 1.2，虚部为 3.4 的复数。
	// c := complex(1.2, 3.4)
	// realPart := real(c)
	// imaginaryPart := imag(c)
	// complex64类型通常在需要进行复数运算的科学计算、信号处理等领域中使用。
	// todo 综上所述，先不考虑Complex64
	case reflect.Complex64:
		panic("still not support complex64")
	// complex128是一种复合类型，表示由两个 64 位浮点数组成的复数类型。
	// 通常在高精度的科学计算、工程计算、数学领域等需要处理复数的场景中使用。
	// todo 综上所述，先不考虑Complex128
	case reflect.Complex128:
		panic("still not support complex128")
	// array是数组类型
	// 只可以使用==或者!=进行比较，不能使用>或者<进行比较，所以不支持
	case reflect.Array:
		panic("still not support array")
	// chan通道类型
	// 只可以使用==或者!=进行比较，不能使用>或者<进行比较，所以不支持
	case reflect.Chan:
		panic("still not support chan")
	// func函数类型
	// 函数类型是一等公民，不能用<、>、==、!=进行比较，所以不支持
	case reflect.Func:
		panic("still not support func")
	// interface接口类型
	// interface{}类型不能直接使用 >、< 进行比较
	// 对于 == 和 !=，当两个 interface{} 类型的值都为具体类型且该具体类型支持比较操作时，可以使用 == 和 != 进行比较；
	// 如果其中一个或两个值是不可比较的类型（如切片、映射、函数等），则进行比较会在编译时出错或在运行时引发 panic。
	case reflect.Interface:
		panic("still not support func")
	// map 是引用类型, 不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，只有当两个 map 的类型完全相同，并且具有相同的键值对时，它们才相等
	// 综上所述，不支持
	case reflect.Map:
		panic("still not support map")
	// pointer是个指针类型，它存储了另一个变量的内存地址
	// 一般情况下，不能直接使用 >、< 比较指针。
	// 对于 == 和 !=，可以用来比较两个指针是否指向同一个变量地址或者都是 nil
	// 综上所述，不支持
	case reflect.Pointer:
		panic("still not support pointer")
	// slice是切片类型，也是引用类型。切片不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，只有当两个切片的类型完全相同，长度相同且所有对应位置的元素都相等时，它们才相等。
	// 综上所述，不支持
	case reflect.Slice:
		panic("still not support slice")
	// string支持>、<、==、!=
	// ==、!=、>、<的比较算法如下，
	// 对于相等性比较（==和!=）：
	//		1. 逐个字节比较两个字符串对应的底层字节。如果所有字节都相等，且两个字符串长度也相同，则认为两个字符串相等。
	//		2. 如果在比较过程中发现有不同的字节，或者两个字符串长度不同，则认为它们不相等。
	// 对于字典序比较（<和>）：
	//		1. 从两个字符串的第一个字节开始比较。
	//		2. 如果对应字节相等，则继续比较下一个字节。
	//		3. 如果一个字符串的当前字节小于另一个字符串的当前字节，则认为该字符串小于另一个字符串；如果一个字符串的当前字节大于另一个字符串的当前字节，则认为该字符串大于另一个字符串。
	//		4. 如果一个字符串在比较过程中先到达末尾（长度较短），则认为它小于另一个字符串（如果另一个字符串还有剩余字节）。
	// 还有一个特性：""空字符串是最小的。
	// "apple" 是小于"applea"的
	// todo: 所以要考虑""<""这种是不满足的
	// 所以直接返回""
	case reflect.String:
		return ""
	// struct 结构体不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，如果结构体的所有字段都是可比较的类型，那么结构体可以进行比较。比较时会逐个比较结构体的对应字段。如果所有字段都相等，则两个结构体相等；
	case reflect.Struct:
		panic("still not support struct")
	// unsafe指针类型，用于进行低级别的内存操作，它可以指向任意类型的变量。
	// 一般情况下，不能直接使用 >、< 对 unsafe.Pointer 进行比较。
	// 对于 == 和 !=，可以用来比较两个 unsafe.Pointer 是否指向同一个内存地址，但这种比较应该非常谨慎地使用，因为不正确的使用可能导致未定义的行为和安全漏洞。
	case reflect.UnsafePointer:
		panic("still not support unsafePointer")
	default:
		return value
	}
}

func ProcessMoreValue(value any) any {
	switch a := reflect.TypeOf(value); a.Kind() {
	case reflect.Bool:
		return value
	case reflect.Int:
		i, ok := value.(int)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(i, math.MaxInt)
	case reflect.Int8:
		i, ok := value.(int8)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(int(i), math.MaxInt)
	case reflect.Int16:
		i, ok := value.(int16)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(int(i), math.MaxInt)
	case reflect.Int32:
		i, ok := value.(int32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(int(i), math.MaxInt)
	case reflect.Int64:
		i, ok := value.(int64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Number(int(i), math.MaxInt)
	case reflect.Uint:
		i, ok := value.(uint)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(i, math.MaxUint)
	case reflect.Uint8:
		i, ok := value.(uint8)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(uint(i), math.MaxUint8)
	case reflect.Uint16:
		i, ok := value.(uint16)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(uint(i), math.MaxUint16)
	case reflect.Uint32:
		i, ok := value.(uint32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(uint(i), math.MaxUint32)
	case reflect.Uint64:
		i, ok := value.(uint64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.UintRange(uint(i), math.MaxUint64)
	//uintptr是一个整数类型，用于存放一个指针值。
	// 实现与指针相关的算术运算：可以对uintptr类型的值进行算术运算，例如加上或减去一个偏移量，以便在内存中进行指针操作。这在一些底层编程场景中很有用，例如与内存映射、不安全的指针操作等结合使用。
	// 辅助垃圾回收器：在某些情况下，uintptr类型的值可以帮助垃圾回收器确定对象是否可达。但需要小心使用，因为不正确的使用可能导致错误的垃圾回收行为。
	// 需要注意的是，使用uintptr类型需要谨慎，因为它涉及到底层的指针操作，可能会导致不安全的行为和错误。在一般的 Go 编程中，应该尽量避免直接使用uintptr类型，除非你非常清楚自己在做什么并且有特定的底层编程需求。
	// uintptr 是可以使用>、<进行比较的
	// 在 Go 语言中，一旦将整数值存储在uintptr类型的变量中，无法直接获取回原始的整数值形式，因为uintptr类型主要用于指针相关的操作，并不直接对应一个整数值的表示形式。
	//如果想要获取原始的整数值，可以考虑在将整数值转换为uintptr之前，将整数值存储在另一个变量中，或者在转换时记录下原始整数值以便后续使用。
	// todo 综上所述，uintptr先不考虑处理
	case reflect.Uintptr:
		panic("still not support uintptr type")
	case reflect.Float32:
		i, ok := value.(float32)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Float32Range(i, math.MaxFloat32)
	case reflect.Float64:
		i, ok := value.(float64)
		if !ok {
			panic("类型异常")
		}
		return gofakeit.Float64Range(i, math.MaxFloat64)
	// complex64是一种复合类型，表示由两个 32 位浮点数组成的复数类型。
	//一个complex64类型的复数由实部和虚部组成，实部和虚部都是float32类型。可以使用内置函数complex来创建一个complex64类型的复数，例如：c := complex(1.2, 3.4)，这里创建了一个实部为 1.2，虚部为 3.4 的复数。
	// c := complex(1.2, 3.4)
	// realPart := real(c)
	// imaginaryPart := imag(c)
	// complex64类型通常在需要进行复数运算的科学计算、信号处理等领域中使用。
	// todo 综上所述，先不考虑Complex64
	case reflect.Complex64:
		panic("still not support complex64")
	// complex128是一种复合类型，表示由两个 64 位浮点数组成的复数类型。
	// 通常在高精度的科学计算、工程计算、数学领域等需要处理复数的场景中使用。
	// todo 综上所述，先不考虑Complex128
	case reflect.Complex128:
		panic("still not support complex128")
	// array是数组类型
	// 只可以使用==或者!=进行比较，不能使用>或者<进行比较，所以不支持
	case reflect.Array:
		panic("still not support array")
	// chan通道类型
	// 只可以使用==或者!=进行比较，不能使用>或者<进行比较，所以不支持
	case reflect.Chan:
		panic("still not support chan")
	// func函数类型
	// 函数类型是一等公民，不能用<、>、==、!=进行比较，所以不支持
	case reflect.Func:
		panic("still not support func")
	// interface接口类型
	// interface{}类型不能直接使用 >、< 进行比较
	// 对于 == 和 !=，当两个 interface{} 类型的值都为具体类型且该具体类型支持比较操作时，可以使用 == 和 != 进行比较；
	// 如果其中一个或两个值是不可比较的类型（如切片、映射、函数等），则进行比较会在编译时出错或在运行时引发 panic。
	case reflect.Interface:
		panic("still not support func")
	// map 是引用类型, 不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，只有当两个 map 的类型完全相同，并且具有相同的键值对时，它们才相等
	// 综上所述，不支持
	case reflect.Map:
		panic("still not support map")
	// pointer是个指针类型，它存储了另一个变量的内存地址
	// 一般情况下，不能直接使用 >、< 比较指针。
	// 对于 == 和 !=，可以用来比较两个指针是否指向同一个变量地址或者都是 nil
	// 综上所述，不支持
	case reflect.Pointer:
		panic("still not support pointer")
	// slice是切片类型，也是引用类型。切片不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，只有当两个切片的类型完全相同，长度相同且所有对应位置的元素都相等时，它们才相等。
	// 综上所述，不支持
	case reflect.Slice:
		panic("still not support slice")
	// string支持>、<、==、!=
	// ==、!=、>、<的比较算法如下，
	// 对于相等性比较（==和!=）：
	//		1. 逐个字节比较两个字符串对应的底层字节。如果所有字节都相等，且两个字符串长度也相同，则认为两个字符串相等。
	//		2. 如果在比较过程中发现有不同的字节，或者两个字符串长度不同，则认为它们不相等。
	// 对于字典序比较（<和>）：
	//		1. 从两个字符串的第一个字节开始比较。
	//		2. 如果对应字节相等，则继续比较下一个字节。
	//		3. 如果一个字符串的当前字节小于另一个字符串的当前字节，则认为该字符串小于另一个字符串；如果一个字符串的当前字节大于另一个字符串的当前字节，则认为该字符串大于另一个字符串。
	//		4. 如果一个字符串在比较过程中先到达末尾（长度较短），则认为它小于另一个字符串（如果另一个字符串还有剩余字节）。
	// 还有一个特性：""空字符串是最小的。
	// "apple" 是小于"applea"的
	// todo: 所以要考虑""<""这种是不满足的
	// 所以直接返回""
	case reflect.String:
		i, ok := value.(string)
		if !ok {
			panic("类型异常")
		}
		return i + "a"
	// struct 结构体不能直接使用 >、< 进行比较。
	// 对于 == 和 !=，如果结构体的所有字段都是可比较的类型，那么结构体可以进行比较。比较时会逐个比较结构体的对应字段。如果所有字段都相等，则两个结构体相等；
	case reflect.Struct:
		panic("still not support struct")
	// unsafe指针类型，用于进行低级别的内存操作，它可以指向任意类型的变量。
	// 一般情况下，不能直接使用 >、< 对 unsafe.Pointer 进行比较。
	// 对于 == 和 !=，可以用来比较两个 unsafe.Pointer 是否指向同一个内存地址，但这种比较应该非常谨慎地使用，因为不正确的使用可能导致未定义的行为和安全漏洞。
	case reflect.UnsafePointer:
		panic("still not support unsafePointer")
	default:
		return value
	}
}
