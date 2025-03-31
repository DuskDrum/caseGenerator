package z3

// #include "go-z3.h"
import "C"

// Sort represents a sort in Z3.
type Sort struct {
	rawCtx  C.Z3_context
	rawSort C.Z3_sort
}

// BoolSort returns the boolean type.
func (c *Context) BoolSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_bool_sort(c.raw),
	}
}

// StringSort returns the string type.
func (c *Context) StringSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_string_sort(c.raw),
	}
}

// IntSort returns the int type.
func (c *Context) IntSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_int_sort(c.raw),
	}
}

// FloatSort returns the int type.
func (c *Context) FloatSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_fpa_sort(c.raw, 8, 24),
	}
}

// DoubleSort returns the int type.
func (c *Context) DoubleSort() *Sort {
	return &Sort{
		rawCtx:  c.raw,
		rawSort: C.Z3_mk_fpa_sort(c.raw, 11, 53),
	}
}

// GetSortKind 获取 sort类型
func GetSortKind(ast *AST) C.Z3_sort_kind {
	sort := C.Z3_get_sort(ast.rawCtx, ast.rawAST)
	kind := C.Z3_get_sort_kind(ast.rawCtx, sort)
	return kind
}
