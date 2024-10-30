package param

import (
	"caseGenerator/common/enum"
	"go/ast"
)

// ValueSpec 基本的变量类型，结构简单，只需要 name、type、value
// var a,b,c int
// var x int = 10
type ValueSpec struct {
	BasicParam
	ParamList []ValueSpecParam
}

type ValueSpecParam struct {
	ParamName  string
	ParamValue Parameter
}

func (s *ValueSpec) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *ValueSpec) GetInstance() Parameter {
	return s
}

// GetZeroValue 返回具体类型对应的零值
func (s *ValueSpec) GetZeroValue() Parameter {
	return s
}

func (s *ValueSpec) GetFormula() string {
	return s.GetName()
}

// ParseValueSpec 解析ast
func (s *ValueSpec) ParseValueSpec(expr *ast.ValueSpec, name string) *ValueSpec {
	vs := &ValueSpec{
		BasicParam: BasicParam{
			ParameterType: enum.PARAMETER_TYPE_VALUE_SPEC,
			Name:          name,
		},
	}
	// 已 names 为准
	// 如果 values宽度大于 1，那么 name和 values是一对一
	if len(expr.Values) == 0 {
		parameter := ParseParameter(expr.Values[0])
		paramList := make([]ValueSpecParam, 0, 10)
		for _, v := range expr.Names {
			param := ValueSpecParam{
				ParamName:  v.Name,
				ParamValue: parameter,
			}
			paramList = append(paramList, param)
		}
		vs.ParamList = paramList
		return vs
	} else {
		paramList := make([]ValueSpecParam, 0, 10)
		for i, v := range expr.Names {
			parameter := ParseParameter(expr.Values[i])
			param := ValueSpecParam{
				ParamName:  v.Name,
				ParamValue: parameter,
			}
			paramList = append(paramList, param)
		}
		vs.ParamList = paramList
		return vs
	}
}
