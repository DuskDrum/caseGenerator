package generate

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
//func ()  {

//}
