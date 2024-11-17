package expr

import (
	"caseGenerator/common/enum"
	_struct "caseGenerator/parser/struct"
	"go/token"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/samber/lo"
)

// Expression 公式
type Expression struct {
	Expr        string               // "a > 10"
	ElementList []string             // ["a",">","10"]
	IdentMap    map[string]*Ident    // {"a":"a"}
	SelectorMap map[string]*Selector // {"astSelector_a_b_c":"a.b.c"}
	CallMap     map[string]*Call     // {"astCall_funca_c_d":funca(c,d)}
	BasicList   []*BasicLit          // ["10"]
}

func ParseExpression(param _struct.Parameter) []*Expression {
	expressionList := make([]*Expression, 0, 10)
	eList := ParseExpressionParam(param)
	expressionList = append(expressionList, eList...)
	return expressionList
}

// ParseExpressionParam 将Binary、Unary解析为Expression，得到两个东西，一个是里面的Ident和func的引用，一个是最终得到的公式
func ParseExpressionParam(param _struct.Parameter) []*Expression {
	switch exprType := param.(type) {
	case *Binary:
		return MockBinary(exprType)
	case *Unary:
		return MockUnary(exprType)
	case *Parent:
		return MockParent(exprType)
	case *Ident:
		elementList := []string{param.GetFormula()}
		identMap := map[string]*Ident{exprType.IdentName: exprType}

		expression := &Expression{
			ElementList: elementList,
			IdentMap:    identMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *Selector:
		elementList := []string{param.GetFormula()}
		key := strings.ReplaceAll(param.GetFormula(), ".", "_")
		selectorMap := map[string]*Selector{"astSelector_" + key: exprType}

		expression := &Expression{
			ElementList: elementList,
			SelectorMap: selectorMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *Call:
		elementList := []string{param.GetFormula()}
		key := strings.ReplaceAll(param.GetFormula(), ".", "_")
		key = strings.ReplaceAll(key, "(", "_")
		key = strings.ReplaceAll(key, ")", "")

		callMap := map[string]*Call{"astCall_" + key: exprType}

		expression := &Expression{
			ElementList: elementList,
			CallMap:     callMap,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	case *BasicLit:
		elementList := []string{param.GetFormula()}
		basicList := []*BasicLit{exprType}

		expression := &Expression{
			ElementList: elementList,
			BasicList:   basicList,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	default:
		elementList := []string{param.GetFormula()}

		expression := &Expression{
			ElementList: elementList,
			Expr:        strings.Join(elementList, " "),
		}
		return []*Expression{expression}
	}
}

// MockBinary mock Binary
func MockBinary(param *Binary) []*Expression {
	expressionList := make([]*Expression, 0, 10)

	// 如果类型是||逻辑或，剪枝只处理X
	// 解析X
	xExpressionList := ParseExpressionParam(param.X)
	// 解析Y
	yExpressionList := ParseExpressionParam(param.Y)

	// 解析Op
	// 逻辑或进行剪枝处理
	if param.Op == token.LOR {
		// 剪枝处理，只处理||左边的公式
		return xExpressionList
		// 如果类型是&&逻辑与，处理X和Y
	} else if param.Op == token.LAND {
		expressionList = append(expressionList, xExpressionList...)
		expressionList = append(expressionList, yExpressionList...)
		return expressionList
	}
	// 如果类型不是逻辑与或者逻辑或，那么手动组装Expression
	identMap := make(map[string]*Ident, 10)
	callMap := make(map[string]*Call, 10)
	selectorMap := make(map[string]*Selector, 10)
	elementList := make([]string, 0, 10)
	basicList := make([]*BasicLit, 0, 10)

	for _, v := range xExpressionList {
		elementList = append(elementList, v.ElementList...)
		basicList = append(basicList, v.BasicList...)
		for mk, mv := range v.SelectorMap {
			selectorMap[mk] = mv
		}
		for mk, mv := range v.CallMap {
			callMap[mk] = mv
		}
		for mk, mv := range v.IdentMap {
			identMap[mk] = mv
		}
	}

	if param.Op == token.EQL {
		// 等于，底层一定不是binary
		// 不是govaluate.EQ.String()==>"="
		elementList = append(elementList, "==")
	} else if param.Op == token.LSS {
		elementList = append(elementList, govaluate.LT.String())
	} else if param.Op == token.GTR {
		elementList = append(elementList, govaluate.GT.String())
	} else if param.Op == token.LEQ {
		elementList = append(elementList, govaluate.LTE.String())
	} else if param.Op == token.GEQ {
		elementList = append(elementList, govaluate.GTE.String())
	}

	// 暴力破解
	if param.Op == token.ADD {
		elementList = append(elementList, govaluate.PLUS.String())
	} else if param.Op == token.SUB {
		elementList = append(elementList, govaluate.MINUS.String())
	} else if param.Op == token.MUL {
		elementList = append(elementList, govaluate.MULTIPLY.String())
	} else if param.Op == token.QUO {
		elementList = append(elementList, govaluate.DIVIDE.String())
	} else if param.Op == token.REM {
		elementList = append(elementList, govaluate.MODULUS.String())
	}
	for _, v := range yExpressionList {
		elementList = append(elementList, v.ElementList...)
		basicList = append(basicList, v.BasicList...)
		for mk, mv := range v.SelectorMap {
			selectorMap[mk] = mv
		}
		for mk, mv := range v.CallMap {
			callMap[mk] = mv
		}
		for mk, mv := range v.IdentMap {
			identMap[mk] = mv
		}
	}

	return []*Expression{{
		IdentMap:    identMap,
		BasicList:   basicList,
		ElementList: elementList,
		CallMap:     callMap,
		Expr:        strings.Join(elementList, " "),
		SelectorMap: selectorMap,
	}}
}

// MockUnary mock Unary
func MockUnary(param *Unary) []*Expression {
	// 解析公式
	eList := ParseExpressionParam(param.Content)

	elementList := make([]string, 0, 10)
	identMap := make(map[string]*Ident, 10)
	callMap := make(map[string]*Call, 10)
	basicList := make([]*BasicLit, 0, 10)
	selectorMap := make(map[string]*Selector, 10)

	if param.Op == token.NOT {
		elementList = append(elementList, govaluate.INVERT.String())
	} else if param.Op == token.SUB {
		elementList = append(elementList, govaluate.NEGATE.String())
	}

	for _, v := range eList {
		elementList = append(elementList, v.ElementList...)
		basicList = append(basicList, v.BasicList...)
		for mk, mv := range v.SelectorMap {
			selectorMap[mk] = mv
		}
		for mk, mv := range v.CallMap {
			callMap[mk] = mv
		}
		for mk, mv := range v.IdentMap {
			identMap[mk] = mv
		}
	}
	// 解析括号

	return []*Expression{{
		Expr:        strings.Join(elementList, " "),
		ElementList: elementList,
		IdentMap:    identMap,
		CallMap:     callMap,
		BasicList:   basicList,
		SelectorMap: selectorMap,
	}}
}

// MockParent mock Parent
func MockParent(param *Parent) []*Expression {
	// 解析公式
	eList := ParseExpressionParam(param.Content)

	elementList := make([]string, 0, 10)
	identMap := make(map[string]*Ident, 10)
	callMap := make(map[string]*Call, 10)
	basicList := make([]*BasicLit, 0, 10)
	selectorMap := make(map[string]*Selector, 10)

	elementList = append(elementList, "(")
	for _, v := range eList {
		elementList = append(elementList, v.ElementList...)
		basicList = append(basicList, v.BasicList...)
		for mk, mv := range v.SelectorMap {
			selectorMap[mk] = mv
		}
		for mk, mv := range v.CallMap {
			callMap[mk] = mv
		}
		for mk, mv := range v.IdentMap {
			identMap[mk] = mv
		}
	}
	// 解析括号
	elementList = append(elementList, ")")

	return []*Expression{{
		Expr:        strings.Join(elementList, " "),
		ElementList: elementList,
		IdentMap:    identMap,
		CallMap:     callMap,
		BasicList:   basicList,
		SelectorMap: selectorMap,
	}}
}

type MockResult interface {
	GetMockValue() any
}

// IdentMockResult ident类型对应的值
type IdentMockResult struct {
	Ident     Ident
	MockValue any
}

func (i *IdentMockResult) GetMockValue() any {
	return i.MockValue
}

// SelectorMockResult selector类型对应的值
type SelectorMockResult struct {
	Selector  Selector
	MockValue any
}

func (s *SelectorMockResult) GetMockValue() any {
	return s.Selector
}

// CallMockResult call类型对应的值
type CallMockResult struct {
	Call      Call
	MockValue any
}

func (c *CallMockResult) GetMockValue() any {
	return c.MockValue
}

// MockExpression  mock 表达式
//  1. 如果basicLit有值，那么有靶子了，给其他变量赋值代入表达式中
//     如果是 int、float，可能是==、!=、>、<。那么找到 basicLit 的最大或者最小值，计算得到需要
//     如果是 string，那么可能是==或者!=，也有可能是>、<
//     如果是 nil， 可能是==或者!=。 nil是属于 Ident 里的
//
// 2. 如果两边都没有靶子，那么将其中一边设置为零值，再继续用第一步的流程(ident 不知道变量类型，所以没办法处理)
func MockExpression(expression *Expression) MockResult {
	// 1. 如果有 basicLit，那么按照 govalue进行解析试算得到最终结果
	if len(expression.BasicList) > 0 {
		return MockBasicExpression(expression)
	}
	// 2. 如果有nil，那么将其他属性都变成 nil
	if len(expression.IdentMap) > 0 {
		for _, v := range expression.IdentMap {
			if strings.EqualFold(v.IdentName, "nil") {

			}
		}
	}
	// 3. 如果没有 basic 也没有 nil，那么先不处理

	return nil
}

// MockBasicExpression mock 有 basic 的表达式
func MockBasicExpression(expression *Expression) MockResult {
	// todo 这种多个basicLit 的类型一般是一样的，不一样就告警出去
	var specificType *enum.SpecificType
	basicValueList := make([]any, 0, 10)
	for _, v := range expression.BasicList {
		if specificType == nil {
			specificType = &v.SpecificType
		} else if &v.SpecificType != specificType {
			panic("多个 basicLit的类型不一样，请检查")
		}
		basicValueList = append(basicValueList, v.Value)
	}
	// 如果类型是 int、float
	if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_INT {
		return MockBasicIntExpression(expression, basicValueList)
	} else if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_FLOAT64 || lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_FLOAT32 {
		return MockBasicFloatExpression(expression, basicValueList)
	} else if lo.FromPtr(specificType) == enum.SPECIFIC_TYPE_STRING {
		return MockBasicStringExpression(expression, basicValueList)
	}
	return nil
}

// MockBasicIntExpression mock int basic 的表达式
func MockBasicIntExpression(expression *Expression, basicValueList []any) MockResult {
	// 如果类型是 int
	var inList []int
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 int
		if value, ok := v.(int); ok {
			inList = append(inList, value)
		} else {
			// 如果转换失败，返回错误
			panic("element is not an int")
		}
	}
	minValue := lo.Min(inList)
	maxValue := lo.Max(inList)

	params := []string{"a", "b", "c"} // 参数名称

	current := make([]int, len(params))

	result := ComposeInt(params, current, 0, minValue, maxValue, expression.Expr)
	if result != nil {
		// 参数一一对应的值
	}

	return nil
}

// MockBasicFloatExpression mock float basic 的表达式
func MockBasicFloatExpression(expression *Expression, basicValueList []any) MockResult {
	// 如果类型是 float
	var inList []float64
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 float
		if value, ok := v.(float64); ok {
			inList = append(inList, value)
		} else if fValue, fok := v.(float32); fok {
			inList = append(inList, float64(fValue))
		} else {
			// 如果转换失败，返回错误
			panic("element is not an float")
		}
	}
	minValue := lo.Min(inList)
	maxValue := lo.Max(inList)

	params := []string{"a", "b", "c"} // 参数名称

	current := make([]float64, len(params))

	result := ComposeFloat(params, current, 0, minValue, maxValue, expression.Expr)
	if result != nil {
		// 参数一一对应的值
	}

	return nil
}

// MockBasicStringExpression mock string basic 的表达式
func MockBasicStringExpression(expression *Expression, basicValueList []any) MockResult {
	// 如果类型是 string
	var strList []string
	strList = append(strList, "")
	for _, v := range basicValueList {
		// 类型断言，将元素转换为 string
		if value, ok := v.(string); ok {
			strList = append(strList, value)
			strList = append(strList, value+"a")
		} else {
			// 如果转换失败，返回错误
			panic("element is not an string")
		}
	}
	params := []string{"a", "b", "c"} // 参数名称

	current := make([]string, len(params))

	result := ComposeString(params, strList, current, 0, expression.Expr)
	if result != nil {
		// 参数一一对应的值
	}
	return nil
}

// ComposeInt 组合int
func ComposeInt(params []string, current []int, index, min, max int, expr string) []int {
	// 如果当前索引超出了参数范围，保存组合并返回
	if index == len(params) {
		//fmt.Printf("current is:%v \n", current)
		// 匹配参数是否满足不等式
		// 创建 map
		parameters := make(map[string]any)
		// 迭代并填充 map
		for i := 0; i < len(params); i++ {
			parameters[params[i]] = current[i]
		}
		exp, err := govaluate.NewEvaluableExpression(expr)
		if err != nil {
			panic(err.Error())
		}

		result, err := exp.Evaluate(parameters)
		if err != nil {
			panic(err.Error())
		}
		// 尝试类型断言为 bool
		if v, ok := result.(bool); ok {
			if v {
				return current
			}
		}
		return nil
	}

	// 遍历当前参数从 0 到 max 的所有可能值
	for i := min; i <= max; i++ {
		current[index] = i
		composeInt := ComposeInt(params, current, index+1, min, max, expr)
		if composeInt != nil {
			return composeInt
		}
	}
	return nil
}

// ComposeFloat 组合float
func ComposeFloat(params []string, current []float64, index int, min, max float64, expr string) []float64 {
	// 如果当前索引超出了参数范围，保存组合并返回
	if index == len(params) {
		//fmt.Printf("current is:%v \n", current)
		// 匹配参数是否满足不等式
		// 创建 map
		parameters := make(map[string]any)
		// 迭代并填充 map
		for i := 0; i < len(params); i++ {
			parameters[params[i]] = current[i]
		}
		exp, err := govaluate.NewEvaluableExpression(expr)
		if err != nil {
			panic(err.Error())
		}

		result, err := exp.Evaluate(parameters)
		if err != nil {
			panic(err.Error())
		}
		// 尝试类型断言为 bool
		if v, ok := result.(bool); ok {
			if v {
				return current
			}
		}
		return nil
	}

	// 遍历当前参数从 0 到 max 的所有可能值
	for i := min; i <= max; i++ {
		current[index] = i
		composeInt := ComposeFloat(params, current, index+1, min, max, expr)
		if composeInt != nil {
			return composeInt
		}
	}
	return nil
}

// ComposeString 组合string
func ComposeString(params []string, values []string, current []string, index int, expr string) []string {
	// 如果当前索引超出了参数范围，保存组合并返回
	if index == len(params) {
		//fmt.Printf("current is:%v \n", current)
		// 匹配参数是否满足不等式
		// 创建 map
		parameters := make(map[string]any)
		// 迭代并填充 map
		for i := 0; i < len(params); i++ {
			parameters[params[i]] = current[i]
		}
		exp, err := govaluate.NewEvaluableExpression(expr)
		if err != nil {
			panic(err.Error())
		}

		result, err := exp.Evaluate(parameters)
		if err != nil {
			panic(err.Error())
		}
		// 尝试类型断言为 bool
		if v, ok := result.(bool); ok {
			if v {
				return current
			}
		}
		return nil
	}

	// 遍历当前参数的所有可能值
	for _, value := range values {
		current[index] = value
		result := ComposeString(params, values, current, index+1, expr)
		if result != nil {
			return result
		}
	}
	return nil
}
