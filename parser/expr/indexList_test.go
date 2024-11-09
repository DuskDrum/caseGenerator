package expr

// 未得到 indexList
//func TestIndexListCase(t *testing.T) {
//	// 定义包含索引列表表达式的 Go 源代码
//	src := `package main
//    func main() {
//        a := [3][3]int{{1,2,3},{4,5,6},{7,8,9}}
//        b := a[1][2]
//    }`
//
//	// 创建文件集
//	fset := token.NewFileSet()
//
//	// 解析源文件
//	file, err := parser.ParseFile(fset, "", src, 0)
//	if err != nil {
//		fmt.Println("解析错误:", err)
//		return
//	}
//
//	// 遍历 AST
//	ast.Inspect(file, func(n ast.Node) bool {
//		// 检查是否为 *ast.IndexListExpr 节点
//		if indexListExpr, ok := n.(*ast.IndexListExpr); ok {
//			fmt.Println("找到索引列表表达式:")
//
//			// 输出被索引的表达式
//			fmt.Printf("  被索引的表达式: %T\n", indexListExpr.X)
//
//			// 输出所有索引值
//			for i, index := range indexListExpr.Indices {
//				fmt.Printf("  索引[%d]: %T\n", i, index)
//			}
//		}
//		return true
//	})
//}
