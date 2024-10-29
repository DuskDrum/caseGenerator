package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

type FuncLit struct {
	BasicParam
}

func (s *FuncLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNCLIT
}

func (s *FuncLit) GetInstance() Parameter {
	return s
}

func (s *FuncLit) GetZeroValue() Parameter {
	panic("implement me")
}

func (s *FuncLit) GetFormula() string {
	panic("implement me")
}

// ParseFuncLit 解析ast
func ParseFuncLit(expr *ast.FuncLit, name string) *FuncLit {
	return nil
}
