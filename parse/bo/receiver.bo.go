package bo

import (
	"caseGenerator/utils"
	"encoding/json"
	"fmt"
)

// ReceiverInfo 接受者，method独有。 receive有很多特性，所以直接定义为InvocationUnary
type ReceiverInfo struct {
	ReceiverName  string
	ReceiverValue InvocationUnary
}

func (r *ReceiverInfo) GetParamName() string {
	return r.ReceiverName
}

func (r *ReceiverInfo) UnmarshalerInfo(jsonString string) {
	dat := utils.Empty[InvocationUnary]()
	if err := json.Unmarshal([]byte(jsonString), &dat); err == nil {
		fmt.Println(dat)
		r.ReceiverValue = dat
	} else {
		fmt.Println(jsonString)
	}
}
