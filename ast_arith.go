package z3

import (
	"unsafe"
)

// #include "go-z3.h"
import "C"

// Add creates an AST node representing adding.
//
// All AST values must be part of the same context.
func (a *AST) Add(args ...*AST) *AST {
	raws := make([]C.Z3_ast, len(args)+1)
	raws[0] = a.rawAST
	for i, arg := range args {
		raws[i+1] = arg.rawAST
	}

	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_add(
			a.rawCtx,
			C.uint(len(raws)),
			(*C.Z3_ast)(unsafe.Pointer(&raws[0]))),
	}
}

// Mul creates an AST node representing multiplication.
//
// All AST values must be part of the same context.
func (a *AST) Mul(args ...*AST) *AST {
	raws := make([]C.Z3_ast, len(args)+1)
	raws[0] = a.rawAST
	for i, arg := range args {
		raws[i+1] = arg.rawAST
	}

	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_mul(
			a.rawCtx,
			C.uint(len(raws)),
			(*C.Z3_ast)(unsafe.Pointer(&raws[0]))),
	}
}

// Sub creates an AST node representing subtraction.
//
// All AST values must be part of the same context.
func (a *AST) Sub(args ...*AST) *AST {
	raws := make([]C.Z3_ast, len(args)+1)
	raws[0] = a.rawAST
	for i, arg := range args {
		raws[i+1] = arg.rawAST
	}

	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_sub(
			a.rawCtx,
			C.uint(len(raws)),
			(*C.Z3_ast)(unsafe.Pointer(&raws[0]))),
	}
}

// Div
// Z3_mk_div
func (a *AST) Div(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_div(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// Rem   %
// Z3_mk_rem
func (a *AST) Rem(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_rem(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// Lt creates a "less than" comparison.
//
// Maps to: Z3_mk_lt
func (a *AST) Lt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_lt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// Le creates a "less than" comparison.
//
// Maps to: Z3_mk_le
func (a *AST) Le(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_le(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// Gt creates a "greater than" comparison.
//
// Maps to: Z3_mk_gt
func (a *AST) Gt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_gt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// Ge creates a "less than" comparison.
//
// Maps to: Z3_mk_ge
func (a *AST) Ge(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_ge(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// Int相关

// AddInt creates an AST node representing adding.
//
// All AST values must be part of the same context.
func (a *AST) AddInt(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_add(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// MulInt creates an AST node representing multiplication.
//
// All AST values must be part of the same context.
func (a *AST) MulInt(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_mul(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// SubInt creates an AST node representing subtraction.
//
// All AST values must be part of the same context.
func (a *AST) SubInt(t *AST) *AST {

	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_sub(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// DivInt
// Z3_mk_div
func (a *AST) DivInt(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_div(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// RemInt   %
// Z3_mk_rem
func (a *AST) RemInt(t *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_rem(
			a.rawCtx,
			a.rawAST,
			t.rawAST,
		),
	}
}

// LtInt creates a "less than" comparison.
//
// Maps to: Z3_mk_lt
func (a *AST) LtInt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_lt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// LeInt creates a "less than" comparison.
//
// Maps to: Z3_mk_le
func (a *AST) LeInt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_le(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// GtInt creates a "greater than" comparison.
//
// Maps to: Z3_mk_gt
func (a *AST) GtInt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_gt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// GeInt creates a "less than" comparison.
//
// Maps to: Z3_mk_ge
func (a *AST) GeInt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_ge(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// Fpa相关

// FpaAdd creates an AST node representing adding.
//
// All AST values must be part of the same context.
// Z3_mk_fpa_add
func (a *AST) FpaAdd(t1, t2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_add(
			a.rawCtx,
			a.rawAST, // Rounding Mode
			t1.rawAST,
			t2.rawAST,
		),
	}
}

// FpaMul creates an AST node representing multiplication.
//
// All AST values must be part of the same context.
// Z3_mk_fpa_mul
func (a *AST) FpaMul(t1, t2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_mul(
			a.rawCtx,
			a.rawAST, // Rounding Mode
			t1.rawAST,
			t2.rawAST,
		),
	}
}

// FpaSub creates an AST node representing subtraction.
//
// All AST values must be part of the same context.
// Z3_mk_fpa_sub
func (a *AST) FpaSub(t1, t2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_sub(
			a.rawCtx,
			a.rawAST, // Rounding Mode
			t1.rawAST,
			t2.rawAST,
		),
	}
}

// FpaLt creates a "less than" comparison.
//
// Maps to: Z3_mk_fpa_lt
func (a *AST) FpaLt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_lt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// FpaLe creates a "less than" comparison.
//
// Maps to: Z3_mk_fpa_leq
func (a *AST) FpaLe(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_leq(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// FpaGt creates a "greater than" comparison.
//
// Maps to: Z3_mk_fpa_gt
func (a *AST) FpaGt(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_gt(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// FpaGe creates a "greater than" comparison.
//
// Maps to: Z3_mk_fpa_geq
func (a *AST) FpaGe(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_geq(a.rawCtx, a.rawAST, a2.rawAST),
	}
}

// FpaEq creates a "greater than" comparison.
//
// Maps to: Z3_mk_fpa_eq
func (a *AST) FpaEq(a2 *AST) *AST {
	return &AST{
		rawCtx: a.rawCtx,
		rawAST: C.Z3_mk_fpa_eq(a.rawCtx, a.rawAST, a2.rawAST),
	}
}
