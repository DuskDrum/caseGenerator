package enum

import (
	"encoding/json"

	"github.com/samber/lo"
)

// SpecificType 具体类型，对应的是 reflect.Type
type SpecificType struct {
	Name      string
	Desc      string
	ZeroValue any
}

var (
	SPECIFIC_TYPE_BOOL          = SpecificType{Name: "Bool", ZeroValue: lo.Empty[bool](), Desc: "布尔类型"}
	SPECIFIC_TYPE_INT           = SpecificType{Name: "Int", ZeroValue: lo.Empty[int](), Desc: "Int"}
	SPECIFIC_TYPE_INT8          = SpecificType{Name: "Int8", ZeroValue: lo.Empty[int8](), Desc: "Int8"}
	SPECIFIC_TYPE_INT16         = SpecificType{Name: "Int16", ZeroValue: lo.Empty[int16](), Desc: "Int16"}
	SPECIFIC_TYPE_INT32         = SpecificType{Name: "Int32", ZeroValue: lo.Empty[int32](), Desc: "Int32"}
	SPECIFIC_TYPE_INT64         = SpecificType{Name: "Int64", ZeroValue: lo.Empty[int64](), Desc: "Int64"}
	SPECIFIC_TYPE_UINT          = SpecificType{Name: "Uint", ZeroValue: lo.Empty[uint](), Desc: "Uint"}
	SPECIFIC_TYPE_UINT8         = SpecificType{Name: "Uint8", ZeroValue: lo.Empty[uint8](), Desc: "Uint8"}
	SPECIFIC_TYPE_UINT16        = SpecificType{Name: "Uint16", ZeroValue: lo.Empty[uint16](), Desc: "Uint16"}
	SPECIFIC_TYPE_UINT32        = SpecificType{Name: "Uint32", ZeroValue: lo.Empty[uint32](), Desc: "Uint32"}
	SPECIFIC_TYPE_UINT64        = SpecificType{Name: "Uint64", ZeroValue: lo.Empty[uint64](), Desc: "Uint64"}
	SPECIFIC_TYPE_UINTPTR       = SpecificType{Name: "Uintptr", ZeroValue: nil, Desc: "uintptr是一个整数类型，用于存放一个指针值。实现与指针相关的算术运算：可以对uintptr类型的值进行算术运算，例如加上或减去一个偏移量，以便在内存中进行指针操作。这在一些底层编程场景中很有用，例如与内存映射、不安全的指针操作等结合使用。辅助垃圾回收器：在某些情况下，uintptr类型的值可以帮助垃圾回收器确定对象是否可达。但需要小心使用，因为不正确的使用可能导致错误的垃圾回收行为。需要注意的是，使用uintptr类型需要谨慎，因为它涉及到底层的指针操作，可能会导致不安全的行为和错误。在一般的 Go 编程中，应该尽量避免直接使用uintptr类型，除非你非常清楚自己在做什么并且有特定的底层编程需求。uintptr 是可以使用>、<进行比较的。在 Go 语言中，一旦将整数值存储在uintptr类型的变量中，无法直接获取回原始的整数值形式，因为uintptr类型主要用于指针相关的操作，并不直接对应一个整数值的表示形式。如果想要获取原始的整数值，可以考虑在将整数值转换为uintptr之前，将整数值存储在另一个变量中，或者在转换时记录下原始整数值以便后续使用。"}
	SPECIFIC_TYPE_FLOAT32       = SpecificType{Name: "Float32", ZeroValue: lo.Empty[float32](), Desc: "Float32"}
	SPECIFIC_TYPE_FLOAT64       = SpecificType{Name: "Float64", ZeroValue: lo.Empty[float64](), Desc: "Float64"}
	SPECIFIC_TYPE_COMPLEX64     = SpecificType{Name: "Complex64", ZeroValue: lo.Empty[complex64](), Desc: "complex64是一种复合类型，表示由两个 32 位浮点数组成的复数类型。不支持使用>, <, ==, !=等比较运算符进行直接比较。不过他的实部和虚部可以比较 real(c1) == real(c2) && imag(c1) == imag(c2)。complex64类型通常在需要进行复数运算的科学计算、信号处理等领域中使用。"}
	SPECIFIC_TYPE_COMPLEX128    = SpecificType{Name: "Complex128", ZeroValue: lo.Empty[complex128](), Desc: "complex128是一种复合类型，表示由两个 64 位浮点数组成的复数类型。complex128类型也不支持使用 ==, !=等比较运算符进行直接比较。通常在高精度的科学计算、工程计算、数学领域等需要处理复数的场景中使用。"}
	SPECIFIC_TYPE_ARRAY         = SpecificType{Name: "Array", ZeroValue: nil, Desc: "Array"}
	SPECIFIC_TYPE_CHAN          = SpecificType{Name: "Chan", ZeroValue: nil, Desc: "Chan"}
	SPECIFIC_TYPE_FUNC          = SpecificType{Name: "Func", ZeroValue: nil, Desc: "函数类型是一等公民，不能用<、>、==、!=进行比较，所以不支持"}
	SPECIFIC_TYPE_INTERFACE     = SpecificType{Name: "Interface", ZeroValue: nil, Desc: "对于 == 和 !=，当两个 interface{} 类型的值都为具体类型且该具体类型支持比较操作时，可以使用 == 和 != 进行比较；如果其中一个或两个值是不可比较的类型（如切片、映射、函数等），则进行比较会在编译时出错或在运行时引发 panic。"}
	SPECIFIC_TYPE_MAP           = SpecificType{Name: "Map", ZeroValue: nil, Desc: "map 是引用类型, == 和 !=，只有当两个 map 的类型完全相同，并且具有相同的键值对时，它们才相等"}
	SPECIFIC_TYPE_POINTER       = SpecificType{Name: "Pointer", ZeroValue: nil, Desc: "pointer是个指针类型，它存储了另一个变量的内存地址。对于 == 和 !=，可以用来比较两个指针是否指向同一个变量地址或者都是 nil"}
	SPECIFIC_TYPE_SLICE         = SpecificType{Name: "Slice", ZeroValue: nil, Desc: "slice是切片类型，也是引用类型。对于 == 和 !=，只有当两个切片的类型完全相同，长度相同且所有对应位置的元素都相等时，它们才相等。"}
	SPECIFIC_TYPE_STRING        = SpecificType{Name: "String", ZeroValue: lo.Empty[string](), Desc: "string支持>、<、==、!=。==、!=、>、<的比较算法如下，对于相等性比较（==和!=）：1. 逐个字节比较两个字符串对应的底层字节。如果所有字节都相等，且两个字符串长度也相同，则认为两个字符串相等。2. 如果在比较过程中发现有不同的字节，或者两个字符串长度不同，则认为它们不相等。对于字典序比较（<和>）：1. 从两个字符串的第一个字节开始比较。2. 如果对应字节相等，则继续比较下一个字节。3. 如果一个字符串的当前字节小于另一个字符串的当前字节，则认为该字符串小于另一个字符串；如果一个字符串的当前字节大于另一个字符串的当前字节，则认为该字符串大于另一个字符串。4. 如果一个字符串在比较过程中先到达末尾（长度较短），则认为它小于另一个字符串（如果另一个字符串还有剩余字节）。还有一个特性：\"\"空字符串是最小的。\"apple\" 是小于\"applea\"的。"}
	SPECIFIC_TYPE_STRUCT        = SpecificType{Name: "Struct", ZeroValue: nil, Desc: "struct 结构体不能直接使用 >、< 进行比较。对于 == 和 !=，如果结构体的所有字段都是可比较的类型，那么结构体可以进行比较。比较时会逐个比较结构体的对应字段。如果所有字段都相等，则两个结构体相等；"}
	SPECIFIC_TYPE_UNSAFEPOINTER = SpecificType{Name: "UnsafePointer", ZeroValue: nil, Desc: "unsafe指针类型，用于进行低级别的内存操作，它可以指向任意类型的变量。一般情况下，不能直接使用 >、< 对 unsafe.Pointer 进行比较。对于 == 和 !=，可以用来比较两个 unsafe.Pointer 是否指向同一个内存地址，但这种比较应该非常谨慎地使用，因为不正确的使用可能导致未定义的行为和安全漏洞。"}
)

