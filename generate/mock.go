package generate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser"
	"go/token"
	"math"
	"time"

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
	if xParam.AstType == enum.PARAM_AST_TYPE_BasicLit {

	}
}

func MockMinInt(min int) int {
	return intRange(min, math.MaxInt)
}

func MockMaxInt(max int) int {
	return intRange(math.MinInt, max)
}

func MockRangeInt(min, max int) int {
	return intRange(min, max)
}

// intRange intRange
func intRange(min, max int) int {
	// 调整整数范围
	randomNumber := gofakeit.Number(min, max)
	return randomNumber
}

func MockMinUint(min uint) uint {
	return uintRange(min, math.MaxUint)
}

func MockMaxUint(max uint) uint {
	return uintRange(math.MinInt, max)
}

func MockRangeUint(min, max uint) uint {
	return uintRange(min, max)
}

// mockInt uint
func uintRange(min, max uint) uint {
	// 调整整数范围
	randomNumber := gofakeit.UintRange(min, max)
	return randomNumber
}

func MockMinFloat32(min float32) float32 {
	return float32Range(min, math.MaxFloat32)
}

func MockMaxFloat32(max float32) float32 {
	return float32Range(math.MinInt, max)
}

func MockRangeFloat32(min, max float32) float32 {
	return float32Range(min, max)
}

// float32Range uint
func float32Range(min, max float32) float32 {
	// 调整整数范围
	randomNumber := gofakeit.Float32Range(min, max)
	return randomNumber
}

func MockMinFloat64(min float64) float64 {
	return float64Range(min, math.MaxFloat64)
}

func MockMaxFloat64(max float64) float64 {
	return float64Range(math.MinInt, max)
}

func MockRangeFloat64(min, max float64) float64 {
	return float64Range(min, max)
}

// float64Range float64
func float64Range(min, max float64) float64 {
	// 调整整数范围
	randomNumber := gofakeit.Float64Range(min, max)
	return randomNumber
}

// MockMinDate date
func MockMinDate(min time.Time) time.Time {
	// 调整时间范围
	return dateRange(min, time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC))
}

// MockMaxDate date
func MockMaxDate(max time.Time) time.Time {
	// 调整时间范围
	return dateRange(time.Time{}, max)
}

// MockRangeDate date
func MockRangeDate(min, max time.Time) time.Time {
	// 调整时间范围
	return dateRange(min, max)
}

// dateRange date
func dateRange(min, max time.Time) time.Time {
	// 调整时间范围
	return gofakeit.DateRange(min, max)
}

// MockSizeSlice 根据尺寸随机生成切片
func MockSizeSlice(size int64, slicePtr any) {
	faker := gofakeit.New(size)
	faker.Slice(slicePtr)
}

// MockSlice 随机生成切片
func MockSlice(slicePtr any) {
	gofakeit.Slice(slicePtr)
}

func MockSizeString(size uint) string {
	return gofakeit.LetterN(size)
}
func MockString() string {
	return gofakeit.Letter()
}

func MockSizeStruct(size int64, v any) error {
	faker := gofakeit.New(size)
	return faker.Struct(v)
}

func MockStruct(v any) error {
	return gofakeit.Struct(v)
}
