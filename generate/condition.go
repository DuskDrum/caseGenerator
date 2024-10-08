package generate

import (
	"caseGenerator/parser"
	"go/token"

	"github.com/samber/lo"
)

// Condition 判断条件
type Condition struct {
}

// GenerateCondition 判断某个方法中的条件，找到最简单的路径  ConditionExprType
// 此方法主要用于解决一个条件语句，怎么能走到里面的每一个条件
// 条件需要考虑：1. 用方法的返回值做逻辑运算符的条件 2. 布尔值做逻辑运算符 3. 某个类的属性做逻辑运算符 4. 多个条件做嵌套
// 每个condition一般是一个值做逻辑运算符，比如等于，那么类似于a == b， 或者 a() == b
// 那么处理逻辑如下：
//
//		方法 a() ==> Mock
//	 参数 a.A ==> Assignment: 1. 赋值里面也是方法：Mock 2. 赋值里面是属性：Request/变量/常量  3. 处理不了
func GenerateCondition(si *parser.ConditionNode) {
	// 1. 判断Cond
	if si.Cond == nil {
		panic("condition's Cond can't be nil")
	}
	condInfo := lo.FromPtr(si.Cond)

}

// MockFuncInstruct func mock指令
type MockFuncInstruct struct {
	// MockFunctionName mock方法名
	MockFunctionName string
	// MockReceiverType mock方法对应的Receiver类型。
	// build4 := Mock((*repo.XXXRepo).XXMethod).Return(A, B).Build()
	MockReceiverType string
	// MockReturnList mock对应的响应列表
	MockReturnList []string
}

// MockParamInstruct 参数 mock指令
type MockParamInstruct struct {
	// MockParamName mock参数名
	MockParamName string
	// MockParamValue mock对应的参数值
	MockParamValue string
}

// 解析condition, 得到mock方法指令、mock参数指令
func generateConditionInfo(info *parser.CondInfo, funcInstructList []MockFuncInstruct, paramInstructList []MockParamInstruct) {
	// 说明到底了， paramValue必有值
	if info.XParam.ValueTag == true {
		// x有值了， y还没有值。暂时处理不了
		if info.YParam.ValueTag != true {
			panic("xparam has value but yParam don't has value")
		}
		xp := info.XParam.ParamValue
		yp := info.YParam.ParamValue
		op := info.Op
		// mock值或者mock方法
		MockInstruct(xp, yp, lo.ToPtr(op), funcInstructList, paramInstructList)
	}
	// 如果是parentTag，一般来说只有把几个或当成与的场景会加括号，暂定不考虑括号
	// 逻辑与，需要继续执行。还在处理逻辑，所以需要递归执行
	if info.Op == token.LAND {
		// 逻辑与代表了左右两边都需要执行
		generateConditionInfo(info.XParam)
		generateConditionInfo(info.YParam)

		// 逻辑或，直接跳过处理
	} else if info.Op == token.LOR {
		// 逻辑否，一般逻辑否只会有paramValue，代表的含义是paramValue == false
	} else if info.Op == token.NOT {
		if info.ParamValue == nil {
			panic("token.Not paramValue should not be nil")
		}
		notValue := info.ParamValue

	}
}
