package expr

import (
	"caseGenerator/common/enum"
	"go/ast"
	"strings"

	"github.com/samber/lo"
)

// FuncType 定义方法特殊结构，方法名、请求参数列表、响应参数列表
// *ast.FuncType 用于表示函数类型。它包含函数的参数和结果的类型信息。
type FuncType struct {
	TypeParams []Field
	Params     []Field
	Results    []Field
}

func (s *FuncType) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_FUNC_TYPE
}

func (s *FuncType) GetFormula() string {
	formula := ""
	if len(s.TypeParams) != 0 {
		// todo 泛型需要确认下逻辑
	}

	if len(s.Params) != 0 {
		formulaList := make([]string, 0, 10)
		for _, v := range s.Params {
			paramFormula := v.GetFormula()
			formulaList = append(formulaList, paramFormula)
		}
		paramStr := strings.Join(formulaList, ", ")
		formula = "(" + paramStr + ")"
	} else {
		formula += "()"
	}

	if len(s.Results) != 0 {
		resultList := make([]string, 0, 10)
		for _, v := range s.Results {
			paramFormula := v.GetFormula()
			resultList = append(resultList, paramFormula)
		}
		resultStr := strings.Join(resultList, ", ")
		formula = "(" + resultStr + ")"
	}

	return formula
}

// ParseFuncType 解析ast
func ParseFuncType(expr *ast.FuncType, af *ast.File) *FuncType {
	ft := &FuncType{}

	if expr.TypeParams != nil {
		typeParamList := make([]Field, 0, 10)
		for _, v := range expr.TypeParams.List {
			field := ParseField(v, af)
			if field != nil {
				typeParamList = append(typeParamList, lo.FromPtr(field))
			}
		}
		ft.TypeParams = typeParamList
	}

	if expr.Params != nil {
		paramsList := make([]Field, 0, 10)
		for _, v := range expr.Params.List {
			field := ParseField(v, af)
			if field != nil {
				paramsList = append(paramsList, lo.FromPtr(field))
			}
		}
		ft.Params = paramsList

	}

	if expr.Results != nil {
		resultsList := make([]Field, 0, 10)
		for _, v := range expr.Results.List {
			field := ParseField(v, af)
			if field != nil {
				resultsList = append(resultsList, lo.FromPtr(field))
			}
		}
		ft.Results = resultsList
	}

	return ft
}
