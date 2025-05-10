package test

type Receiver1 struct {
	Name   string
	Age    int
	IsPass bool
}

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
