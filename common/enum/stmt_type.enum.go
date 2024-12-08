package enum

import (
	"encoding/json"
)

type StmtType BaseEnum

var (
	STMT_TYPE_ASSIGN       = StmtType{Name: "LogicExpression", Desc: "赋值类型"}
	STMT_TYPE_BLOCK        = StmtType{Name: "Block", Desc: "代码块类型"}
	STMT_TYPE_BRANCH       = StmtType{Name: "Branch", Desc: "分支语句"}
	STMT_TYPE_CASECLAUSE   = StmtType{Name: "CaseClause", Desc: "switch语句中的一个case子句"}
	STMT_TYPE_COMMONCLAUSE = StmtType{Name: "CommonClause", Desc: "select语句中的case子句"}
	STMT_TYPE_DECL         = StmtType{Name: "Decl", Desc: "声明语句"}
	STMT_TYPE_DEFER        = StmtType{Name: "Defer", Desc: "延迟语句"}
	STMT_TYPE_EMPTY        = StmtType{Name: "Empty", Desc: "空语句"}
	STMT_TYPE_EXPR         = StmtType{Name: "Expr", Desc: "表达式语句"}
	STMT_TYPE_FOR          = StmtType{Name: "For", Desc: "循环语句"}
	STMT_TYPE_GO           = StmtType{Name: "Go", Desc: "协程语句"}
	STMT_TYPE_IF           = StmtType{Name: "If", Desc: "条件语句"}
	STMT_TYPE_INCDEC       = StmtType{Name: "IncDec", Desc: "自增、自减语句"}
	STMT_TYPE_LABELED      = StmtType{Name: "Labeled", Desc: "标签语句"}
	STMT_TYPE_RANGE        = StmtType{Name: "Range", Desc: "范围语句"}
	STMT_TYPE_RETURN       = StmtType{Name: "Return", Desc: "返回语句"}
	STMT_TYPE_SELECT       = StmtType{Name: "Select", Desc: "Select语句"}
	STMT_TYPE_SEND         = StmtType{Name: "Send", Desc: "发送语句"}
	STMT_TYPE_SWITCH       = StmtType{Name: "Switch", Desc: "Switch语句"}
	STMT_TYPE_TYPESWITCH   = StmtType{Name: "TypeSwitch", Desc: "TypeSwitch语句"}
)

var ALL_STMT_TYPE = map[string]StmtType{
	STMT_TYPE_ASSIGN.Name:       STMT_TYPE_ASSIGN,
	STMT_TYPE_BLOCK.Name:        STMT_TYPE_BLOCK,
	STMT_TYPE_BRANCH.Name:       STMT_TYPE_BRANCH,
	STMT_TYPE_CASECLAUSE.Name:   STMT_TYPE_CASECLAUSE,
	STMT_TYPE_COMMONCLAUSE.Name: STMT_TYPE_COMMONCLAUSE,
	STMT_TYPE_DECL.Name:         STMT_TYPE_DECL,
	STMT_TYPE_DEFER.Name:        STMT_TYPE_DEFER,
	STMT_TYPE_EMPTY.Name:        STMT_TYPE_EMPTY,
	STMT_TYPE_EXPR.Name:         STMT_TYPE_EXPR,
	STMT_TYPE_FOR.Name:          STMT_TYPE_FOR,
	STMT_TYPE_GO.Name:           STMT_TYPE_GO,
	STMT_TYPE_IF.Name:           STMT_TYPE_IF,
	STMT_TYPE_INCDEC.Name:       STMT_TYPE_INCDEC,
	STMT_TYPE_LABELED.Name:      STMT_TYPE_LABELED,
	STMT_TYPE_RANGE.Name:        STMT_TYPE_RANGE,
	STMT_TYPE_RETURN.Name:       STMT_TYPE_RETURN,
	STMT_TYPE_SELECT.Name:       STMT_TYPE_SELECT,
	STMT_TYPE_SEND.Name:         STMT_TYPE_SEND,
	STMT_TYPE_SWITCH.Name:       STMT_TYPE_SWITCH,
	STMT_TYPE_TYPESWITCH.Name:   STMT_TYPE_TYPESWITCH,
}

func (at *StmtType) MarshalJSON() ([]byte, error) {
	// 自定义序列化逻辑
	return json.Marshal(at.Name)
}

func (at *StmtType) UnmarshalJSON(data []byte) error {
	// 自定义序列化逻辑
	var name string
	err := json.Unmarshal(data, &name)
	if err != nil {
		return err
	}
	assignmentType, ok := ALL_STMT_TYPE[name]
	if ok {
		at.Name = assignmentType.Name
		at.Desc = assignmentType.Desc
	}
	return nil
}
