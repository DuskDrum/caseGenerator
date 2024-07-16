package parser

import (
	"caseGenerator/common/enum"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"go/token"
	"log"
)

type Assignment struct {
	Param
	// 赋值的位置
	ParamIndex int `json:"paramIndex"`
	// 赋值类型
	AssignmentType enum.AssignmentType `json:"assignmentType,omitempty"`
	// 赋值的目标值的值
	StrategyFormulaValue any `json:"strategyFormulaValue,omitempty"`
}

type CompositeLitValue struct {
	// 记录类型
	*Param
	// 记录每个值
	Values []ParamValue
}

func (c CompositeLitValue) ToString() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// CallLitValue x1, x2, x3 := aa.bb(cc,dd,ee)
type CallLitValue struct {
	// 记录方法的类型
	*Param
	// 记录每个方法调用的值和类型
	Values []*ParamValue
}

func (c CallLitValue) ToString() string {
	marshal, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

// AssignmentBinaryNode 组装节点的赋值节点
type AssignmentBinaryNode struct {
	// Param为Y节点的信息
	Param
	Op       token.Token
	Children *AssignmentBinaryNode
}

// AssignmentSlice slice的存储
type AssignmentSlice struct {
	// Param为Y节点的信息
	Param
	LowIndex  string
	HighIndex string
}

// AssignmentUnaryNode 组装节点的赋值节点
type AssignmentUnaryNode struct {
	Param
	// 如果是arrow36，代表了是chan的调用
	Op token.Token
}

func (a AssignmentUnaryNode) ToString(value string) string {
	if a.Op == token.SUB {
		return "-" + value
	} else {
		panic("未知的的Unary操作符")
	}
}

// ParseAssignment 解析赋值语句
// 存在问题： 无法从ast中推测出等式左边的参数是什么类型，只能通过上下文推测
func (s *SourceInfo) ParseAssignment(n ast.Node) []*Assignment {
	if n == nil {
		return nil
	}
	// 赋值列表等于一个 ":=" 或者 "=="，那么列表的strategy是一样的，param有多个
	// 由于有些ast解析不出参数的type，那么param只有名字
	//adi := make([]Assignment, 0, 10)

	switch node := n.(type) {
	case *ast.AssignStmt:
		// 如果Rhs大于1，那么左右的数量是一样的
		if len(node.Rhs) > 1 {
			return s.parseMultipleRhsAssignment(node)
		} else {
			return s.parseSingleRhsAssignment(node)
		}
	case *ast.GenDecl:
		return s.parseGenDeclAssignment(node)
	case *ast.DeclStmt:
		switch nd := node.Decl.(type) {
		case *ast.GenDecl:
			return s.parseGenDeclAssignment(nd)
		default:
			panic("不支持此类型")
		}
		//bo.AppendAssignmentDetailInfoToList(adi)
	// 这种是没有响应值的function
	case *ast.ExprStmt:
		var as Assignment
		as.AssignmentType = enum.ASSIGNMENT_TYPE_EXPRSTMT
		switch nd := node.X.(type) {
		case *ast.CallExpr:
			// 没有响应值的function，没有响应信息
			paramRequests := make([]*ParamValue, 0, 10)
			for _, v := range nd.Args {
				reqInfos := s.ParseParamRequest(v)
				for _, reqInfo := range reqInfos {
					paramRequests = append(paramRequests, reqInfo)
				}
			}
			param := s.ParamParse(nd.Fun)
			as.StrategyFormulaValue = CallLitValue{
				Param:  param,
				Values: paramRequests,
			}
		default:
			panic("不支持此类型")
		}
		return []*Assignment{&as}
	}

	return nil
}

// parseMultipleRhsAssignment 有多个右值的赋值的解析
func (s *SourceInfo) parseMultipleRhsAssignment(node *ast.AssignStmt) []*Assignment {
	if len(node.Lhs) != len(node.Rhs) {
		panic("左右数量不一致")
	}
	assignments := make([]*Assignment, 0, 10)
	// 左边：变量，层级调用
	// 右边：类型断言，赋值，方法
	for i := range node.Lhs {
		as := s.parseAssignmentDetail(node, i, node.Rhs[i])
		assignments = append(assignments, as)
	}
	return assignments
}

// parseSingleRhsAssignment 只有一个右值的赋值的解析
func (s *SourceInfo) parseSingleRhsAssignment(node *ast.AssignStmt) []*Assignment {
	if len(node.Rhs) != 1 {
		panic("符号右边数量不为1")
	}
	assignments := make([]*Assignment, 0, 10)
	// 左边：变量，层级调用
	// 右边：类型断言，赋值，方法
	for i := range node.Lhs {
		as := s.parseAssignmentDetail(node, i, node.Rhs[0])
		assignments = append(assignments, as)
	}
	return assignments
}

func (s *SourceInfo) parseAssignmentDetail(node *ast.AssignStmt, index int, rhsNode ast.Expr) *Assignment {
	// 有一些赋值这一句无法判断出等式右边变量的类型，只能靠变量名上下文联系
	parseResult := s.ParamParse(node.Lhs[index])
	if parseResult.Name == "" && parseResult.Type != "" {
		parseResult.Name = parseResult.Type
		parseResult.Type = ""
	}
	var as Assignment
	as.Param = *parseResult
	as.ParamIndex = index
	// 解析对应下标右边的赋值
	switch nRhsType := rhsNode.(type) {
	case *ast.CallExpr:
		// 没有响应值的function，没有响应信息
		paramRequests := make([]*ParamValue, 0, 10)
		for _, v := range nRhsType.Args {
			reqInfos := s.ParseParamRequest(v)
			for _, reqInfo := range reqInfos {
				paramRequests = append(paramRequests, reqInfo)
			}
		}
		as.AssignmentType = enum.ASSIGNMENT_TYPE_CALL
		param := s.ParamParse(nRhsType.Fun)
		as.StrategyFormulaValue = CallLitValue{
			Param:  param,
			Values: paramRequests,
		}
		// 类型断言可以是 a.(type) 也可以是A.B.C.(type)
	case *ast.TypeAssertExpr:
		// 类型断言已在上面处理了
	// 如果是aa("","") + bb("","")的情况需要处理这个语法树
	case *ast.UnaryExpr:
		result := s.ParamParse(nRhsType.X)
		as.StrategyFormulaValue = AssignmentUnaryNode{Param: *result, Op: nRhsType.Op}
		as.AssignmentType = enum.ASSIGNMENT_TYPE_UNARY
	case *ast.BinaryExpr:
		result := s.parseBinaryAssignment(nRhsType)
		as.StrategyFormulaValue = result
		as.AssignmentType = enum.ASSIGNMENT_TYPE_BINARY
	// 构造类型
	case *ast.CompositeLit:
		// 如果Incomplete为true，那么这个类型是不完整的
		if nRhsType.Incomplete {
			log.Fatal("compositeLit is Incomplete")
		}
		// 构造的type是param， 构造里面的内容共同生成了它的值
		compositeType := s.ParamParse(nRhsType.Type)
		// 解析composite的内容
		values := make([]ParamValue, 0, 10)
		for _, v := range nRhsType.Elts {
			paramValue := s.ParamParseValue(v)
			values = append(values, *paramValue)
		}
		as.StrategyFormulaValue = CompositeLitValue{Param: compositeType, Values: values}
		as.AssignmentType = enum.ASSIGNMENT_TYPE_COMPOSITE
	case *ast.BasicLit:
		as.AssignmentType = enum.ASSIGNMENT_TYPE_BASICLIT
		as.StrategyFormulaValue = nRhsType.Value
	case *ast.FuncLit:
		//requestList, responseList := s.ParseFuncTypeParamParseResult(nRhsType.Type)
		funcType := s.parseFuncType(nRhsType.Type)
		as.StrategyFormulaValue = funcType
		as.AssignmentType = enum.ASSIGNMENT_TYPE_FUNCTION
	case *ast.Ident:
		as.StrategyFormulaValue = nRhsType.Name
		as.AssignmentType = enum.ASSIGNMENT_TYPE_IDENT
	case *ast.StarExpr:
		param := s.ParamParse(nRhsType)
		as.Type = param.Type
		as.AssignmentType = enum.ASSIGNMENT_TYPE_STAR
	case *ast.SelectorExpr:
		param := s.ParamParse(nRhsType)
		as.Type = param.Type
		as.AssignmentType = enum.ASSIGNMENT_TYPE_SELECTOR
	case *ast.SliceExpr:
		// pinter[1:2]这种格式
		as.AssignmentType = enum.ASSIGNMENT_TYPE_SLICE
		// 主体解析
		param := s.ParamParse(nRhsType.X)
		low := s.ParamParseValue(nRhsType.Low)
		high := s.ParamParseValue(nRhsType.High)
		as.StrategyFormulaValue = AssignmentSlice{
			Param: *param,
			LowIndex: lo.TernaryF(low.Value == "", func() string {
				return low.Type
			}, func() string {
				return low.Value
			}),
			HighIndex: lo.TernaryF(high.Value == "", func() string {
				return high.Type
			}, func() string {
				return high.Value
			}),
		}
	default:
		panic("不支持此类型")
	}
	return &as
}

func (s *SourceInfo) parseBinaryAssignment(node *ast.BinaryExpr) *AssignmentBinaryNode {
	var ab AssignmentBinaryNode
	ab.Op = node.Op
	parse := s.ParamParse(node.Y)
	ab.Param = *parse
	// 如果左边是一个二元表达式，那么继续解析
	if xNode, ok := node.X.(*ast.BinaryExpr); ok {
		ab.Children = s.parseBinaryAssignment(xNode)
	} else {
		param := s.ParamParse(node.X)
		ab.Children = &AssignmentBinaryNode{
			Param:    *param,
			Children: nil,
		}
	}
	return &ab
}

func (s *SourceInfo) parseGenDeclAssignment(node *ast.GenDecl) []*Assignment {
	assignments := make([]*Assignment, 0, 10)
	// 目前遇到的例子，nd.Specs结构中只会有一个元素
	for _, ndSpec := range node.Specs {
		switch npVa := ndSpec.(type) {
		case *ast.ValueSpec:
			// 大于1说明name和value要一一对应
			if len(npVa.Values) > 1 {
				if len(npVa.Values) != len(npVa.Names) {
					panic("左右数量不一致")
				}
				for i := range npVa.Names {
					var as Assignment
					as.Param = *s.ParamParse(npVa.Names[i])
					as.ParamIndex = i
					as.AssignmentType = enum.ASSIGNMENT_TYPE_VALUESPEC
					as.StrategyFormulaValue = s.ParamParseValue(npVa.Values[i])
					assignments = append(assignments, &as)
				}
			} else if len(npVa.Values) == 1 {
				for i := range npVa.Names {
					var as Assignment
					as.Param = *s.ParamParse(npVa.Names[i])
					as.ParamIndex = i
					as.AssignmentType = enum.ASSIGNMENT_TYPE_VALUESPEC
					as.StrategyFormulaValue = s.ParamParseValue(npVa.Values[0])
					assignments = append(assignments, &as)
				}
			} else {
				// 这种的type一般有值，没有值就跳过处理了
				if npVa.Type != nil {
					for i := range npVa.Names {
						var as Assignment
						as.Param = *s.ParamParse(npVa.Names[i])
						as.ParamIndex = i
						as.AssignmentType = enum.ASSIGNMENT_TYPE_VALUESPEC
						as.StrategyFormulaValue = s.ParamParseValue(npVa.Type)
						assignments = append(assignments, &as)
					}
				}
			}
		case *ast.TypeSpec:
			fmt.Println("进来了")
		case *ast.ImportSpec:
			fmt.Println("跳过import的解析...")
		}

	}
	return assignments
}
