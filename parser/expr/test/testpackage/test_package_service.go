package testpackage

type TestService struct{}

func (TestService) Add(a, b int) int {
	return a + b
}
