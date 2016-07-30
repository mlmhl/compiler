package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/goutil/container"
)

//
// Symbol's definition
//

type SymbolList struct {
	symbols *container.Trie
}

func newSymbolList() *SymbolList {
	return &SymbolList{
		symbols: container.NewTrie(),
	}
}

// get symbol's index
func (symbolList *SymbolList) Get(symbol string) int {
	return symbolList.symbols.Get(symbol)
}

// add a new symbol
func (symbolList *SymbolList) Put(symbol string) {
	symbolList.symbols.Put(symbol, symbolList.symbols.Size())
}

func (symbolList *SymbolList) Contains(symbol string) bool {
	return symbolList.symbols.Contains(symbol)
}

//
// Context
//

type Context struct {
	symbolList *SymbolList
	variables  *container.Trie
	functions  *container.Trie
	outContext *Context
}

func NewContext(symbolList *SymbolList, outContext *Context) *Context {
	if symbolList == nil {
		symbolList = newSymbolList()
	}
	return &Context{
		symbolList: symbolList,
		variables:  nil,
		functions:  nil,
		outContext: outContext,
	}
}

func (context *Context) GetFunctionList() []*Function {
	functions := []*Function{}
	for iterator := context.functions.Iterator(); ; iterator.HasNext() {
		functions = append(functions, iterator.Next().(*Function))
	}
	return functions
}

func (context *Context) AddVariable(name string, declaration *Declaration) errors.Error {
	if context.variables == nil {
		context.variables = container.NewTrie()
	}
	if context.variables.Contains(name) {
		return errors.NewDuplicateDeclarationError(name,
			context.variables.Get(name).(*Declaration).GetLocation(), declaration.GetLocation())
	}
	context.variables.Put(name, declaration)

	context.addSymbolIfNotExist(name)

	return nil
}

func (context *Context) AddFunction(name string, function *Function) errors.Error {
	if context.functions == nil {
		context.functions = container.NewTrie()
	}
	if context.functions.Contains(name) {
		return errors.NewDuplicateFunctionDefinitionError(name,
			context.functions.Get(name).(*Function).GetLocation(), function.GetLocation())
	}
	context.functions.Put(name, function)

	context.addSymbolIfNotExist(name)

	return nil
}

func (context *Context) addSymbolIfNotExist(symbol string) {
	// Put a new symbol if not exist
	if !context.symbolList.Contains(symbol) {
		context.symbolList.Put(symbol)
	}
}
