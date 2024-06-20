package parser

import (
	"caseGenerator/common/utils"
	"caseGenerator/parse/visitor_v2"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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
func (p Package) ParsePackage(fileDir string) {
	// 1. 遍历文件夹，找到所有私有方法的调用
	// 2. 遍历文件夹，找到所有的私有变量的调用
	// 3. 解析出module名
	p.ModuleName = utils.GetModulePath()
	// 4. 解析packageName

	// 5. packagePath就是文件目录
	p.PackagePath = fileDir
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
					param, err := extractFileFunction(funcDecl, path)
					if err != nil {
						return false
					}
					linkedStr := param.GenerateMockGoLinked()
					// 组装map的key, key是每个方法的目录+funcName
					key := privateFunctionLinked.GenerateKey(path, param.funcName)
					privateFunctionLinked.AddLinkedInfo(key, linkedStr)
				}
				return true
			})
		}
		return nil
	})
	return privateFunctionLinked, err
}

func (s *SourceInfo) extractFileFunction(funcDecl *ast.FuncDecl, filepath string) (functionParam FunctionDeclare, err error) {
	defer func() {
		if err := recover(); err != nil {
			_ = fmt.Errorf("extractFile parse error: %s", err)
		}
	}()
	// 1. 解析request列表
	vReq := visitor_v2.Request{}
	vReq.Parse(funcDecl.Type.Params)
	// 2. 解析response列表
	vResp := visitor_v2.Response{}
	vResp.Parse(funcDecl.Type.Results)
	// 3. 解析function name
	funcName := funcDecl.Name

	functionParam = FunctionDeclare{
		RequestList:    vReq.RequestList,
		responseList:   vResp.ResponseList,
		funcName:       funcName.Name,
		moduleName:     utils.GetModulePath(),
		filepath:       filepath,
		ImportPkgPaths: nil,
	}

	return
}
