package visitor_v2

import (
	"caseGenerator/generate"
	"go/ast"
)

type Response struct {
	ResponseList []generate.ResponseDetail
}

func (r Response) Parse(list *ast.FieldList) {
	if list == nil {
		return
	}
	fields := list.List
	if len(fields) > 0 {
		responseDetails := make([]generate.ResponseDetail, 0, 10)

		for _, v := range fields {
			if len(v.Names) > 0 {
				for _, name := range v.Names {
					param := ParseParamWithoutInit(v.Type, name.Name)
					responseDetails = append(responseDetails, generate.ResponseDetail{
						ParamName:       param.ParamName,
						ParamType:       param.ParamType,
						ParamInitValue:  param.ParamInitValue,
						ParamCheckValue: param.ParamCheckValue,
						IsEllipsis:      param.IsEllipsis,
					})
				}
			} else if v.Type != nil {
				param := ParseParamWithoutInit(v.Type, "param")
				responseDetails = append(responseDetails, generate.ResponseDetail{
					ParamName:       param.ParamName,
					ParamType:       param.ParamType,
					ParamInitValue:  param.ParamInitValue,
					ParamCheckValue: param.ParamCheckValue,
					IsEllipsis:      param.IsEllipsis,
				})
			}

		}
		r.ResponseList = responseDetails

	}

}
