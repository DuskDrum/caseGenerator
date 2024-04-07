package parse

import (
	"caseGenerator/generate"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func Extract(path string, excludedPaths ...string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Process error
		if err != nil {
			return err
		}

		// Skip excluded paths
		for _, p := range excludedPaths {
			if p == path {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Skip vendor and all directories beginning with a .
		if info.IsDir() && (info.Name() == "vendor" || (len(info.Name()) > 1 && info.Name()[0] == '.')) {
			return filepath.SkipDir
		}

		// Only process go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		// Everything is fine here, extract if path is a file
		if !info.IsDir() {
			hasSuffix := strings.HasSuffix(path, "_test.go")
			if hasSuffix {
				return nil
			}
			if err = extractFile(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func extractFile(filename string) (err error) {
	// Parse file and create the AST
	var fset = token.NewFileSet()
	var f *ast.File
	if f, err = parser.ParseFile(fset, filename, nil, parser.ParseComments); err != nil {
		return
	}

	// 1. 组装package、method信息
	_ = Package{PackagePath: "", PackageName: f.Name.Name}
	// 2. 组装import信息
	ipInfo := importParse(f)
	marshal, err := json.Marshal(&ipInfo)
	if err != nil {
		return
	}
	fmt.Printf("import信息:%v", string(marshal))
	// 组装所有方法
	methods := make([]string, 0, 10)
	// Loop in comment groups
OuterLoop:
	for _, cg := range f.Decls {
		decl, ok := cg.(*ast.FuncDecl)
		if !ok {
			fmt.Print("不处理非 ast.FuncDecl 的内容")
			continue
		}
		// 跳过某些方法的设置
		if decl.Doc != nil {
			for _, com := range decl.Doc.List {
				// 跳过处理
				if strings.Contains(com.Text, "skipcodeautoparse") {
					fmt.Print("skip auto parse...")
					continue OuterLoop
				}
			}
		}

		methods = append(methods, decl.Name.Name)
		// 3. 组装method方法
		methodInfo := Method{MethodName: decl.Name.Name}
		fmt.Print("\n开始处理:" + decl.Name.Name)

		// 4. 组装响应的信息,使用walk，walk特定的类型
		var resultVisit ResultVisitor
		ast.Walk(&resultVisit, cg)
		marshal, err := json.Marshal(&resultVisit)
		if err != nil {
			continue
		}
		fmt.Printf("结果信息:%v \n", string(marshal))
		// 5. 判断所有要改动到的condition
		//var conditionVisitor ConditionVisitor
		//ast.Walk(&conditionVisitor, cg)
		//conMarshal, err := json.Marshal(&conditionVisitor)
		//if err != nil {
		//	continue
		//}
		//fmt.Printf("条件信息:%v \n", string(conMarshal))
		// 6. 定义所有receiver.
		var paramVisitor ParamVisitor
		var rec Receiver
		if decl.Recv != nil {
			var recvName string
			switch recvType := decl.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				switch astStartExpr := recvType.X.(type) {
				case *ast.Ident:
					recvName = astStartExpr.Name
				// 下标类型，实际上是泛型。泛型先不处理
				case *ast.IndexExpr:
					continue
				}
			case *ast.Ident:
				recvName = recvType.Name
			default:
				recvName = decl.Recv.List[0].Names[0].Name
			}
			rec = Receiver{
				ReceiverName: recvName,
			}
			switch typeType := decl.Recv.List[0].Type.(type) {
			case *ast.StarExpr:
				switch astStartExpr := typeType.X.(type) {
				case *ast.Ident:
					rec.ReceiverValue = InvocationUnary{InvocationName: astStartExpr.Name}
				// 下标类型，实际上是泛型。泛型先不处理
				case *ast.IndexExpr:
					continue
				}
			case *ast.Ident:
				rec.ReceiverValue = InvocationUnary{InvocationName: typeType.Name}
			}
			paramVisitor.AddDetail(recvName, &rec)
			paramVisitorMarshal, err := json.Marshal(&paramVisitor)
			if err != nil {
				continue
			}
			fmt.Printf("接受者信息:%v \n", string(paramVisitorMarshal))
		}
		// 7. 定义所有request
		dbs := ParseRequestInfo(decl, ipInfo)
		dbsMarshal, err := json.Marshal(&dbs)
		if err != nil {
			continue
		}
		fmt.Printf("请求参数信息:%v \n", string(dbsMarshal))
		// 8. 获取所有赋值、变量
		ast.Walk(&paramVisitor, cg)
		paramVMarshal, err := json.Marshal(&paramVisitor)
		if err != nil {
			continue
		}
		fmt.Printf("参数信息:%v \n", string(paramVMarshal))
		// 9. 开始处理数据，request信息，名称和类型
		rdList := make([]generate.RequestDetail, 0, 10)
		if len(dbs) > 0 {
			for _, db := range dbs {
				rdList = append(rdList, db)
			}
		}
		importPackageList := make([]string, 0, 10)
		// 10. 开始处理receiver
		uu, _ := uuid.NewUUID()
		cd := generate.CaseDetail{
			CaseName:    uu.String(),
			MethodName:  methodInfo.MethodName,
			RequestList: rdList,
		}
		if decl.Recv != nil {
			cd.ReceiverType = "utils.Empty[" + rec.ReceiverValue.InvocationName + "]()"
			importPackageList = append(importPackageList, "\"slp/reconcile/core/common/utils\"")
		}
		// 11. 开始处理mock

		// 12. 处理import

		packageName := f.Name.Name
		// ./core/service/complex_conditon_service.go，去掉前缀和后缀
		newFilename := strings.Replace(filename, "./", "", 1)
		newFilename = strings.TrimSuffix(newFilename, ".go")
		// 获取最后一个斜杠的位置
		lastIndex := strings.LastIndex(newFilename, "/")
		// 如果找到了斜杠
		if lastIndex == -1 {
			log.Fatalf("路径存在问题")
		}
		// 分割字符串
		filePath := newFilename[:lastIndex+1]
		fileName := newFilename[lastIndex+1:]

		gm := generate.GenMeta{
			Package:        packageName,
			FileName:       fileName,
			FilePath:       filePath,
			ImportPkgPaths: importPackageList,
			CaseDetailList: []generate.CaseDetail{cd},
			MethodName:     methodInfo.MethodName,
		}
		generate.GenFile(gm)
	}

	return
}

func importParse(af *ast.File) Import {
	importList := make([]string, 0, 10)
	importMap := make(map[string]string, 10)
	for _, importSpec := range af.Imports {
		if importSpec.Name == nil {
			importList = append(importList, importSpec.Path.Value)
		} else {
			importMap[importSpec.Name.Name] = importSpec.Path.Value
		}
	}
	return Import{
		ImportList:     importList,
		AliasImportMap: importMap,
	}
}

// containsMethod 判断是否存在于列表
func containsMethod(methodList []string, methodName string) bool {
	if len(methodList) == 0 {
		return false
	}
	for _, m := range methodList {
		if m == methodName {
			return true
		}
	}
	return false
}

// ParseRequestInfo RequestName:  "ctx",
// RequestType:  "context.Context",
// RequestValue: "context.Background()",
func ParseRequestInfo(decl *ast.FuncDecl, ipInfo Import) []generate.RequestDetail {
	dbs := make([]generate.RequestDetail, 0, 10)
	for i, requestParam := range decl.Type.Params.List {
		// "_" 这种不处理了
		var db generate.RequestDetail

		name := requestParam.Names[0].Name
		if name == "_" {
			db.RequestName = "param" + strconv.Itoa(i)
		} else {
			db.RequestName = name
		}
		importPaths := make([]string, 0, 10)

		switch dbType := requestParam.Type.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(dbType)
			db.RequestType = expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
			if expr == "context.Context" {
				db.RequestValue = "context.Background()"
				importPaths = append(importPaths, "\"context\"")
			} else if expr == "time.Time" {
				db.RequestValue = "time.Now()"
				importPaths = append(importPaths, "\"time\"")
			} else {
				db.RequestValue = "utils.Empty[" + expr + "]()"
				importPaths = append(importPaths, "\"slp/reconcile/core/common/utils\"")
			}
		case *ast.Ident:
			db.RequestType = dbType.Name
			if db.RequestType == "string" {
				db.RequestValue = "\"\""
			} else if db.RequestType == "int" {
				db.RequestValue = "0"
			} else if db.RequestType == "int64" {
				db.RequestValue = "0"
			} else if db.RequestType == "any" {
				db.RequestValue = "\"\""
			} else if db.RequestType == "bool" {
				db.RequestValue = "true"
			} else {
				db.RequestValue = "utils.Empty[" + dbType.Name + "]()"
				importPaths = append(importPaths, "\"slp/reconcile/core/common/utils\"")
			}
			// 指针类型
		case *ast.StarExpr:
			switch starXType := dbType.X.(type) {
			case *ast.SelectorExpr:
				expr := GetRelationFromSelectorExpr(starXType)
				db.RequestType = "*" + expr
				if strings.Contains(expr, ".") {
					parts := strings.Split(expr, ".")
					firstField := parts[0]
					importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
				}
				if expr == "context.Context" {
					db.RequestValue = "context.Background()"
					importPaths = append(importPaths, "\"context\"")
				} else {
					db.RequestValue = "lo.ToPtr(" + expr + "{})"
					importPaths = append(importPaths, "\"github.com/samber/lo\"")
				}
			case *ast.Ident:
				db.RequestType = "*" + starXType.Name
				if starXType.Name == "string" {
					db.RequestValue = "lo.ToPtr(\"\")"
				} else {
					db.RequestValue = "lo.ToPtr(" + starXType.Name + "{})"
				}
				importPaths = append(importPaths, "\"github.com/samber/lo\"")
			case *ast.ArrayType:
				requestType, importList := parseArrayType(starXType, ipInfo)
				for _, v := range importList {
					importPaths = append(importPaths, v)
				}
				db.RequestType = "*" + requestType
				// todo: 怎么初始化 *[]string呢
				db.RequestValue = "&" + requestType + "{}"
				importPaths = append(importPaths, "\"github.com/samber/lo\"")

			default:
				log.Fatalf("未知类型...")
			}
		case *ast.FuncType:
			db.RequestType = "func()"
			db.RequestValue = "nil"
		case *ast.InterfaceType:
			// 啥也不做
			db.RequestType = "interface{}"
			db.RequestValue = "nil"
		case *ast.ArrayType:
			requestType, importList := parseArrayType(dbType, ipInfo)
			for _, v := range importList {
				importPaths = append(importPaths, v)
			}
			db.RequestType = requestType
			db.RequestValue = "make(" + requestType + ", 0, 10)"
		case *ast.MapType:
			requestType, portList := parseMapType(dbType, ipInfo)
			for _, v := range portList {
				importPaths = append(importPaths, v)
			}
			db.RequestType = requestType
			db.RequestValue = "make(" + db.RequestType + ", 10)"
		case *ast.Ellipsis:
			continue
		case *ast.ChanType:
			var chanInfo string
			// 处理value
			switch eltType := dbType.Value.(type) {
			case *ast.SelectorExpr:
				expr := GetRelationFromSelectorExpr(eltType)
				chanInfo = expr
				if strings.Contains(expr, ".") {
					parts := strings.Split(expr, ".")
					firstField := parts[0]
					importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
				}
			case *ast.Ident:
				chanInfo = eltType.Name
			case *ast.StarExpr:
				switch starXType := eltType.X.(type) {
				case *ast.SelectorExpr:
					expr := GetRelationFromSelectorExpr(starXType)
					chanInfo = "*" + expr
					if strings.Contains(expr, ".") {
						parts := strings.Split(expr, ".")
						firstField := parts[0]
						importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
					}
				case *ast.Ident:
					chanInfo = "*" + starXType.Name
				default:
					log.Fatalf("未知类型...")
				}
			case *ast.InterfaceType:
				chanInfo = "any"
			default:
				log.Fatalf("未知类型...")
			}
			if dbType.Dir == ast.RECV {
				db.RequestType = "<-chan " + chanInfo
			} else {
				db.RequestType = "chan<- " + chanInfo
			}
			db.RequestValue = "make(" + db.RequestType + ")"
		case *ast.IndexExpr:
			// 下标类型，一般是泛型，处理不了
			continue
		default:
			log.Fatalf("未知类型...")
		}
		db.ImportPkgPath = importPaths
		dbs = append(dbs, db)
	}
	return dbs
}

func parseMapType(mpType *ast.MapType, ipInfo Import) (string, []string) {
	var keyInfo, valueInfo string
	importPaths := make([]string, 0, 10)
	// 处理key
	switch eltType := mpType.Key.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		keyInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		keyInfo = eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			keyInfo = "*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			keyInfo = "*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		keyInfo = "any"
	default:
		log.Fatalf("未知类型...")
	}

	// 处理value
	switch eltType := mpType.Value.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		valueInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		valueInfo = eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			valueInfo = "*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			valueInfo = "*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		valueInfo = "any"
	case *ast.MapType:
		var portList []string
		valueInfo, portList = parseMapType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.ArrayType:
		var portList []string
		valueInfo, portList = parseArrayType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo, importPaths
}

func parseArrayType(dbType *ast.ArrayType, ipInfo Import) (string, []string) {
	var requestType string
	importPaths := make([]string, 0, 10)
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		requestType = "[]" + eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			requestType = "[]*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			requestType = "[]*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:
		var portList []string

		requestType, portList = parseArrayType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.MapType:
		var portList []string
		requestType, portList = parseMapType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return requestType, importPaths
}

type ParamParseResult struct {
	// 参数名
	ParamName string
	// 参数类型
	ParamType string
	// 参数初始化值
	ParamInitValue string
	// 参数校验值
	ParamCheckValue string
	// 参数引入的依赖包
	ImportPkgPath []string
}

// ParseParamList RequestName:  "ctx",
// RequestType:  "context.Context",
// RequestValue: "context.Background()",
func ParseParamList(fieldList []*ast.Field, ipInfo Import) []ParamParseResult {
	dbs := make([]ParamParseResult, 0, 10)
	for i, requestParam := range fieldList {
		// "_" 这种不处理了
		name := requestParam.Names[0].Name
		if name == "_" {
			name = "param" + strconv.Itoa(i)
		}
		db := parseParam(requestParam.Type, name, ipInfo)
		if db == nil {
			continue
		}
		dbs = append(dbs, *db)
	}
	return dbs
}

// parseParam 处理参数
func parseParam(expr ast.Expr, name string, ipInfo Import) *ParamParseResult {
	// "_" 这种不处理了
	var db ParamParseResult

	db.ParamName = name
	importPaths := make([]string, 0, 10)

	switch dbType := expr.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(dbType)
		db.ParamType = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
		if expr == "context.Context" {
			db.ParamInitValue = "context.Background()"
			importPaths = append(importPaths, "\"context\"")
		} else if expr == "time.Time" {
			db.ParamInitValue = "time.Now()"
			importPaths = append(importPaths, "\"time\"")
		} else {
			db.ParamInitValue = "utils.Empty[" + expr + "]()"
			importPaths = append(importPaths, "\"slp/reconcile/core/common/utils\"")
		}
	case *ast.Ident:
		db.ParamType = dbType.Name
		if db.ParamType == "string" {
			db.ParamInitValue = "\"\""
		} else if db.ParamType == "int" {
			db.ParamInitValue = "0"
		} else if db.ParamType == "int64" {
			db.ParamInitValue = "0"
		} else if db.ParamType == "any" {
			db.ParamInitValue = "\"\""
		} else if db.ParamType == "bool" {
			db.ParamInitValue = "true"
		} else {
			db.ParamInitValue = "utils.Empty[" + dbType.Name + "]()"
			importPaths = append(importPaths, "\"slp/reconcile/core/common/utils\"")
		}
		// 指针类型
	case *ast.StarExpr:
		param := parseParam(dbType.X, name, ipInfo)
		if param.ParamType == "context.Context" {
			db.ParamInitValue = "context.Background()"
			importPaths = append(importPaths, "\"context\"")
		} else if param.ParamType == "string" {
			db.ParamInitValue = "lo.ToPtr(\"\")"
			importPaths = append(importPaths, "\"github.com/samber/lo\"")
		} else {
			db.ParamInitValue = "lo.ToPtr(" + param.ParamType + "{})"
			importPaths = append(importPaths, "\"github.com/samber/lo\"")
		}
		db.ParamType = "*" + param.ParamType
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
	case *ast.FuncType:
		db.ParamType = "func()"
		db.ParamInitValue = "nil"
	case *ast.InterfaceType:
		// 啥也不做
		db.ParamType = "interface{}"
		db.ParamInitValue = "nil"
	case *ast.ArrayType:
		requestType, importList := parseParamArrayType(dbType, ipInfo)
		for _, v := range importList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + requestType + ", 0, 10)"
	case *ast.MapType:
		requestType, portList := parseParamMapType(dbType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
		db.ParamType = requestType
		db.ParamInitValue = "make(" + db.ParamType + ", 10)"
	case *ast.Ellipsis:
		return nil
	case *ast.ChanType:
		// 处理value
		param := parseParam(dbType.Value, name, ipInfo)
		if len(db.ImportPkgPath) > 0 {
			for _, v := range db.ImportPkgPath {
				importPaths = append(importPaths, v)
			}
		}
		if dbType.Dir == ast.RECV {
			db.ParamType = "<-chan " + param.ParamType
		} else {
			db.ParamType = "chan<- " + param.ParamType
		}
		db.ParamInitValue = "make(" + param.ParamType + ")"
	case *ast.IndexExpr:
		// 下标类型，一般是泛型，处理不了
		return nil
	default:
		log.Fatalf("未知类型...")
	}
	db.ImportPkgPath = importPaths
	return &db
}

func parseParamMapType(mpType *ast.MapType, ipInfo Import) (string, []string) {
	var keyInfo, valueInfo string
	importPaths := make([]string, 0, 10)
	// 处理key
	switch eltType := mpType.Key.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		keyInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		keyInfo = eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			keyInfo = "*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			keyInfo = "*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		keyInfo = "any"
	default:
		log.Fatalf("未知类型...")
	}

	// 处理value
	switch eltType := mpType.Value.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		valueInfo = expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		valueInfo = eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			valueInfo = "*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			valueInfo = "*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		valueInfo = "any"
	case *ast.MapType:
		var portList []string
		valueInfo, portList = parseParamMapType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.ArrayType:
		var portList []string
		valueInfo, portList = parseParamArrayType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	default:
		log.Fatalf("未知类型...")
	}
	return "map[" + keyInfo + "]" + valueInfo, importPaths
}

func parseParamArrayType(dbType *ast.ArrayType, ipInfo Import) (string, []string) {
	var requestType string
	importPaths := make([]string, 0, 10)
	switch eltType := dbType.Elt.(type) {
	case *ast.SelectorExpr:
		expr := GetRelationFromSelectorExpr(eltType)
		requestType = "[]" + expr
		if strings.Contains(expr, ".") {
			parts := strings.Split(expr, ".")
			firstField := parts[0]
			importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
		}
	case *ast.Ident:
		requestType = "[]" + eltType.Name
	case *ast.StarExpr:
		switch starXType := eltType.X.(type) {
		case *ast.SelectorExpr:
			expr := GetRelationFromSelectorExpr(starXType)
			requestType = "[]*" + expr
			if strings.Contains(expr, ".") {
				parts := strings.Split(expr, ".")
				firstField := parts[0]
				importPaths = append(importPaths, ipInfo.GetImportPath(firstField))
			}
		case *ast.Ident:
			requestType = "[]*" + starXType.Name
		default:
			log.Fatalf("未知类型...")
		}
	case *ast.InterfaceType:
		requestType = "[]any"
	case *ast.ArrayType:
		var portList []string

		requestType, portList = parseParamArrayType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}
	case *ast.MapType:
		var portList []string
		requestType, portList = parseParamMapType(eltType, ipInfo)
		for _, v := range portList {
			importPaths = append(importPaths, v)
		}

	default:
		log.Fatalf("未知类型...")
	}
	return requestType, importPaths
}
