package ast

import "github.com/mlmhl/compiler/common"

//
// identifier
//
type Identifier struct {
	name     string
	location *common.Location
}

func NewIdentifier(name string, location *common.Location) *Identifier {
	return &Identifier{
		name:     name,
		location: location,
	}
}

func (identifier *Identifier) GetName() string {
	return identifier.name;
}

func (identifier *Identifier) GetLocation() *common.Location {
	return identifier.location
}

//
// prototype of base types
//
const (
	INTEGER_TYPE = &integerType{}
	FLOAT_TYPE   = &floatType{}
	STRING_TYPE  = &stringType{}
	BOOL_TYPE    = &boolType{}
	NULL_TYPE    = &nullType{}
)

type Type interface{}

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

//
// Use identifier's location as Declaration's location.
//
type Declaration struct {
	typ         Type
	identifier  *Identifier
	initializer Expression
}
