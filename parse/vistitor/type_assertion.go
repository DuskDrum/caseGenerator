package vistitor

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"github.com/samber/lo"
	"go/ast"
	"log"
)

type TypeAssertionVisitor struct {
	TypeAssertionSlice []TypeAssertionMapBo
}

type TypeAssertionMapBo struct {
	Key                  string
	TypeAssertionContent *bo.ParamParseResult
}

// AddTypeAssertionSlice todo 单例模式
func (v *TypeAssertionVisitor) AddTypeAssertionSlice(key string, value *bo.ParamParseResult) {
	if v.TypeAssertionSlice == nil {
		v.TypeAssertionSlice = make([]TypeAssertionMapBo, 0, 10)
	}
	mapBo := TypeAssertionMapBo{
		Key:                  key,
		TypeAssertionContent: value,
	}
	v.TypeAssertionSlice = append(v.TypeAssertionSlice, mapBo)
}

func (v *TypeAssertionVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	switch node := n.(type) {
	case *ast.TypeAssertExpr:
		identName := node.X.(*ast.Ident)
		if identName == nil {
			log.Fatalf("未成功解析出节点的名称...")
		}
		// 获取 typeParam
		parse := ParamParse(node.Type, "")
		v.AddTypeAssertionSlice(identName.Name, parse)
	}

	return v
}

// CombinationTypeAssertionRequest 排列组合所有类型断言的可能性
func (v *TypeAssertionVisitor) CombinationTypeAssertionRequest(reqList []generate.RequestDetail) [][]generate.RequestDetail {
	// 1. 先按照key转为map<key,slice>
	typeAssertionRequestDetail := lo.Map(v.TypeAssertionSlice, func(item TypeAssertionMapBo, index int) generate.RequestDetail {
		return generate.RequestDetail{
			RequestName:  item.TypeAssertionContent.ParamName,
			RequestType:  item.TypeAssertionContent.ParamType,
			RequestValue: item.TypeAssertionContent.ParamInitValue,
			IsEllipsis:   false,
		}
	})
	combinations := make([][]generate.RequestDetail, 0, 10)
	// 2. 排列组合, 使用嵌套循环
	for _, v := range typeAssertionRequestDetail {
		detail := make([]generate.RequestDetail, 0, 10)
		for _, v2 := range reqList {
			if v.RequestName == v2.RequestName {
				detail = append(detail, v)
			} else {
				detail = append(detail, v2)
			}
		}
		combinations = append(combinations, detail)
	}
	return combinations
}
