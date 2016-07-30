package ast

import "github.com/mlmhl/compiler/common"

// Use Identifier's location as Parameter's location.
type Parameter struct {
	typ        Type
	identifier *Identifier
}

func NewParameter(typ Type, identifier *Identifier) *Parameter {
	return &Parameter{
		typ: typ,
		identifier: identifier,
	}
}

func (parameter *Parameter) GetType() Type {
	return parameter.typ
}

// Use Expression's location as Argument's location.
type Argument struct {
	expression Expression
}

func NewArgument(expression Expression) *Argument {
	return &Argument{
		expression: expression,
	}
}

type Function struct {
	returnType Type
	identifier *Identifier
	parameters []*Parameter
	block      *FunctionBlock
}

func NewFunction(typ Type, identifier *Identifier,
parameters []Parameter, block *FunctionBlock) {
	return &Function{
		returnType: typ,
		identifier: identifier,
		parameters: parameters,
		block: block,
	}
}

func (function *Function) GetName() string {
	return function.identifier.GetName()
}

func (function *Function) GetLocation() *common.Location {
	// use identifier's location as function's location
	return function.identifier.GetLocation()
}

func (function *Function) Fix(context *Context) {
}
