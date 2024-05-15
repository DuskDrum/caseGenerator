package vistitor

import (
	"caseGenerator/generate"
	"caseGenerator/parse/bo"
	"fmt"
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
func (v *TypeAssertionVisitor) CombinationTypeAssertionRequest(reqList []generate.RequestDetail) []generate.RequestDetail {
	// 1. 先按照key转为map<key,slice>
	typeAssertionSliceMap := lo.GroupBy(v.TypeAssertionSlice, func(item TypeAssertionMapBo) string {
		return item.Key
	})
	// 2. 排列组合
	for i := range typeAssertionSliceMap {
		fmt.Print(i)
	}
	// todo 代写
	return nil
}
