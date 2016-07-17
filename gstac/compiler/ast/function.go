package ast

// Use Identifier's location as Parameter's location.
type Parameter struct {
	typ  Type
	identifier *Identifier
}

// Use Expression's location as Argument's location.
type Argument struct {
	expression Expression
}
