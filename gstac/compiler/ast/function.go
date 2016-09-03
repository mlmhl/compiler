package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/executable"
	"github.com/mlmhl/goutil/encoding"
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

func (parameter *Parameter) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(encoding.DefaultEncoder.Int(parameter.typ.GetOffset()))
	exe.AppendSlice(encoding.DefaultEncoder.Int(context.GetSymbolIndex(parameter.identifier.GetName())))
	return nil
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

func (argument *Argument) CastTo(destType Type, context *Context) errors.Error {
	var err errors.Error
	argument.expression, err = argument.expression.CastTo(destType, context)
	return err
}

func (argument *Argument) Generate(context *Context, exe *executable.Executable) errors.Error {
	return argument.expression.Generate(context, exe)
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

func (function *Function) Generate(context *Context, exe *executable.Executable) errors.Error {
	// generate return type
	exe.AppendSlice(encoding.DefaultEncoder.Int(function.returnType.GetOffset()))

	// generate function name
	exe.AppendSlice(encoding.DefaultEncoder.Int(context.GetSymbolIndex(function.GetName())))

	// generate parameters' count and each parameter
	exe.AppendSlice(encoding.DefaultEncoder.Int(len(function.parameters)))
	for _, param := range(function.parameters) {
		if err := param.Generate(context, exe); err != nil {
			return err
		}
	}

	// generate function body
	return function.block.Generate(context, exe)
}