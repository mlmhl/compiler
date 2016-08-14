package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
)

type Parameter struct {
	typ        Type
	identifier *Identifier

	// Use type's location as Parameter's location.
	location *common.Location
}

func NewParameter(typ Type, identifier *Identifier, location *common.Location) *Parameter {
	return &Parameter{
		typ: typ,
		identifier: identifier,

		location: location,
	}
}

func (parameter *Parameter) GetType() Type {
	return parameter.typ
}

func (parameter *Parameter) GetIdentifier() *Identifier {
	return parameter.identifier
}

func (parameter *Parameter) GetLocation() *common.Location {
	return parameter.location
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

func (argument *Argument) Fix(context *Context) errors.Error {
	var err errors.Error
	argument.expression, err = argument.expression.Fix(context)
	return err
}

func (argument *Argument) CastTo(destType Type) errors.Error {
	var err errors.Error
	argument.expression, err = argument.expression.CastTo(destType)
	return err
}

type Function struct {
	returnType Type
	identifier *Identifier
	parameters []*Parameter
	block      *FunctionBlock
}

func NewFunction(typ Type, identifier *Identifier,
parameters []*Parameter, block *FunctionBlock) *Function {
	return &Function{
		returnType: typ,
		identifier: identifier,
		parameters: parameters,
		block: block,
	}
}

func (function *Function) GetType() Type {
	return function.returnType
}

func (function *Function) GetName() string {
	return function.identifier.GetName()
}

func (function *Function) GetParameterList() []*Parameter {
	return function.parameters
}

func (function *Function) GetLocation() *common.Location {
	// use identifier's location as function's location
	return function.identifier.GetLocation()
}

func (function *Function) Fix(context *Context) errors.Error {
	localContext := NewContext(context.GetSymbolList(), context, context.GetOutFunctionDefinition())

	// Add parameters as declarations
	for _, param := range(function.GetParameterList()) {
		declaration := localContext.GetVariable(param.GetIdentifier().GetName())
		if declaration != nil {
			return errors.NewParameterDuplicatedDefinitionError(param.GetIdentifier().GetName(),
			declaration.GetLocation(), param.GetIdentifier().GetLocation())
		}
		localContext.AddVariable(param.GetIdentifier().GetName(),
		NewDeclaration(param.GetType(), param.GetIdentifier(), nil, param.GetLocation()))
	}

	return function.block.Fix(context)
}
