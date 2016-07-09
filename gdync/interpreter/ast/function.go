package ast

import (
	"fmt"

	"github.com/mlmhl/compiler/common"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

var nativeFunctions []Function = []Function{
	NewPrintfFunction(),
}

func GetNativeFunctions() []Function {
	return nativeFunctions
}

// function signature parameter
type Parameter struct {
	identifier *types.Identifier
}

func NewParameter(identifier *types.Identifier) *Parameter {
	return &Parameter{
		identifier: identifier,
	}
}

// function invoke argument
type Argument struct {
	expression Expression
}

func NewArgument(expression Expression) *Argument {
	return &Argument{
		expression: expression,
	}
}

type Function interface {
	GetName() string

	// Get function definition point's location
	GetLocation() *common.Location

	Evaluate(arguments []types.Value, env *Environment) (types.Value, gerror.Error)
}

//
// native function
//

type PrintfFunction struct {
}

func NewPrintfFunction() *PrintfFunction {
	return &PrintfFunction{}
}

func (f *PrintfFunction) GetName() string {
	return "Printf"
}

func (f *PrintfFunction) GetLocation() *common.Location {
	return nil
}

func (f *PrintfFunction) Evaluate(arguments []types.Value, env *Environment) (types.Value, gerror.Error) {
	if len(arguments) == 0 {
		return nil, gerror.NewArgumentTooFewError(f.GetName(), 1, 0, nil)
	} else {
		format := arguments[0]
		if format.GetType() != types.STRING_TYPE {
			return nil, gerror.NewTypeMismatchError(types.STRING_TYPE.String(),
				format.GetType().String(), format.GetValue(), nil)
		}
		values := []interface{}{}
		for i := 1; i < len(arguments); i++ {
			values = append(values, arguments[i].GetValue())
		}
		fmt.Printf(format.GetValue().(string), values...)
	}

	return nil, nil
}

//
// Custom function
//

type CustomFunction struct {
	block      *Block
	parameters []*Parameter

	// use identifier's location as function's location
	identifier *types.Identifier
}

func NewCustomFunction(identifier *types.Identifier,
	parameters []*Parameter, block *Block) *CustomFunction {
	return &CustomFunction{
		block:      block,
		parameters: parameters,

		identifier: identifier,
	}
}

func (function *CustomFunction) GetName() string {
	return function.identifier.GetName()
}

func (function *CustomFunction) GetLocation() *common.Location {
	return function.identifier.GetLocation()
}

func (function *CustomFunction) Evaluate(arguments []types.Value,
	env *Environment) (types.Value, gerror.Error) {
	if len(arguments) < len(function.parameters) {
		return nil, gerror.NewArgumentTooFewError(
			function.GetName(), len(function.parameters), len(arguments), nil)
	}

	if len(arguments) > len(function.parameters) {
		return nil, gerror.NewArgumentTooManyError(
			function.GetName(), len(function.parameters), len(arguments), nil)
	}

	for i, parameter := range function.parameters {
		argument := arguments[i]
		env.AddLocalVariable(types.NewVariable(parameter.identifier, argument))
	}

	result, err := function.block.Execute(env)
	if err != nil {
		return nil, err
	}

	if result.GetType() == RETURN_STATEMENT_RESULT {
		return result.GetValue(), nil
	} else {
		// no return statement, return a null value
		return types.NewValue(types.NULL_TYPE, nil), nil
	}
}
