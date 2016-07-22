package ast

import "github.com/mlmhl/compiler/common"

// Use Identifier's location as Parameter's location.
type Parameter struct {
	typ        Type
	identifier *Identifier
}

func (parameter *Parameter) GetType() Type {
	return parameter.typ
}

// Use Expression's location as Argument's location.
type Argument struct {
	expression Expression
}

type Function struct {
	// use identifier's location as function's location
	identifier *Identifier
	parameters []Parameter
	block      *Block
}

func (function *Function) GetName() string {
	return function.identifier.GetName()
}

func (function *Function) GetLocation() *common.Location {
	return function.identifier.GetLocation()
}
