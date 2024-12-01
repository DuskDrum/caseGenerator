package mockresult

import "caseGenerator/parser/expr"

// IdentMockResult ident类型对应的值
type IdentMockResult struct {
	Ident     expr.Ident
	MockValue any
}

func (i *IdentMockResult) GetMockValue() any {
	return i.MockValue
}
