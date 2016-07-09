package interpreter

import (
	"github.com/mlmhl/compiler/parser"
	"github.com/mlmhl/compiler/gdync/interpreter/ast"
	"github.com/mlmhl/compiler/gdync/interpreter/clog"
)

//
// Using  Recursive Descent Parsing and Operator-Precedence Parsing to build the AST
//

type Interpreter struct {
	env    *ast.Environment // global scope context
	parser *parser.Parser

	statements []ast.Statement

	logger *clog.Logger
}

func NewInterpreter() *Interpreter {
	logger, err := clog.NewLogger("stdout")
	if err != nil {
		logger = nil
	}
	return &Interpreter{
		env:    ast.NewEnvironment(ast.VariableSet{}, ast.FunctionSet{}, true),
		parser: parser.NewParser(),

		statements: []ast.Statement{},

		logger: logger,
	}
}

func (interpreter *Interpreter) Interpret(file string) {
	if err := interpreter.parser.Parse(file); err != nil {
		interpreter.logger.InternalError(err)
	}

	interpreter.initNativeFunctions()
	interpreter.create()
	interpreter.execute()
}

func (interpreter *Interpreter) initNativeFunctions() {
	for _, function := range ast.GetNativeFunctions() {
		interpreter.env.AddFunction(function)
	}
}

func (interpreter *Interpreter) create() {
	interpreter.compileUnit()
}

func (interpreter *Interpreter) execute() {
	for _, statement := range interpreter.statements {
		if _, err := statement.Execute(interpreter.env); err != nil {
			interpreter.logger.RuntimeError(err)
		}
	}
}