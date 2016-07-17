package ast

import (
	"github.com/mlmhl/goutil/container"
)

//
// Symbol's definition
//
type Symbol string

type SymbolList []Symbol

//
// Context
//
type Context struct {
	symbolList SymbolList

	variables *container.Trie
	functions *container.Trie

	outContext *Context
}
