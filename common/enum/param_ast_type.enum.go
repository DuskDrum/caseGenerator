package enum

type ParamAstType struct {
	Name string
}

var (
	PARAM_AST_TYPE_1 = ParamAstType{Name: "1"}
	PARAM_AST_TYPE_2 = ParamAstType{Name: "2"}
	PARAM_AST_TYPE_3 = ParamAstType{Name: "3"}
)

var ALL_PARAM_AST_TYPE = map[string]ParamAstType{
	PARAM_AST_TYPE_1.Name: PARAM_AST_TYPE_1,
	PARAM_AST_TYPE_2.Name: PARAM_AST_TYPE_2,
	PARAM_AST_TYPE_3.Name: PARAM_AST_TYPE_3,
}
