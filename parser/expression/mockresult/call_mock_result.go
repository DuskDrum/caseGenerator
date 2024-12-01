package mockresult

import "caseGenerator/parser/expr"

// CallMockResult call类型对应的值
type CallMockResult struct {
	Call      expr.Call
	MockValue any
}

func (c *CallMockResult) GetMockValue() any {
	return c.MockValue
}
