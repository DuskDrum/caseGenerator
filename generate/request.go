package generate

import (
	"caseGenerator/common/constants"
	"caseGenerator/common/enum"
	"caseGenerator/common/utils"
	"caseGenerator/parser"
	"strings"
)

// Request 记录要生成的请求信息
type Request struct {
	RequestCaseList []RequestCase
	// ReceiverName  代表就是单测中receiver的名字，默认是`convert`
	ReceiverName string
	// ReceiverType 代表就是单测中receiver的类型，比如说是`某个struct`、`int`等
	ReceiverType string
	// ReceiverValue 代表了单测中要设置的值，默认是用`lo.Empty[int]`
	ReceiverValue string
}

// RequestCase 可能一个请求想有多种返回，比如说想给请求设置lo.Empty[]、又想给请求设置fakeit、又想自定义一些字段，这是可以有多个RequestCase
type RequestCase struct {
	RequestList       []RequestDetail
	RequestNameString string
}

type RequestDetail struct {
	Name  string
	Value string
}

// 记录receiver，方法名，请求参数
func (r *Request) generateRequest(si *parser.SourceInfo) {
	// 1. 组装receiver
	receiver := si.FunctionDeclare.Receiver
	r.ReceiverName = constants.RECEIVER_NAME
	r.ReceiverType = receiver.Type
	r.ReceiverValue = utils.LoEmpty(receiver.Type)
	// 2. 组装lo.Empty[]的请求参数
	empty := RequestCase{
		RequestList: make([]RequestDetail, 0, 10),
	}
	var emptyRequestNameString string
	for _, v := range si.RequestList {
		empty.RequestList = append(empty.RequestList, RequestDetail{
			Name:  v.Name,
			Value: utils.LoEmpty(v.Type),
		})
		if v.AstType == enum.PARAM_AST_TYPE_Ellipsis {
			emptyRequestNameString += "tt.args." + v.Name + "... , "
		} else {
			emptyRequestNameString += "tt.args." + v.Name + ", "
		}
	}
	emptyRequestNameString = strings.TrimRight(emptyRequestNameString, ", ")
	empty.RequestNameString = emptyRequestNameString
	// 3. 组装fakeit的请求参数
	fakeit := RequestCase{
		RequestList: make([]RequestDetail, 0, 10),
	}
	var fakeitRequestNameString string
	for _, v := range si.RequestList {
		fakeit.RequestList = append(empty.RequestList, RequestDetail{
			Name:  v.Name,
			Value: "gofakeit.Struct(" + v.Type + ")",
		})
		if v.AstType == enum.PARAM_AST_TYPE_Ellipsis {
			fakeitRequestNameString += "tt.args." + v.Name + "... , "
		} else {
			fakeitRequestNameString += "tt.args." + v.Name + ", "
		}
	}
	fakeitRequestNameString = strings.TrimRight(fakeitRequestNameString, ", ")
	fakeit.RequestNameString = fakeitRequestNameString
	// 4. 暂时只有这两种类型，直接组装为request列表
	caseList := make([]RequestCase, 0, 10)
	caseList = append(caseList, empty)
	caseList = append(caseList, fakeit)
}
