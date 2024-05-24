package bo

// MockInstruct Mock指令
type MockInstruct struct {
	MockResponseParam  []ParamParseResult
	MockFunction       string
	MockFunctionParam  []ParamParseResult
	MockFunctionResult []ParamParseResult
}
