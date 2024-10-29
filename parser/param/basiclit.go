package param

import (
	"caseGenerator/common/enum"
	"go/ast"
	"go/token"

	"github.com/samber/lo"
)

// BasicLit 基本数据类型，token.INT 等
type BasicLit struct {
	BasicParam
	BasicValue
	Kind token.Token
}

func (s *BasicLit) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_BASICLIT
}

func (s *BasicLit) GetInstance() Parameter {
	return s
}

func (s *BasicLit) GetZeroValue() Parameter {
	switch s.Kind {
	case token.INT:
		s.Value = 0
	case token.FLOAT:
		s.Value = lo.Empty[float32]()
	case token.CHAR:
		s.Value = lo.Empty[string]()
	case token.STRING:
		s.Value = lo.Empty[string]()
	default:
		s.Value = nil
	}
	return s
}

func (s *BasicLit) GetFormula() string {
	return s.Name
}

// ParseBasicLit 解析ast
func ParseBasicLit(expr *ast.BasicLit, name string) *Array {
	return nil
}
