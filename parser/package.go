package parser

// Package 包相关信息，包含包中私有的方法、私有的变量、私有的常量等。也包含包的名称、路径、module名
type Package struct {
	PackageName string
	PackagePath string
	ModuleName  string
	// key是方法名
	PrivateFunction map[string]FunctionDeclare
	// key是变量名
	PrivateParam map[string]Assignment
}
