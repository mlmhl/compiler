package compiler

import (
	"fmt"
	"github.com/mlmhl/compiler/gstac/compiler/ast"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/token"
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
	cursor := compiler.parser.GetCursor()

	function, err := compiler.functionDefinition()
	if err == nil {
		// This is a function definition.
		err = compiler.globalContext.AddFunction(function.GetName(), function)
		if err != nil {
			compiler.logger.CompileError(err)
		}
	} else {
		// rollback if failed to pass a `type` token
		compiler.parser.Seek(cursor)

		// This is a statement
		compiler.statements = append(compiler.statements, compiler.statement())
	}

	compiler.parser.Commit()
}

func (compiler *Compiler) functionDefinition() (*ast.Function, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	var err errors.Error
	var tok *token.Token

	var typ ast.Type
	var identifier *ast.Identifier
	var paramList []*ast.Parameter

	// parser return type
	typ, err = compiler.typeSpecifier()
	if err != nil {
		return nil, err
	}

	// parser identifier
	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.IDENTIFIER_ID {
		return nil, errors.NewSyntaxError(
			fmt.Sprintf("%s is invalid function identifier",
				token.GetDescription(tok.GetType())),
			tok.GetLocation())
	}
	identifier = ast.NewIdentifier(tok.GetValue(), tok.GetLocation())

	// parser parameters
	paramList, err = compiler.parameterList()
	if err != nil {
		return nil, err
	}

	// parser block
	statements := compiler.statementListForBlock()

	return ast.NewFunction(typ, identifier, paramList,
		ast.NewFunctionBlock(statements)), nil
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

func (compiler *Compiler) parameterList() ([]*ast.Parameter, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {
		return nil, errors.NewSyntaxError(
			"Function parameter list should start with "+
				token.GetDescription(tok.GetType()),
			tok.GetLocation())
	}

	parameterList := []*ast.Parameter{}

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() == token.RSP_ID {
		// empty parameter list
		return parameterList, nil
	} else {
		parser.RollBack(1)
	}

	for {
		typ, err := compiler.typeSpecifier()
		if err != nil {
			return nil, err
		}

		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.IDENTIFIER_ID {
			return nil, errors.NewSyntaxError(
				fmt.Sprintf("Can't use %s as a identifier",
					token.GetDescription(tok.GetType())),
				tok.GetLocation(),
			)
		}
		parameterList = append(parameterList, ast.NewParameter(typ,
			ast.NewIdentifier(tok.GetValue(), tok.GetLocation())))

		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.RSP_ID {
			break
		} else if tok.GetType() != token.COMMA_ID {
			return nil, errors.NewSyntaxError(
				fmt.Sprintf("Can't use %s in function parameter list",
					token.GetDescription(tok.GetType())),
				tok.GetLocation())
		}
	}

	return parameterList, nil
}

func (compiler *Compiler) statementListForBlock() []ast.Statement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error
	var statements []ast.Statement

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LLP_ID {
		logger.CompileError(errors.NewSyntaxError(
			"Block should start with "+token.GetDescription(tok.GetType()),
			tok.GetLocation()))
	}

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() == token.RLP_ID {
		return nil
	} else {
		parser.RollBack(1)
	}

	statements = compiler.statementList()

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.RLP_ID {
		logger.CompileError(errors.NewSyntaxError(
			"Block should stop with a "+token.GetDescription(tok.GetType()),
			tok.GetLocation()))
	}

	return statements
}

func (compiler *Compiler) statementList() []ast.Statement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	statements := []ast.Statement{}

	for {
		statements = append(statements, compiler.statement())

		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		parser.RollBack(1)

		// statement list is around by large parentheses,
		// so a right large parentheses means statement list ended,
		if tok.GetType() == token.RLP_ID {
			break
		}
	}

	return statements
}

func (compiler *Compiler) statement() ast.Statement {
	parser := compiler.parser
	logger := compiler.logger

	tok, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	parser.RollBack(1)

	switch tok.GetType() {
	case token.IF_ID:
		return compiler.ifStatement()
	case token.FOR_ID:
		return compiler.forStatement()
	case token.WHILE_ID:
		return compiler.whileStatement()
	case token.CONTINUE_ID:
		return compiler.continueStatement()
	case token.BREAK_ID:
		return compiler.breakStatement()
	case token.RETURN_ID:
		return compiler.returnStatement()
	case token.BOOL_TYPE_ID:
		fallthrough
	case token.INTEGER_TYPE_ID:
		fallthrough
	case token.FLOAT_TYPE_ID:
		fallthrough
	case token.STRING_TYPE_ID:
		return compiler.declarationStatement()
	default:
		return compiler.expressionStatement()
	}
}

func (compiler *Compiler) ifStatement() ast.Statement {
	parser := compiler.parser

	var tok *token.Token

	// The first token's type must be IF_ID
	tok, _ = parser.Next()
	ifStatement := ast.NewIfStatement(tok.GetLocation())

	// parse condition expression
	ifStatement.SetCondition(compiler.conditionExpression())

	// parse if's block
	ifStatement.SetIfBlock(ast.NewIfBlock(compiler.statementListForBlock()))

	// parse elif statements, if exists
	ifStatement.SetElifStatements(compiler.elifStatements())

	// parse else statement, if exists
	ifStatement.SetElseBlock(compiler.elseStatement())

	return ifStatement
}

func (compiler *Compiler) elifStatements() ast.ElifStatement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.ELIF_ID {
		parser.RollBack(1)
		return nil
	}

	elifStatement := ast.NewElifStatement(tok.GetLocation())
	elifStatement.SetCondition(compiler.conditionExpression())
	elifStatement.SetBlock(ast.NewIfBlock(compiler.statementListForBlock()))

	return elifStatement
}

func (compiler *Compiler) elseStatement() ast.ElseStatement {
	
}

func (compiler *Compiler) forStatement() ast.Statement {

}

func (compiler *Compiler) whileStatement() ast.Statement {

}

func (compiler *Compiler) continueStatement() ast.Statement {

}

func (compiler *Compiler) breakStatement() ast.Statement {

}

func (compiler *Compiler) returnStatement() ast.Statement {

}

func (compiler *Compiler) declarationStatement() ast.Statement {

}

func (compiler *Compiler) expressionStatement() ast.Statement {

}

func (compiler *Compiler) conditionExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {
		logger.CompileError(errors.NewSyntaxError(
			"Condition expression should start with " +
			token.GetDescription(tok.GetType()),
			tok.GetLocation()))
	}

	expression := compiler.expression()

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.RSP_ID {
		logger.CompileError(errors.NewSyntaxError(
			"Condition expression should end with " +
			token.GetDescription(tok.GetType()),
			tok.GetLocation()))
	}

	return expression
}

func (compiler *Compiler) expression() ast.Expression {

}
