package types

import (
	"github.com/mlmhl/compiler/common"
)

type Identifier struct {
	name string
	location *common.Location
}

func NewIdentifier(name string, location *common.Location) *Identifier {
	return &Identifier{
		name: name,
		location: location,
	}
}

func (identifier *Identifier) GetName() string {
	return identifier.name
}
func (identifier *Identifier) GetLocation() *common.Location {
	return identifier.location
}