var ALL_SPECIFIC_TYPE = map[string]SpecificType{
	SPECIFIC_TYPE_BOOL.Name:          SPECIFIC_TYPE_BOOL,
	SPECIFIC_TYPE_INT.Name:           SPECIFIC_TYPE_INT,
	SPECIFIC_TYPE_INT8.Name:          SPECIFIC_TYPE_INT8,
	SPECIFIC_TYPE_INT16.Name:         SPECIFIC_TYPE_INT16,
	SPECIFIC_TYPE_INT32.Name:         SPECIFIC_TYPE_INT32,
	SPECIFIC_TYPE_INT64.Name:         SPECIFIC_TYPE_INT64,
	SPECIFIC_TYPE_UINT.Name:          SPECIFIC_TYPE_UINT,
	SPECIFIC_TYPE_UINT8.Name:         SPECIFIC_TYPE_UINT8,
	SPECIFIC_TYPE_UINT16.Name:        SPECIFIC_TYPE_UINT16,
	SPECIFIC_TYPE_UINT32.Name:        SPECIFIC_TYPE_UINT32,
	SPECIFIC_TYPE_UINT64.Name:        SPECIFIC_TYPE_UINT64,
	SPECIFIC_TYPE_UINTPTR.Name:       SPECIFIC_TYPE_UINTPTR,
	SPECIFIC_TYPE_FLOAT32.Name:       SPECIFIC_TYPE_FLOAT32,
	SPECIFIC_TYPE_FLOAT64.Name:       SPECIFIC_TYPE_FLOAT64,
	SPECIFIC_TYPE_COMPLEX64.Name:     SPECIFIC_TYPE_COMPLEX64,
	SPECIFIC_TYPE_COMPLEX128.Name:    SPECIFIC_TYPE_COMPLEX128,
	SPECIFIC_TYPE_ARRAY.Name:         SPECIFIC_TYPE_ARRAY,
	SPECIFIC_TYPE_CHAN.Name:          SPECIFIC_TYPE_CHAN,
	SPECIFIC_TYPE_FUNC.Name:          SPECIFIC_TYPE_FUNC,
	SPECIFIC_TYPE_INTERFACE.Name:     SPECIFIC_TYPE_INTERFACE,
	SPECIFIC_TYPE_MAP.Name:           SPECIFIC_TYPE_MAP,
	SPECIFIC_TYPE_POINTER.Name:       SPECIFIC_TYPE_POINTER,
	SPECIFIC_TYPE_SLICE.Name:         SPECIFIC_TYPE_SLICE,
	SPECIFIC_TYPE_STRING.Name:        SPECIFIC_TYPE_STRING,
	SPECIFIC_TYPE_STRUCT.Name:        SPECIFIC_TYPE_STRUCT,
	SPECIFIC_TYPE_UNSAFEPOINTER.Name: SPECIFIC_TYPE_UNSAFEPOINTER,
}

func (at *SpecificType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Name)
}

func (at *SpecificType) UnmarshalJSON(data []byte) error {
	// 自定义序列化逻辑
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	assignmentType, ok := ALL_SPECIFIC_TYPE[name]
	if ok {
		at.Name = assignmentType.Name
		at.Desc = assignmentType.Desc
	}
	return nil
}
