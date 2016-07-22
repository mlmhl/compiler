package compiler

import (
	"github.com/mlmhl/compiler/gstac/compiler/ast"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/token"
	"github.com/mlmhl/goutil/container"
)

func (compiler *Compiler) compileUnit() {
	parser := compiler.parser
	logger := compiler.logger

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.FINISHED_ID {
			break
		}
		parser.RollBack(tok)
		compiler.definitionOrStatement()
	}
}

func (compiler *Compiler) definitionOrStatement() {
	function, err := compiler.functionDefinition()
	if err == nil {
		// This is a function definition.
		err = compiler.globalContext.AddFunction(function.GetName(), function)
		if err != nil {
			compiler.logger.CompileError(err)
		}
	} else {
		// This is a statement
		statement, err := compiler.statement()
		if err != nil {
			compiler.logger.CompileError(err)
		} else {
			compiler.statements = append(compiler.statements, statement)
		}
	}
}

func (compiler *Compiler) functionDefinition() (*ast.Function, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	var typ ast.Type
	var err errors.Error
	var tok *token.Token
	var paramList []*ast.Parameter

	typ, err = compiler.typeSpecifier()
	if err != nil {
		return nil, err
	}

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {

	}
}

func (compiler *Compiler) statement() (ast.Statement, errors.Error) {
}

func (compiler *Compiler) basicTypeSpecifier() (ast.Type, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	tok, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}

	switch tok.GetType() {
	case token.BOOL_TYPE_ID:
		return ast.BOOL_TYPE, nil
	case token.INTEGER_TYPE_ID:
		return ast.INTEGER_TYPE, nil
	case token.FLOAT_TYPE_ID:
		return ast.FLOAT_TYPE, nil
	case token.STRING_TYPE_ID:
		return ast.STRING_TYPE, nil
	case token.NULL_ID:
		return ast.NULL_TYPE, nil
	default:
		return nil, errors.NewUnsupportedTypeError(
			token.GetDescription(tok.GetType()), tok.GetLocation())
	}
}

// Up to now, array type specifier is the only composite type.
func (compiler *Compiler) typeSpecifier() (ast.Type, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	var base Type
	var err errors.Error

	// rollback if failed to pass a `type` token
	cursor := parser.GetCursor()
	defer func() {
		if err != nil {
			parser.Seek(cursor)
		}
	}()

	base, err = compiler.basicTypeSpecifier()
	if err != nil {
		return nil, err
	}

	var lTok *token.Token
	var rTok *token.Token
	deriveTags := []ast.DeriveTag{}

	for {
		lTok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if lTok.GetType() != token.LMP_ID {
			parser.RollBack(1)
			break
		}

		rTok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if rTok.GetType() != token.RMP_ID {
			err = errors.NewParenthesesNotMatchedError(token.GetDescription(token.LMP_ID),
			token.GetDescription(token.RMP_ID), lTok.GetLocation(), rTok.GetLocation())
			break
		}

		deriveTags = append(deriveTags, ast.NewArrayderive())
	}

	return ast.NewDeriveType(base, deriveTags), nil
}