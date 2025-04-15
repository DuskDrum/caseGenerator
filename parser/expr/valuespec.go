package expr

import (
	"caseGenerator/common/enum"
	"caseGenerator/parser/struct"
	"go/ast"
)

// ValueSpec 基本的变量类型，结构简单，只需要 name、type、value
// var a,b,c int
// var x int = 10
type ValueSpec struct {
	ParamList []ValueSpecParam
}

type ValueSpecParam struct {
	ParamName string
	ParamType _struct.Parameter
}

func (s *ValueSpec) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_IDENT
}

func (s *ValueSpec) GetFormula() string {
	formula := ""
	for _, p := range s.ParamList {
		formula += "var " + p.ParamName + " " + p.ParamType.GetFormula() + "\n"
	}
	return formula
}

// ParseValueSpec 解析ast
func ParseValueSpec(expr *ast.ValueSpec, af *ast.File) *ValueSpec {
	vs := &ValueSpec{}
	// 已 names 为准
	// 如果 values宽度大于 1，那么 name和 values是一对一
	if len(expr.Values) == 0 {
		parameter := ParseParameter(expr.Values[0], af)
		paramList := make([]ValueSpecParam, 0, 10)
		for _, v := range expr.Names {
			param := ValueSpecParam{
				ParamName: v.Name,
				ParamType: parameter,
			}
			paramList = append(paramList, param)
		}
		vs.ParamList = paramList
		return vs
	} else {
		paramList := make([]ValueSpecParam, 0, 10)
		for i, v := range expr.Names {
			parameter := ParseParameter(expr.Values[i], af)
			param := ValueSpecParam{
				ParamName: v.Name,
				ParamType: parameter,
			}
			paramList = append(paramList, param)
		}
		vs.ParamList = paramList
		return vs
	}
}
