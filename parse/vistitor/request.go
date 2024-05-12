package vistitor

import (
	"caseGenerator/generate"
	"caseGenerator/parse"
	"fmt"
	"go/ast"
	"strconv"
)

type RequestVisitor struct {
}

func (v *RequestVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}
	if list, ok := n.(*ast.FieldList); ok {
		dbs := make([]generate.RequestDetail, 0, 10)
		// 1. 解析
		for i, requestParam := range list.List {
			// "_" 这种不处理了
			var db generate.RequestDetail
			if requestParam.Names == nil && requestParam.Type != nil {
				db.RequestName = "param" + strconv.Itoa(i)
				// todo typeParam
				result := parse.ParamParse(requestParam.Type, db.RequestName)
				if result == nil {
					fmt.Println(result)
				}
				db.RequestType = result.ParamType
				db.RequestValue = result.ParamInitValue
				db.ImportPkgPath = result.ImportPkgPath
				db.IsEllipsis = result.IsEllipsis
				dbs = append(dbs, db)
				continue
			}

			names := requestParam.Names
			for j, name := range names {
				if name.Name == "_" {
					db.RequestName = "param" + strconv.Itoa(i) + strconv.Itoa(j)
				} else {
					db.RequestName = name.Name
				}
				// todo typeParamMap
				result := parse.ParamParse(requestParam.Type, name.Name)
				if result == nil {
					fmt.Println(result)
				}
				db.RequestType = result.ParamType
				db.RequestValue = result.ParamInitValue
				db.ImportPkgPath = result.ImportPkgPath
				db.IsEllipsis = result.IsEllipsis
				dbs = append(dbs, db)
			}
		}
		// 2. 开始处理数据，request信息，名称和类型
		if len(dbs) > 0 {
			for _, db := range dbs {
				parse.AppendRequestDetailToList(db)
			}
		}

	}

	return v
}
