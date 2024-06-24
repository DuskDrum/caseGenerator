package parser

import (
	"caseGenerator/common/utils"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

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

// ParsePackage 解析包的基础信息
func (p Package) ParsePackage(si *SourceInfo, fileDir string) {
	// 1. 遍历文件夹，找到所有私有方法的调用
	funcMap, err := si.extractPrivateFile(fileDir)
	if err != nil {
		panic("extract private file error!")
	}
	p.PrivateFunction = funcMap
	// 2. 遍历文件夹，找到所有的私有变量的调用
	assignmentMap, err := si.extractPrivateAssignment(fileDir)
	if err != nil {
		panic("extract private assignment error!")
	}
	p.PrivateParam = assignmentMap
	// 3. 解析出module名
	p.ModuleName = utils.GetModulePath()
	// 4. 解析packageName
	var packageName string
	index := strings.LastIndex(fileDir, "/")
	if index > 0 {
		packageName = fileDir[:index]
	} else {
		packageName = ""
	}
	p.PackageName = packageName
	// 5. packagePath就是文件目录
	p.PackagePath = fileDir
}

// 遍历文件夹，找到所有私有变量的详情
func (s *SourceInfo) extractPrivateAssignment(fileDir string) (map[string]Assignment, error) {
	funcMap := make(map[string]Assignment, 10)
	// 1. 遍历解析文件夹
	err := filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Parse the file
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}

			// Traverse the AST
			ast.Inspect(file, func(node ast.Node) bool {
				// Check if the node is a function declaration
				//funcDecl, ok := node.(*ast.AssignStmt)
				//if !ok {
				//	return true
				//}
				// 判断是否是私有方法
				//if unicode.IsLower([]rune(funcDecl.Name.Name)[0]) {
				//	// 1. 解析key， key 是 path + "." + funcDecl.Name.Name
				//	key := path + "." + funcDecl.Name.Name
				//
				//	functionParam, err := s.extractPrivateAssignment(funcDecl, path)
				//	if err != nil {
				//		return false
				//	}
				//	funcMap[key] = functionParam
				//}
				return true
			})
		}
		return nil
	})
	return funcMap, err
}

// extractAssignmentFunction filepath 应该以 .go 后缀结尾
func (s *SourceInfo) extractAssignmentFunction(funcDecl *ast.FuncDecl, filepath string) (functionParam FunctionDeclare, err error) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("extractFile parse error: %s", err)
		}
	}()
	if !strings.HasSuffix(filepath, ".go") {
		panic("filepath should be suffix with .go")
	}
	// 1. 解析request列表、response列表
	reqList, respList := s.ParseFuncTypeParamParseResult(funcDecl.Type)
	// 2. 解析泛型
	// todo 注意私有方法的泛型调用
	genericsMap := s.ParseGenericsMap(funcDecl)
	// 3. 解析receiver
	receiver := s.ParseReceiver(funcDecl)
	// 4. 将filePath 转为.go文件和包名目录
	var fileName, functionPath string
	index := strings.LastIndex(filepath, "/")
	if index > 0 {
		fileName = filepath[index+1:]
		functionPath = filepath[:index]
	} else {
		fileName = filepath
		functionPath = ""
	}
	functionParam = FunctionDeclare{
		RequestList:  reqList,
		ResponseList: respList,
		FunctionName: funcDecl.Name.Name,
		FunctionBasic: FunctionBasic{
			FunctionPath: functionPath,
			FileName:     fileName,
			GenericsMap:  genericsMap,
			Receiver:     receiver,
		},
	}
	return
}

// 遍历文件夹，找到所有私有方法的调用
func (s *SourceInfo) extractPrivateFile(fileDir string) (map[string]FunctionDeclare, error) {
	funcMap := make(map[string]FunctionDeclare, 10)
	// 1. 遍历解析文件夹
	err := filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Parse the file
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
			if err != nil {
				return err
			}

			// Traverse the AST
			ast.Inspect(file, func(node ast.Node) bool {
				// Check if the node is a function declaration
				funcDecl, ok := node.(*ast.FuncDecl)
				if !ok {
					return true
				}
				// 判断是否是私有方法
				if unicode.IsLower([]rune(funcDecl.Name.Name)[0]) {
					// 1. 解析key， key 是 path + "." + funcDecl.Name.Name
					key := path + "." + funcDecl.Name.Name

					functionParam, err := s.extractFileFunction(funcDecl, path)
					if err != nil {
						return false
					}
					funcMap[key] = functionParam
				}
				return true
			})
		}
		return nil
	})
	return funcMap, err
}

// extractFileFunction filepath 应该以 .go 后缀结尾
func (s *SourceInfo) extractFileFunction(funcDecl *ast.FuncDecl, filepath string) (functionParam FunctionDeclare, err error) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("extractFile parse error: %s", err)
		}
	}()
	if !strings.HasSuffix(filepath, ".go") {
		panic("filepath should be suffix with .go")
	}
	// 1. 解析request列表、response列表
	reqList, respList := s.ParseFuncTypeParamParseResult(funcDecl.Type)
	// 2. 解析泛型
	// todo 注意私有方法的泛型调用
	genericsMap := s.ParseGenericsMap(funcDecl)
	// 3. 解析receiver
	receiver := s.ParseReceiver(funcDecl)
	// 4. 将filePath 转为.go文件和包名目录
	var fileName, functionPath string
	index := strings.LastIndex(filepath, "/")
	if index > 0 {
		fileName = filepath[index+1:]
		functionPath = filepath[:index]
	} else {
		fileName = filepath
		functionPath = ""
	}
	functionParam = FunctionDeclare{
		RequestList:  reqList,
		ResponseList: respList,
		FunctionName: funcDecl.Name.Name,
		FunctionBasic: FunctionBasic{
			FunctionPath: functionPath,
			FileName:     fileName,
			GenericsMap:  genericsMap,
			Receiver:     receiver,
		},
	}
	return
}
