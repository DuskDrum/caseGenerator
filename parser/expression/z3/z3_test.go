package z3

import (
	"fmt"
	"testing"

	"caseGenerator/go-z3"
)

func TestZ3Case(t *testing.T) {
	// 创建 Z3 上下文
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer ctx.Close()

	// 定义整数变量
	x := ctx.Const(ctx.Symbol("x"), ctx.IntSort())
	y := ctx.Const(ctx.Symbol("y"), ctx.IntSort())

	// 添加约束：x > 10 且 x + y = 20
	solver := ctx.NewSolver()
	solver.Assert(x.Gt(ctx.Int(10, ctx.IntSort())))
	solver.Assert(x.Add(y).Eq(ctx.Int(20, ctx.IntSort())))

	// 求解并输出结果
	if solver.Check() == z3.True {
		model := solver.Model()
		fmt.Println("x =", model.Eval(x).Int())
		fmt.Println("y =", model.Eval(y).Int())
	}
}

// TestZ3Case2 测试复杂的用例， x、y、z
func TestZ3Case2(t *testing.T) {
	// 创建 Z3 上下文
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	// 定义整数变量
	x := ctx.Const(ctx.Symbol("x"), ctx.IntSort())
	y := ctx.Const(ctx.Symbol("y"), ctx.IntSort())
	z := ctx.Const(ctx.Symbol("z"), ctx.IntSort())

	// 添加约束：x > 10 且 x + y = 20 且 x = 15 且 x + z < 10
	solver := ctx.NewSolver()
	solver.Assert(x.Gt(ctx.Int(10, ctx.IntSort())))
	solver.Assert(x.Add(y).Eq(ctx.Int(20, ctx.IntSort())))
	solver.Assert(x.Eq(ctx.Int(15, ctx.IntSort())))
	solver.Assert(x.Add(z).Lt(ctx.Int(10, ctx.IntSort())))

	// 求解并输出结果
	if solver.Check() == z3.True {
		model := solver.Model()
		fmt.Println("x =", model.Eval(x).Int())
		fmt.Println("y =", model.Eval(y).Int())
		fmt.Println("z =", model.Eval(z).Int())
	}
}

// TestZ3Case3 测试复杂的用例， a、b、c、d
func TestZ3Case3(t *testing.T) {
	// 创建 Z3 上下文
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	// 定义整数变量
	a := ctx.Const(ctx.Symbol("a"), ctx.IntSort())
	b := ctx.Const(ctx.Symbol("b"), ctx.IntSort())
	c := ctx.Const(ctx.Symbol("c"), ctx.IntSort())
	d := ctx.Const(ctx.Symbol("d"), ctx.BoolSort())

	// 添加约束：a > 10 且 a + b = 20 且 a = 15 且 a + c < 10 且 !d
	solver := ctx.NewSolver()
	solver.Assert(a.Gt(ctx.Int(10, ctx.IntSort())))
	solver.Assert(a.Add(b).Eq(ctx.Int(20, ctx.IntSort())))
	solver.Assert(a.Eq(ctx.Int(15, ctx.IntSort())))
	solver.Assert(a.Add(c).Lt(ctx.Int(10, ctx.IntSort())))
	solver.Assert(d.Not())

	// 求解并输出结果
	if solver.Check() == z3.True {
		model := solver.Model()
		fmt.Println("a =", model.Eval(a).Int())
		fmt.Println("b =", model.Eval(b).Int())
		fmt.Println("c =", model.Eval(c).Int())
		fmt.Println("d =", model.Eval(d).String())
	}
}

// TestZ3Case4 测试复杂的用例， a、b、c、d
func TestZ3Case4(t *testing.T) {
	// 创建 Z3 上下文
	config := z3.NewConfig()
	ctx := z3.NewContext(config)
	defer func(ctx *z3.Context) {
		err := ctx.Close()
		if err != nil {
			panic(err.Error())
		}
	}(ctx)

	// 定义整数变量
	a := ctx.Const(ctx.Symbol("a"), ctx.FloatSort())

	// 添加约束：a > 10 且 a + b = 20 且 a = 15 且 a + c < 10 且 !d
	solver := ctx.NewSolver()
	as := a.Gt(ctx.Float(10.1, ctx.FloatSort()))

	solver.Assert(as)

	// 求解并输出结果
	if solver.Check() == z3.True {
		model := solver.Model()
		fmt.Println("a =", model.Eval(a).Float())
	}
}
