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

	isPriorityOf(other Type) bool

	isPriorityOfNull() bool
	isPriorityOfBool() bool
	isPriorityOfInteger() bool
	isPriorityOfFloat() bool
	isPriorityOfString() bool
}

type baseType struct {
	name string
}

func (typ *baseType) GetName() string {
	return typ.name
}

func (typ *baseType) isPriorityOfNull() bool {
	return true
}

func (typ *baseType) isPriorityOfString() bool {
	return false
}

type integerType struct {
	baseType
}

func (typ *integerType) isPriorityOf(other Type) bool {
	return other.isPriorityOfInteger()
}

func (typ *integerType) isPriorityOfBool() bool {
	return true
}

func (typ *integerType) isPriorityOfInteger() bool {
	return false
}

func (typ *integerType) ifPriorityOfFloat() bool {
	return false
}

type floatType struct {
	baseType
}

func (typ *floatType) isPriorityOf(other Type) bool {
	return other.isPriorityOfFloat()
}

func (typ *floatType) isPriorityOfBool() bool {
	return true
}

func (typ *floatType) isPriorityOfInteger() bool {
	return true
}

func (typ *floatType) ifPriorityOfFloat() bool {
	return false
}

type stringType struct {
	baseType
}

func (typ *stringType) isPriorityOf(other Type) bool {
	return other.isPriorityOfString()
}

func (typ *stringType) isPriorityOfBool() bool {
	return true
}

func (typ *stringType) isPriorityOfInteger() bool {
	return true
}

func (typ *stringType) ifPriorityOfFloat() bool {
	return true
}

type boolType struct {
	baseType
}

func (typ *boolType) isPriorityOf(other Type) bool {
	return other.isPriorityOfBool()
}

func (typ *boolType) isPriorityOfBool() bool {
	return false
}

func (typ *boolType) isPriorityOfInteger() bool {
	return false
}

func (typ *boolType) ifPriorityOfFloat() bool {
	return false
}

type nullType struct {
	baseType
}

func (typ *nullType) isPriorityOf(other Type) bool {
	return other.isPriorityOfNull()
}

func (typ *nullType) isPriorityOfBool() bool {
	return false
}

func (typ *nullType) isPriorityOfInteger() bool {
	return false
}

func (typ *nullType) ifPriorityOfFloat() bool {
	return false
}

//
// derived type
//

type DeriveTag interface {
	GetTag() string
}

type FunctionDeriveTag struct {
	tag       string
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
		tag:       tag,
		paramList: paramList,
	}
}

func (functionDeriveTag *FunctionDeriveTag) GetTag() string {
	return functionDeriveTag.tag
}

type ArrayDerive struct{}

func NewArrayDerive() *ArrayDerive {
	return &ArrayDerive{}
}

func (arrayDerive *ArrayDerive) GetTag() string {
	return "[]"
}

type DeriveType struct {
	name       string
	base       Type
	deriveTags []DeriveTag
}

func NewDeriveType(base Type, deriveTags []DeriveTag) *DeriveType {
	name := base.GetName()
	for _, tag := range deriveTags {
		name = append(name, []byte(tag.GetTag()))
	}
	return &DeriveType{
		name:       name,
		base:       base,
		deriveTags: deriveTags,
	}
}

func (deriveType *DeriveType) GetName() string {
	return deriveType.name
}

func (deriveType *DeriveType) isPriorityOf(typ Type) {
	panic("Can't invoke `isPriorityOf` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfNull() {
	panic("Can't invoke `isPriorityOfNull` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfBool() {
	panic("Can't invoke `isPriorityOfBool` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfInteger() {
	panic("Can't invoke `isPriorityOfInteger` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfFloat() {
	panic("Can't invoke `isPriorityOfFloat` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfString() {
	panic("Can't invoke `isPriorityOfString` on DeriveType")
}

type Declaration struct {
	typ         Type
	identifier  *Identifier
	initializer Expression

	// Use identifier's location as Declaration's location.
}

func NewDeclaration(typ Type, identifier *Identifier, initializer Expression,
	location *common.Location) *Declaration {
	return &Declaration{
		typ: typ,
		identifier: identifier,
		initializer: initializer,
	}
}

func (declaration *Declaration) GetName() string {
	return declaration.identifier.GetName()
}

func (declaration *Declaration) GetLocation() *common.Location {
	return declaration.identifier.GetLocation()
}

func (declaration *Declaration) Fix(context *Context) {
	declaration.initializer.Fix(context)
	declaration.initializer.CastTo(declaration.typ)
}