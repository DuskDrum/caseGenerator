package generate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser"
	"go/token"
	"reflect"

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
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
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
