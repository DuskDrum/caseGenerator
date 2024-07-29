package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// SourceInfoExt 源信息
type SourceInfoExt struct {
	// 包信息
	Package
	// 包中的文件信息
	SourceInfo []*SourceInfo
}

// SourceInfo 这是解析包里每个文件得到的信息
type SourceInfo struct {
	// 依赖信息
	Import
	// 方法信息
	FunctionDeclare
	// 赋值信息列表
	AssignmentList []*Assignment
	// 条件语句列表
	ConditionList []*ConditionNode
}

// ParseSource 主入口，根据func解析出所有内容，并返回出来
func (s *SourceInfoExt) ParseSource(fileDir string) {
	// 1. 解析package相关信息, 其中包含了整个包下面的私有方法和私有变量
	// 2. 解析fileDir，遍历其中的文件，并解析文件中的方法列表
	// 要保证这个fileDir是最底层的目录，这样才能保证其下文件中的包名是一样的
	sourceInfoList := make([]*SourceInfo, 0, 10)
	parseGoFile(fileDir, func(funcDecl *ast.FuncDecl, filepath string) {
		var sourceInfo SourceInfo
		s.Package = sourceInfo.ParsePackage(fileDir)

		// 3. 解析方法信息，其中包含了方法的入参、响应列表、泛型、receiver等信息
		// 会跳过"_test.go"的文件解析
		// filepath 是以".go"结尾的全路径
		functionDeclare := sourceInfo.extractFileFunction(funcDecl, filepath)
		if functionDeclare != nil {
			sourceInfo.FunctionDeclare = lo.FromPtr(functionDeclare)
		}
		conditionNodes := make([]*ConditionNode, 0, 10)
		assignmentList := make([]*Assignment, 0, 10)
		for _, body := range funcDecl.Body.List {
			// 3. 遍历解析方法中所有赋值
			assignment := sourceInfo.ParseAssignment(body)
			if len(assignment) > 0 {
				assignmentList = append(assignmentList, assignment...)
			}
			// 4. 遍历解析方法中的所有条件
			switch node := body.(type) {
			case *ast.IfStmt, *ast.SwitchStmt:
				conditionNode := sourceInfo.parseCondition(node)
				if conditionNode != nil {
					conditionNodes = append(conditionNodes, conditionNode)
				}
			}
		}
		sourceInfo.ConditionList = conditionNodes
		sourceInfo.AssignmentList = assignmentList

		sourceInfoList = append(sourceInfoList, &sourceInfo)
	})
	s.SourceInfo = sourceInfoList
}

func parseGoFile(path string, funcDeclParse func(funcDecl *ast.FuncDecl, filename string), excludedPaths ...string) {
	// walk的处理目录下的每个go文件
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 根据入参 过滤某些路径
		for _, p := range excludedPaths {
			if p == path {
				return nil
			}
		}
		// 跳过目录的处理
		if info.IsDir() {
			//return filepath.SkipDir
			return nil
		}
		// 跳过不是go文件的处理， 跳过测试文件
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// 开始解析文件
		var fset = token.NewFileSet()
		var f *ast.File
		if f, err = parser.ParseFile(fset, path, nil, parser.ParseComments); err != nil {
			return nil
		}
	OuterLoop:
		for _, cg := range f.Decls {
			// 只处理funcDecl类型的ast
			decl, ok := cg.(*ast.FuncDecl)
			if !ok {
				continue
			}
			// 标记，解析的时候可以跳过
			if decl.Doc != nil {
				for _, com := range decl.Doc.List {
					// 跳过处理
					if strings.Contains(com.Text, "skipcodeautoparse") {
						fmt.Print("skip auto parse...")
						continue OuterLoop
					}
				}
			}
			// 处理每个方法的解析
			funcDeclParse(decl, path)
		}
		return nil
	})
	if err != nil {
		panic(err.Error())
	}
}
