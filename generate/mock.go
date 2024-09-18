package generate

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser"
	"go/token"

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

// MockInstruct 根据符号左边和右边，组装mock指令
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
