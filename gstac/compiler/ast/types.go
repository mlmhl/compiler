package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/goutil/encoding"
	"github.com/mlmhl/compiler/gstac/executable"
	"github.com/mlmhl/compiler/gstac/errors"
)

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
var (
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
	GetBaseType() Type

	// 0: bool
	// 1: float
	// 2: float
	// 3: object(string, array)
	// 4 : null
	GetOffset() int

	Equal(other Type) bool
	IsDeriveType() bool

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

func (typ *baseType) GetBaseType() Type {
	panic("Can't invoke GetBaseType on this type")
}

func (typ *baseType) GetOffset() int {
	panic("Can't invoke GetOffset on this type")
}

func (typ *baseType) Equal(other Type) bool {
	return typ.name == other.GetName()
}

func (typ *baseType) IsDeriveType() bool {
	return false
}

func (typ *baseType) isPriorityOfNull() bool {
	return true
}

func (typ *baseType) isPriorityOfString() bool {
	return false
}

//
// null type
//
type nullType struct {
	baseType
}

func (typ *nullType) GetBaseType() Type {
	return typ
}

func(typ *nullType) GetOffset() int {
	return 4
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

func (typ *nullType) isPriorityOfFloat() bool {
	return false
}

//
// bool type
//
type boolType struct {
	baseType
}

func (typ *boolType) GetBaseType() Type {
	return typ
}

func (typ *boolType) GetOffset() int {
	return 0
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

func (typ *boolType) isPriorityOfFloat() bool {
	return false
}

//
// integer type
//
type integerType struct {
	baseType
}

func (typ *integerType) GetBaseType() Type {
	return typ
}

func (typ *integerType) GetOffset() int {
	return 1
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

func (typ *integerType) isPriorityOfFloat() bool {
	return false
}

type floatType struct {
	baseType
}

func (typ *floatType) GetBaseType() Type {
	return typ
}

func (typ *floatType) GetOffset() int {
	return 2
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

func (typ *floatType) isPriorityOfFloat() bool {
	return false
}

type stringType struct {
	baseType
}

func (typ *stringType) GetBaseType() Type {
	return typ
}

func (typ *stringType) GetOffset() int {
	return 3
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

func (typ *stringType) isPriorityOfFloat() bool {
	return true
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
		tag = append(tag, parameter.GetType().GetName()...)
		tag = append(tag, ',')
	}

	if len(tag) == 0 {
		tag = append(tag, ')')
	} else {
		tag[len(tag)-1] = ')'
	}

	return &FunctionDeriveTag{
		tag:       string(tag),
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
	baseType
	base       Type
	deriveTags []DeriveTag
}

func NewDeriveType(base Type, deriveTags []DeriveTag) *DeriveType {
	name := []byte(base.GetName())
	for _, tag := range deriveTags {
		name = append(name, tag.GetTag()...)
	}
	return &DeriveType{
		baseType: baseType{
			name: string(name),
		},
		base:       base,
		deriveTags: deriveTags,
	}
}

func (typ *DeriveType) GetBaseType() Type {
	return typ.base
}

func (typ *DeriveType) GetOffset() {
	return 3
}

func (typ *DeriveType) IsDeriveType() bool {
	return true
}

func (deriveType *DeriveType) isPriorityOf(typ Type) bool {
	panic("Can't invoke `isPriorityOf` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfNull() bool {
	panic("Can't invoke `isPriorityOfNull` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfBool() bool {
	panic("Can't invoke `isPriorityOfBool` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfInteger() bool {
	panic("Can't invoke `isPriorityOfInteger` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfFloat() bool {
	panic("Can't invoke `isPriorityOfFloat` on DeriveType")
}

func (deriveType *DeriveType) isPriorityOfString() bool {
	panic("Can't invoke `isPriorityOfString` on DeriveType")
}

type Declaration struct {
	typ         Type
	identifier  *Identifier
	initializer Expression

	// Use type's location as Declaration's location.
	location *common.Location
}

func NewDeclaration(typ Type, identifier *Identifier, initializer Expression,
	location *common.Location) *Declaration {
	return &Declaration{
		typ:         typ,
		identifier:  identifier,
		initializer: initializer,

		location: location,
	}
}

func (declaration *Declaration) GetName() string {
	return declaration.identifier.GetName()
}

func (declaration *Declaration) GetType() Type {
	return declaration.typ
}

func (declaration *Declaration) GetLocation() *common.Location {
	return declaration.location
}

func (declaration *Declaration) Fix(context *Context) {
	declaration.initializer.Fix(context)
	declaration.initializer.CastTo(declaration.typ, context)
}

func (declaration *Declaration) Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error) {
	buffer := []byte{}

	// generate location
	buffer = append(buffer, declaration.location.Encode()...)

	// generate variable name's index in symbol list
	buffer = append(buffer, encoding.DefaultEncoder.Int(context.GetSymbolIndex(
		declaration.identifier.GetName())))

	// generate initializer
	expressionCode, err := declaration.initializer.Generate(context)
	if err != nil {
		return nil, err
	}
	buffer = append(buffer, expressionCode...)
	return buffer, nil
}