package types

import (
	"github.com/mlmhl/compiler/common"
)

type Variable struct {
	name  *Identifier
	value Value
}

func NewVariable(name *Identifier, value Value) *Variable {
	return &Variable{
		name:  name,
		value: value,
	}
}

func (variable *Variable) SetValue(value Value) *Variable {
	variable.value = value
	return variable
}

func (variable *Variable) GetName() string {
	return variable.name.GetName()
}

func (variable *Variable) GetValue() Value {
	return variable.value
}

func (variable *Variable) GetLocation() *common.Location {
	return variable.name.location
}