package dict

type ExampleDict struct {
	Name     string
	Age      int32
	IsDelete bool
}

type ReceiverDict struct {
}

// TestReceiverFunc 用来模拟Receive方法
func (r *ReceiverDict) TestReceiverFunc(v string) string {
	return v
}
