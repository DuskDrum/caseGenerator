package mockresult

import "caseGenerator/parser/expr"

// SelectorMockResult selector类型对应的值
type SelectorMockResult struct {
	Selector  expr.Selector
	MockValue any
}

func (s *SelectorMockResult) GetMockValue() any {
	return s.Selector
}
