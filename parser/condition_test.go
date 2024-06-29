package parser

//
//// TestParseAssignment 测试assignment的测试用例，方法调用类型
//func TestParseAssignment_call(t *testing.T) {
//	parseFile("../example/assignment/conv_function.go")
//}
//
//// TestParseAssignment 测试assignment的测试用例，构造类型
//func TestParseAssignment_composite(t *testing.T) {
//	parseFile("../example/assignment/conv_struct.go")
//}
//
//// TestParseAssignment_anonymous_function
//func TestParseAssignment_anonymous_function(t *testing.T) {
//	parseFile("../example/assignment/anonymous_function.go")
//}
//
//// TestParseAssignment_built_in_function
//func TestParseAssignment_built_in_function(t *testing.T) {
//	parseFile("../example/assignment/built_in_function.go")
//}
//
//// TestParseAssignment_closure_function
//func TestParseAssignment_closure_function(t *testing.T) {
//	parseFile("../example/assignment/closure_function.go")
//}
//
//// TestParseAssignment_const
//func TestParseAssignment_const(t *testing.T) {
//	parseFile("../example/assignment/const.go")
//}
//
//// TestParseAssignment_function
//func TestParseAssignment_function(t *testing.T) {
//	parseFile("../example/assignment/function.go")
//}
//
//// TestParseAssignment_inner_function
//func TestParseAssignment_inner_function(t *testing.T) {
//	parseFile("../example/assignment/inner_function.go")
//}
//
//// TestParseAssignment_new
//func TestParseAssignment_new(t *testing.T) {
//	parseFile("../example/assignment/new.go")
//}
//
//// TestParseAssignment_receiver_function
//func TestParseAssignment_receiver_function(t *testing.T) {
//	parseFile("../example/assignment/receiver_function.go")
//}
//
//// TestParseAssignment_var
//func TestParseAssignment_var(t *testing.T) {
//	parseFile("../example/assignment/var.go")
//}
//
//// 解析文件
//func parseFile(path string) {
//	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
//		// Process error
//		if err != nil {
//			return err
//		}
//
//		// Only process go files
//		if !info.IsDir() && filepath.Ext(path) != ".go" {
//			return nil
//		}
//
//		// Everything is fine here, extract if path is a file
//		if !info.IsDir() {
//			hasSuffix := strings.HasSuffix(path, "_test.go")
//			if hasSuffix {
//				return nil
//			}
//
//			// Parse file and create the AST
//			var fset = token.NewFileSet()
//			var f *ast.File
//			if f, err = parser.ParseFile(fset, path, nil, parser.ParseComments); err != nil {
//				return nil
//			}
//
//			// 组装所有方法
//			for _, cg := range f.Decls {
//				decl, ok := cg.(*ast.FuncDecl)
//				if ok {
//					AssignmentWalk := AssignmentWalk{}
//					ast.Walk(&AssignmentWalk, decl)
//				}
//				genDecl, ok := cg.(*ast.GenDecl)
//				if ok {
//					si := SourceInfo{}
//					assignment := si.ParseAssignment(genDecl)
//					if assignment != nil {
//						marshal, err := json.Marshal(assignment)
//						if err != nil {
//							fmt.Print("AssignmentWalk walk err: " + err.Error() + "\n")
//						} else {
//							fmt.Print("AssignmentWalk walk: ", string(marshal)+"\n")
//						}
//					}
//				}
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		return
//	}
//}
//
//type AssignmentWalk struct {
//}
//
//func (v *AssignmentWalk) Visit(n ast.Node) ast.Visitor {
//	if n == nil {
//		return v
//	}
//	var assignment []*Assignment
//	switch node := n.(type) {
//	case *ast.AssignStmt:
//		si := SourceInfo{}
//		assignment = si.ParseAssignment(node)
//	case *ast.DeclStmt:
//		si := SourceInfo{}
//		assignment = si.ParseAssignment(node)
//	case *ast.GenDecl:
//		si := SourceInfo{}
//		assignment = si.ParseAssignment(node)
//	// 这种是没有响应值的function
//	case *ast.ExprStmt:
//		si := SourceInfo{}
//		assignment = si.ParseAssignment(node)
//	}
//	if assignment == nil {
//		return v
//	}
//	marshal, err := json.Marshal(assignment)
//	if err != nil {
//		fmt.Print("AssignmentWalk walk err: " + err.Error() + "\n")
//	} else {
//		fmt.Print("AssignmentWalk walk: ", string(marshal)+"\n")
//	}
//
//	return v
//}
