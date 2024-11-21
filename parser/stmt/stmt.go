package stmt

//
//// Stmt 参数接口类型，将参数需要的方法定义出来
//type Stmt interface {
//	// Express 生成表达式. key为其变量name
//	Express() []map[string]StmtExpression
//}
//
//// StmtExpression stmt的表达式，记录了参数的变动, 参数也可以直接重新赋值
//type StmtExpression struct {
//	Name      string
//	InitParam *_struct.Parameter
//	Type      enum.StmtType
//	// 参数变动列表
//	ExpressionExpr        string   // "a > 10"
//	ExpressionElementList []string // ["a",">","10"]
//	ExpressionMap         map[string]_struct.Parameter
//}
