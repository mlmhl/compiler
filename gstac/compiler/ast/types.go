package ast

//
// prototype of base types
//
const (
	INTEGER_TYPE = &integerType{}
	FLOAT_TYPE = &floatType{}
	STRING_TYPE = &stringType{}
	BOOL_TYPE = &boolType{}
	NULL_TYPe = &nullType{}
)

type Type interface {}

type baseType struct {
}

type integerType struct {
	baseType
}

type floatType struct {
	baseType
}

type stringType struct {
	baseType
}

type boolType struct {
	baseType
}

type nullType struct {
	baseType
}

//
// derived type
//

type FunctionDerive struct {
	paramList []Parameter
}

type ArrayDerive struct {
}