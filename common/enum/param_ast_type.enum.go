package enum

import "encoding/json"

type ParamAstType struct {
	Name string
	Desc string
}

var (
	PARAM_AST_TYPE_Comment         = ParamAstType{Name: "Comment", Desc: "表示一行注释 // 或者 / /"}
	PARAM_AST_TYPE_CommentGroup    = ParamAstType{Name: "CommentGroup", Desc: "表示多行注释"}
	PARAM_AST_TYPE_Field           = ParamAstType{Name: "Field", Desc: "表示结构体中的一个定义或者变量,或者函数签名当中的参数或者返回值"}
	PARAM_AST_TYPE_FieldList       = ParamAstType{Name: "FieldList", Desc: "表示以”{}”或者”()”包围的Filed列表"}
	PARAM_AST_TYPE_BadExpr         = ParamAstType{Name: "BadExpr", Desc: "用来表示错误表达式的占位符"}
	PARAM_AST_TYPE_Ident           = ParamAstType{Name: "Ident", Desc: "比如报名,函数名,变量名"}
	PARAM_AST_TYPE_Ellipsis        = ParamAstType{Name: "Ellipsis", Desc: "省略号表达式,比如参数列表的最后一个可以写成arg..."}
	PARAM_AST_TYPE_BasicLit        = ParamAstType{Name: "BasicLit", Desc: "基本字面值,数字或者字符串"}
	PARAM_AST_TYPE_FuncLit         = ParamAstType{Name: "FuncLit", Desc: "函数定义"}
	PARAM_AST_TYPE_ParenExpr       = ParamAstType{Name: "ParenExpr", Desc: "括号表达式,被括号包裹的表达式"}
	PARAM_AST_TYPE_StarExpr        = ParamAstType{Name: "StarExpr", Desc: "指针类型"}
	PARAM_AST_TYPE_UnaryExpr       = ParamAstType{Name: "UnaryExpr", Desc: "一元表达式"}
	PARAM_AST_TYPE_BinaryExpr      = ParamAstType{Name: "BinaryExpr", Desc: "二元表达式"}
	PARAM_AST_TYPE_KeyValueExp     = ParamAstType{Name: "KeyValueExp", Desc: "key:value"}
	PARAM_AST_TYPE_ArrayType       = ParamAstType{Name: "ArrayType", Desc: "数组类型"}
	PARAM_AST_TYPE_StructType      = ParamAstType{Name: "StructType", Desc: "结构体类型"}
	PARAM_AST_TYPE_FuncType        = ParamAstType{Name: "FuncType", Desc: "函数类型"}
	PARAM_AST_TYPE_InterfaceType   = ParamAstType{Name: "InterfaceType", Desc: "接口类型"}
	PARAM_AST_TYPE_MapType         = ParamAstType{Name: "MapType", Desc: "map类型"}
	PARAM_AST_TYPE_ChanType        = ParamAstType{Name: "ChanType", Desc: "管道类型"}
	PARAM_AST_TYPE_BadStmt         = ParamAstType{Name: "BadStmt", Desc: "错误的语句"}
	PARAM_AST_TYPE_EmptyStmt       = ParamAstType{Name: "EmptyStmt", Desc: "空语句"}
	PARAM_AST_TYPE_ExprStmt        = ParamAstType{Name: "ExprStmt", Desc: "包含单独的表达式语句"}
	PARAM_AST_TYPE_SendStmt        = ParamAstType{Name: "SendStmt", Desc: "chan发送语句"}
	PARAM_AST_TYPE_IncDecStmt      = ParamAstType{Name: "IncDecStmt", Desc: "自增或者自减语句"}
	PARAM_AST_TYPE_GoStmt          = ParamAstType{Name: "GoStmt", Desc: "Go语句"}
	PARAM_AST_TYPE_DeferStmt       = ParamAstType{Name: "DeferStmt", Desc: "延迟语句"}
	PARAM_AST_TYPE_ReturnStmt      = ParamAstType{Name: "ReturnStmt", Desc: "语句"}
	PARAM_AST_TYPE_BranchStmt      = ParamAstType{Name: "BranchStmt", Desc: "continue"}
	PARAM_AST_TYPE_BlockStmt       = ParamAstType{Name: "BlockStmt", Desc: "包裹"}
	PARAM_AST_TYPE_IfStmt          = ParamAstType{Name: "IfStmt", Desc: "语句"}
	PARAM_AST_TYPE_TypeSwitchStmt  = ParamAstType{Name: "TypeSwitchStmt", Desc: "x:=y.(type)"}
	PARAM_AST_TYPE_CommClause      = ParamAstType{Name: "CommClause", Desc: "<-:"}
	PARAM_AST_TYPE_SelectStmt      = ParamAstType{Name: "SelectStmt", Desc: "语句"}
	PARAM_AST_TYPE_ForStmt         = ParamAstType{Name: "ForStmt", Desc: "语句"}
	PARAM_AST_TYPE_RangeStmt       = ParamAstType{Name: "RangeStmt", Desc: "语句"}
	PARAM_AST_TYPE_Spec            = ParamAstType{Name: "Spec", Desc: "type"}
	PARAM_AST_TYPE_Import          = ParamAstType{Name: "Import", Desc: "Spec"}
	PARAM_AST_TYPE_Value           = ParamAstType{Name: "Value", Desc: "Spec"}
	PARAM_AST_TYPE_Type            = ParamAstType{Name: "Type", Desc: "Spec"}
	PARAM_AST_TYPE_BadDecl         = ParamAstType{Name: "BadDecl", Desc: "错误申明"}
	PARAM_AST_TYPE_GenDecl         = ParamAstType{Name: "GenDecl", Desc: "a)"}
	PARAM_AST_TYPE_FuncDecl        = ParamAstType{Name: "FuncDecl", Desc: "函数申明"}
	PARAM_AST_TYPE_CompositeLit    = ParamAstType{Name: "CompositeLit", Desc: "构造类型,比如{1,2,3,4}"}
	PARAM_AST_TYPE_SelectorExpr    = ParamAstType{Name: "SelectorExpr", Desc: "选择结构,类似于a.b的结构"}
	PARAM_AST_TYPE_IndexExpr       = ParamAstType{Name: "IndexExpr", Desc: "expr[expr]"}
	PARAM_AST_TYPE_SliceExpr       = ParamAstType{Name: "SliceExpr", Desc: "expr[low:mid:high]"}
	PARAM_AST_TYPE_CallExpr        = ParamAstType{Name: "CallExpr", Desc: "expr()"}
	PARAM_AST_TYPE_TypeAssertExpr  = ParamAstType{Name: "TypeAssertExpr", Desc: "X.(type)"}
	PARAM_AST_TYPE_DeclStmt        = ParamAstType{Name: "DeclStmt", Desc: "在语句列表里的申明"}
	PARAM_AST_TYPE_LabeledStmt     = ParamAstType{Name: "LabeledStmt", Desc: "indent:stmt"}
	PARAM_AST_TYPE_AssignStmt      = ParamAstType{Name: "AssignStmt", Desc: "赋值语句"}
	PARAM_AST_TYPE_CaseClause      = ParamAstType{Name: "CaseClause", Desc: "语句"}
	PARAM_AST_TYPE_SwitchStmt      = ParamAstType{Name: "SwitchStmt", Desc: "语句"}
	PARAM_AST_TYPE_BasicLitUnary   = ParamAstType{Name: "BasicLitUnary", Desc: "基础"}
	PARAM_AST_TYPE_InvocationUnary = ParamAstType{Name: "InvocationUnary", Desc: "常量、变量名"}
	PARAM_AST_TYPE_ParamUnary      = ParamAstType{Name: "ParamUnary", Desc: "常量的层级调用"}
	PARAM_AST_TYPE_CallUnary       = ParamAstType{Name: "CallUnary", Desc: "方法调用、也有层级调用"}
)

