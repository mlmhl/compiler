package ast

import (
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
	symbolList  *SymbolList
	variables   *container.Trie
	functions   *container.Trie

	outFunction *Function // out function definition
	outContext  *Context  // out level context
}

func NewContext(symbolList *SymbolList, outContext *Context, function *Function) *Context {
	if symbolList == nil {
		symbolList = newSymbolList()
	}
	return &Context{
		symbolList: symbolList,
		variables:  nil,
		functions:  nil,

		outFunction: function,
		outContext: outContext,
	}
}

func (context *Context) GetVariable(name string) *Declaration {
	if context.variables == nil {
		return nil
	}
	if context.variables.Contains(name) {
		return context.variables.Get(name)
	} else {
		if context.outContext != nil {
			return context.outContext.GetVariable(name)
		} else {
			return nil
		}
	}
}

func (context *Context) GetFunction(name string) *Function {
	if context.functions == nil {
		return nil
	}
	if context.functions.Contains(name) {
		return context.functions.Get(name)
	} else {
		if context.outContext != nil {
			return context.outContext.GetFunction(name)
		} else {
			return nil
		}
	}
}

func (context *Context) GetSymbolList() *SymbolList {
	return context.symbolList
}

func (context *Context) GetFunctionList() []*Function {
	functions := []*Function{}
	for iterator := context.functions.Iterator(); ; iterator.HasNext() {
		functions = append(functions, iterator.Next().(*Function))
	}
	return functions
}

func (context *Context) GetOutFunctionDefinition() *Function {
	return context.outFunction
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
