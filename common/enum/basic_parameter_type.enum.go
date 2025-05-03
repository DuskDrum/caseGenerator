package enum

import (
	"go/ast"
	"go/token"
)

type BasicParameterType BaseEnum

var (
	BASIC_PARAMETER_TYPE_INT     = BasicParameterType{Name: "int", Desc: "int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64"}
	BASIC_PARAMETER_TYPE_FLOAT   = BasicParameterType{Name: "float", Desc: "float32"}
	BASIC_PARAMETER_TYPE_DOUBLE  = BasicParameterType{Name: "double", Desc: "float64"}
	BASIC_PARAMETER_TYPE_BOOL    = BasicParameterType{Name: "bool", Desc: "bool"}
	BASIC_PARAMETER_TYPE_STRING  = BasicParameterType{Name: "string", Desc: "string"}
	BASIC_PARAMETER_TYPE_UNKNOWN = BasicParameterType{Name: "unknown", Desc: "unknown"}
)

var ALL_BASIC_PARAMETER_TYPE = map[string]BasicParameterType{
	BASIC_PARAMETER_TYPE_INT.Name:     BASIC_PARAMETER_TYPE_INT,
	BASIC_PARAMETER_TYPE_FLOAT.Name:   BASIC_PARAMETER_TYPE_FLOAT,
	BASIC_PARAMETER_TYPE_DOUBLE.Name:  BASIC_PARAMETER_TYPE_DOUBLE,
	BASIC_PARAMETER_TYPE_BOOL.Name:    BASIC_PARAMETER_TYPE_BOOL,
	BASIC_PARAMETER_TYPE_STRING.Name:  BASIC_PARAMETER_TYPE_STRING,
	BASIC_PARAMETER_TYPE_UNKNOWN.Name: BASIC_PARAMETER_TYPE_UNKNOWN,
}

func GetParamTypeByBasicLit(lit *ast.BasicLit) BasicParameterType {
	switch lit.Kind {
	case token.INT:
		return BASIC_PARAMETER_TYPE_INT
	case token.FLOAT:
		return BASIC_PARAMETER_TYPE_FLOAT
	case token.STRING:
		return BASIC_PARAMETER_TYPE_STRING
	default:
		return BASIC_PARAMETER_TYPE_UNKNOWN
	}
}