var ALL_PARAM_AST_TYPE = map[string]ParamAstType{
	PARAM_AST_TYPE_Comment.Name:         PARAM_AST_TYPE_Comment,
	PARAM_AST_TYPE_CommentGroup.Name:    PARAM_AST_TYPE_CommentGroup,
	PARAM_AST_TYPE_Field.Name:           PARAM_AST_TYPE_Field,
	PARAM_AST_TYPE_FieldList.Name:       PARAM_AST_TYPE_FieldList,
	PARAM_AST_TYPE_BadExpr.Name:         PARAM_AST_TYPE_BadExpr,
	PARAM_AST_TYPE_Ident.Name:           PARAM_AST_TYPE_Ident,
	PARAM_AST_TYPE_Ellipsis.Name:        PARAM_AST_TYPE_Ellipsis,
	PARAM_AST_TYPE_BasicLit.Name:        PARAM_AST_TYPE_BasicLit,
	PARAM_AST_TYPE_FuncLit.Name:         PARAM_AST_TYPE_FuncLit,
	PARAM_AST_TYPE_CompositeLit.Name:    PARAM_AST_TYPE_CompositeLit,
	PARAM_AST_TYPE_ParenExpr.Name:       PARAM_AST_TYPE_ParenExpr,
	PARAM_AST_TYPE_SelectorExpr.Name:    PARAM_AST_TYPE_SelectorExpr,
	PARAM_AST_TYPE_IndexExpr.Name:       PARAM_AST_TYPE_IndexExpr,
	PARAM_AST_TYPE_SliceExpr.Name:       PARAM_AST_TYPE_SliceExpr,
	PARAM_AST_TYPE_TypeAssertExpr.Name:  PARAM_AST_TYPE_TypeAssertExpr,
	PARAM_AST_TYPE_CallExpr.Name:        PARAM_AST_TYPE_CallExpr,
	PARAM_AST_TYPE_StarExpr.Name:        PARAM_AST_TYPE_StarExpr,
	PARAM_AST_TYPE_UnaryExpr.Name:       PARAM_AST_TYPE_UnaryExpr,
	PARAM_AST_TYPE_BinaryExpr.Name:      PARAM_AST_TYPE_BinaryExpr,
	PARAM_AST_TYPE_KeyValueExp.Name:     PARAM_AST_TYPE_KeyValueExp,
	PARAM_AST_TYPE_ArrayType.Name:       PARAM_AST_TYPE_ArrayType,
	PARAM_AST_TYPE_StructType.Name:      PARAM_AST_TYPE_StructType,
	PARAM_AST_TYPE_FuncType.Name:        PARAM_AST_TYPE_FuncType,
	PARAM_AST_TYPE_InterfaceType.Name:   PARAM_AST_TYPE_InterfaceType,
	PARAM_AST_TYPE_MapType.Name:         PARAM_AST_TYPE_MapType,
	PARAM_AST_TYPE_ChanType.Name:        PARAM_AST_TYPE_ChanType,
	PARAM_AST_TYPE_BadStmt.Name:         PARAM_AST_TYPE_BadStmt,
	PARAM_AST_TYPE_DeclStmt.Name:        PARAM_AST_TYPE_DeclStmt,
	PARAM_AST_TYPE_EmptyStmt.Name:       PARAM_AST_TYPE_EmptyStmt,
	PARAM_AST_TYPE_LabeledStmt.Name:     PARAM_AST_TYPE_LabeledStmt,
	PARAM_AST_TYPE_ExprStmt.Name:        PARAM_AST_TYPE_ExprStmt,
	PARAM_AST_TYPE_SendStmt.Name:        PARAM_AST_TYPE_SendStmt,
	PARAM_AST_TYPE_IncDecStmt.Name:      PARAM_AST_TYPE_IncDecStmt,
	PARAM_AST_TYPE_AssignStmt.Name:      PARAM_AST_TYPE_AssignStmt,
	PARAM_AST_TYPE_GoStmt.Name:          PARAM_AST_TYPE_GoStmt,
	PARAM_AST_TYPE_DeferStmt.Name:       PARAM_AST_TYPE_DeferStmt,
	PARAM_AST_TYPE_ReturnStmt.Name:      PARAM_AST_TYPE_ReturnStmt,
	PARAM_AST_TYPE_BranchStmt.Name:      PARAM_AST_TYPE_BranchStmt,
	PARAM_AST_TYPE_BlockStmt.Name:       PARAM_AST_TYPE_BlockStmt,
	PARAM_AST_TYPE_IfStmt.Name:          PARAM_AST_TYPE_IfStmt,
	PARAM_AST_TYPE_CaseClause.Name:      PARAM_AST_TYPE_CaseClause,
	PARAM_AST_TYPE_SwitchStmt.Name:      PARAM_AST_TYPE_SwitchStmt,
	PARAM_AST_TYPE_TypeSwitchStmt.Name:  PARAM_AST_TYPE_TypeSwitchStmt,
	PARAM_AST_TYPE_CommClause.Name:      PARAM_AST_TYPE_CommClause,
	PARAM_AST_TYPE_SelectStmt.Name:      PARAM_AST_TYPE_SelectStmt,
	PARAM_AST_TYPE_ForStmt.Name:         PARAM_AST_TYPE_ForStmt,
	PARAM_AST_TYPE_RangeStmt.Name:       PARAM_AST_TYPE_RangeStmt,
	PARAM_AST_TYPE_Spec.Name:            PARAM_AST_TYPE_Spec,
	PARAM_AST_TYPE_Import.Name:          PARAM_AST_TYPE_Import,
	PARAM_AST_TYPE_Value.Name:           PARAM_AST_TYPE_Value,
	PARAM_AST_TYPE_Type.Name:            PARAM_AST_TYPE_Type,
	PARAM_AST_TYPE_BadDecl.Name:         PARAM_AST_TYPE_BadDecl,
	PARAM_AST_TYPE_GenDecl.Name:         PARAM_AST_TYPE_GenDecl,
	PARAM_AST_TYPE_FuncDecl.Name:        PARAM_AST_TYPE_FuncDecl,
	PARAM_AST_TYPE_SelectorExpr.Name:    PARAM_AST_TYPE_SelectorExpr,
	PARAM_AST_TYPE_IndexExpr.Name:       PARAM_AST_TYPE_IndexExpr,
	PARAM_AST_TYPE_SliceExpr.Name:       PARAM_AST_TYPE_SliceExpr,
	PARAM_AST_TYPE_CallExpr.Name:        PARAM_AST_TYPE_CallExpr,
	PARAM_AST_TYPE_TypeAssertExpr.Name:  PARAM_AST_TYPE_TypeAssertExpr,
	PARAM_AST_TYPE_DeclStmt.Name:        PARAM_AST_TYPE_DeclStmt,
	PARAM_AST_TYPE_LabeledStmt.Name:     PARAM_AST_TYPE_LabeledStmt,
	PARAM_AST_TYPE_AssignStmt.Name:      PARAM_AST_TYPE_AssignStmt,
	PARAM_AST_TYPE_CaseClause.Name:      PARAM_AST_TYPE_CaseClause,
	PARAM_AST_TYPE_SwitchStmt.Name:      PARAM_AST_TYPE_SwitchStmt,
	PARAM_AST_TYPE_BasicLitUnary.Name:   PARAM_AST_TYPE_BasicLitUnary,
	PARAM_AST_TYPE_InvocationUnary.Name: PARAM_AST_TYPE_InvocationUnary,
	PARAM_AST_TYPE_ParamUnary.Name:      PARAM_AST_TYPE_ParamUnary,
	PARAM_AST_TYPE_CallUnary.Name:       PARAM_AST_TYPE_CallUnary,
}

func (at *ParamAstType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Name)
}

func (at *ParamAstType) UnmarshalJSON(data []byte) error {
	// 自定义序列化逻辑
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	assignmentType, ok := ALL_PARAM_AST_TYPE[name]
	if ok {
		at.Name = assignmentType.Name
		at.Desc = assignmentType.Desc
	}
	return nil
}
