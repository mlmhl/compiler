package compiler

import (
	"github.com/mlmhl/compiler/gstac/compiler/ast"
	"github.com/mlmhl/compiler/gstac/parser"
	"github.com/mlmhl/compiler/gdync/interpreter/clog"
)

type Compiler struct {
	parser *parser.Parser
	globalContext *ast.Context
	statements []ast.Statement

	logger *clog.Logger
}

func NewCompiler() *Compiler {
	return &Compiler{
		parser: parser.NewParser(),
		globalContext: ast.NewContext(nil, nil),
		statements: []ast.Statement{},
	}
}

func (compiler *Compiler) Compile(fileName string) {
	err := compiler.parser.Parse(fileName)
	if err != nil {
		compiler.logger.InternalError(err)
	}
	compiler.create()
	compiler.fix()
	compiler.generate()
}

// Syntax Analysis
func (compiler *Compiler) create() {

}

// Semantic Analysis
func (compiler *Compiler) fix() {

}

func (compiler *Compiler) generate() {

}