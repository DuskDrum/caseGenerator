package test

type Receiver1 struct {
	Name   string
	Age    int
	IsPass bool
}

var PackageGlobalNumber = 10
var PackageGlobalString = "hello"
var PackageGlobalBool = true
var ZeroInt1 = ZeroInt
var GlobalReceiver1 = Receiver1{Name: "GlobalName"}
var GlobalReceiverName = GlobalReceiver1.Name
var GlobalReceiverAge = GlobalReceiver1.Age
var GlobalReceiverIsPass = GlobalReceiver1.IsPass
var GlobalCallName = getReceiver().Name // 有可能是call.属性
var GlobalFunction = Ident()            // 有可能是call

func (r Receiver1) add(a int, b int) int {
	result, err := r.doAddInt(a, b)
	if err != nil {
		return a + b
	}
	return result
}

func (r Receiver1) doAddInt(a int, b int) (int, error) {
	return 0, nil
}

func (r Receiver1) doAddInt8(a int, b int) (int8, error) {
	return 0, nil
}

func (r Receiver1) doAddInt16(a int, b int) (int16, error) {
	return 0, nil
}

func getReceiver() Receiver1 {
	return Receiver1{}
}
