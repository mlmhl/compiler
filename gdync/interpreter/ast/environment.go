package ast

import (
	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

type FunctionSet map[string]Function
type VariableSet map[string]*types.Variable

// Context for each process point.
// Global variables in global statement will be put into localVariables.
type Environment struct {
	localVariables  VariableSet
	globalVariables VariableSet

	functions FunctionSet
}

func NewEnvironment(globals VariableSet,
	functions FunctionSet, isGlobal bool) *Environment {
	var localVariables VariableSet
	if isGlobal {
		localVariables = nil
	} else {
		localVariables = VariableSet{}
	}

	return &Environment{
		localVariables:  localVariables,
		globalVariables: globals,

		functions: functions,
	}
}

func (env *Environment) IsGlobal() bool {
	return env.localVariables == nil
}

func (env *Environment) GetGlobalVariables() VariableSet {
	return env.globalVariables
}

func (env *Environment) GetFunctions() FunctionSet {
	return env.functions
}

func (env *Environment) GetGlobalVariable(id *types.Identifier) *types.Variable {
	return getVariable(env.globalVariables, id)
}

func (env *Environment) GetLocalVariable(id *types.Identifier) *types.Variable {
	if env.IsGlobal() {
		panic("Can't get local variable in global scope!")
	}
	return getVariable(env.localVariables, id)
}

func (env *Environment) AddGlobalVariable(variable *types.Variable) {
	if !env.IsGlobal() {
		panic("Can't add global variable in local scope!")
	}
	env.globalVariables[variable.GetName()] = variable
}

func (env *Environment) AddLocalVariable(variable *types.Variable) {
	if env.IsGlobal() {
		panic("Can't add local variable in global scope!")
	}
	env.localVariables[variable.GetName()] = variable
}

func (env *Environment) GetFunction(id *types.Identifier) Function {
	if env.functions == nil {
		// No function exist
		return nil
	}
	if function, ok := env.functions[id.GetName()]; ok {
		return function
	} else {
		return nil
	}
}

// The caller must check if the function already exists.
func (env *Environment) AddFunction(function Function) {
	if !env.IsGlobal() {
		panic("Can't add function in local scope!")
	}

	if env.functions == nil {
		env.functions = map[string]Function{}
	}

	env.functions[function.GetName()] = function
}

func getVariable(variables map[string]*types.Variable, id *types.Identifier) *types.Variable {
	if v, ok := variables[id.GetName()]; ok {
		return v
	} else {
		return nil
	}
}
