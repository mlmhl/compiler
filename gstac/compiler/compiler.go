package compiler

import (
	"github.com/mlmhl/compiler/gstac/compiler/ast"
	"github.com/mlmhl/compiler/gstac/parser"
	"github.com/mlmhl/compiler/gdync/interpreter/clog"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/executable"
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
		globalContext: ast.NewContext(nil, nil, nil),
		statements: []ast.Statement{},
	}
}

func (compiler *Compiler) Compile(sourceFile string, executableFile string) {
	err := compiler.parser.Parse(sourceFile)
	if err != nil {
		compiler.logger.InternalError(err)
	}
	compiler.create()
	compiler.fix()
	compiler.generate(executableFile)
}

// Syntax Analysis
func (compiler *Compiler) create() {
	compiler.compileUnit()
}

// Semantic Analysis
func (compiler *Compiler) fix() {
	var err errors.Error

	for _, statement := range(compiler.statements) {
		if err = statement.Fix(compiler.globalContext); err != nil {
			compiler.logger.CompileError(err)
		}
	}
	for _, function := range(compiler.globalContext.GetFunctionList()) {
		if err = function.Fix(compiler.globalContext); err != nil {
			compiler.logger.CompileError(err)
		}
	}
}

func (compiler *Compiler) generate(fileName string) {
	executable := executable.NewExecutable()
	executable.Open(fileName)

	// write symbols to file
	executable.AddSymbolList(compiler.globalContext.GetSymbolList().Encode())

	// write functions to file
	executable.BeginFunction()
	for _, function := range(compiler.globalContext.GetFunctionList()) {
		code, err := function.Generate(executable)
		if err != nil {
			compiler.logger.CompileError(err)
		}
		executable.AddFunction(function.GetName(), code)
	}
	executable.EndFunction()

	// write global code to file
	for _, statement := range(compiler.statements) {
		code, err := statement.Generate(executable)
		if err != nil {
			compiler.logger.CompileError(err)
		}
		executable.AddGlobalCode(code)
	}


}