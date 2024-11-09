package expr

import (
	"caseGenerator/common/enum"
	"go/ast"

	"github.com/samber/lo"
)

// Struct 构体类型定义。
// 包含结构体的字段列表，描述了结构体的整体结构信息
// 一整个 struct，比如 type XX struct{}是放在 ast.TypeSpec 里的， 名字是 XX
// 他的 Type 是  ast.StructType， structType是一个 field 列表， 里面是 name+type
type Struct struct {
	FieldList []Field
}

func (s *Struct) GetType() enum.ParameterType {
	return enum.PARAMETER_TYPE_STRUCT
}

func (s *Struct) GetFormula() string {
	formula := "struct{\n"
	for _, v := range s.FieldList {
		formula += v.GetFormula() + "\n"
	}
	formula += "}"
	return ""
}

// ParseStruct 解析ast
func ParseStruct(expr *ast.StructType) *Struct {
	s := &Struct{}
	fields := expr.Fields
	fieldList := make([]Field, 0, 10)
	if fields != nil {
		for _, field := range fields.List {
			pf := ParseField(field)
			if pf != nil {
				fieldList = append(fieldList, lo.FromPtr(pf))
			}
		}
	}
	s.FieldList = fieldList
	return s
}
