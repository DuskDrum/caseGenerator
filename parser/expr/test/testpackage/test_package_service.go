package testpackage

type TestService struct{}

func (TestService) Add(a, b int) int {
	return a + b
}

func (*TestService) PtrAdd(a, b int) int {
	return a + b
}
