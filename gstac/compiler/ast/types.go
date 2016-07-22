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
	return identifier.name
}

func (identifier *Identifier) GetLocation() *common.Location {
	return identifier.location
}

//
// prototype of base types
//
const (
	INTEGER_TYPE = &integerType{
		baseType: baseType{
			name: "int",
		},
	}
	FLOAT_TYPE = &floatType{
		baseType: baseType{
			name: "float",
		},
	}
	STRING_TYPE = &stringType{
		baseType: baseType{
			name: "string",
		},
	}
	BOOL_TYPE = &boolType{
		baseType: baseType{
			name: "bool",
		},
	}
	NULL_TYPE = &nullType{
		baseType: baseType{
			name: "null",
		},
	}
)

type Type interface {
	GetName() string
}

type baseType struct {
	name string
}

func (typ *baseType) GetName() string {
	return typ.name
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

type DeriveTag interface {
	GetTag() string
}

type FunctionDeriveTag struct {
	tag string
	paramList []Parameter
}

func NewFunctionDeriveTag(paramList []Parameter) *FunctionDeriveTag {
	tag := []byte{'('}

	for _, parameter := range paramList {
		tag = append(tag, []byte(parameter.GetType().GetName()))
		tag = append(tag, ',')
	}

	if len(tag) == 0 {
		tag = append(tag, ')')
	} else {
		tag[len(tag)-1] = ')'
	}

	return &FunctionDeriveTag{
		tag: tag,
		paramList: paramList,
	}
}

func (functionDeriveTag *FunctionDeriveTag) GetTag() string {
	return functionDeriveTag.tag
}

type ArrayDerive struct{}

func NewArrayderive() *ArrayDerive {
	return &ArrayDerive{}
}

func (arrayDerive *ArrayDerive) GetTag() string {
	return "[]"
}

type DeriveType struct {
	name string
	base      Type
	deriveTags []DeriveTag
}

func NewDeriveType(base Type, deriveTags []DeriveTag) *DeriveType {
	name := base.GetName();
	for _, tag := range(deriveTags) {
		name = append(name, []byte(tag.GetTag()))
	}
	return &DeriveType{
		name: name,
		base: base,
		deriveTags: deriveTags,
	}
}

func (deriveType *DeriveType) GetName() string {
	return deriveType.name
}

//
// Use identifier's location as Declaration's location.
//
type Declaration struct {
	typ         Type
	identifier  *Identifier
	initializer Expression

	// location for type keyword
	location *common.Location
}

func (declaration *Declaration) GetLocation() *common.Location {
	return declaration.location
}
