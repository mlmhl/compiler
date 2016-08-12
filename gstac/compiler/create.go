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

		deriveTags = append(deriveTags, ast.NewArrayDerive())
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

	var statement ast.Statement

	switch tok.GetType() {
	case token.IF_ID:
		statement = compiler.ifStatement()
	case token.FOR_ID:
		statement = compiler.forStatement()
	case token.WHILE_ID:
		statement = compiler.whileStatement()
	case token.CONTINUE_ID:
		statement = compiler.continueStatement()
	case token.BREAK_ID:
		statement = compiler.breakStatement()
	case token.RETURN_ID:
		statement = compiler.returnStatement()
	case token.BOOL_TYPE_ID:
		fallthrough
	case token.INTEGER_TYPE_ID:
		fallthrough
	case token.FLOAT_TYPE_ID:
		fallthrough
	case token.STRING_TYPE_ID:
		statement = compiler.declarationStatement()
	default:
		statement = compiler.expressionStatement()
	}

	// statement maybe ended with a `;`
	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.SEMICOLON_ID {
		parser.RollBack(tok)
	}

	return statement
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

func (compiler *Compiler) elifStatements() []*ast.ElifStatement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	elifStatements := []*ast.ElifStatement{}

	for {
		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.ELIF_ID {
			parser.RollBack(1)
			break
		}

		elifStatement := ast.NewElifStatement(tok.GetLocation())
		elifStatement.SetCondition(compiler.conditionExpression())
		elifStatement.SetBlock(ast.NewIfBlock(compiler.statementListForBlock()))
		elifStatements = append(elifStatements, elifStatement)
	}

	return elifStatements
}

func (compiler *Compiler) elseStatement() *ast.ElseStatement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.ELSE_ID {
		parser.RollBack(1)
		return nil
	}

	return ast.NewElseStatement(
		ast.NewIfBlock(compiler.statementListForBlock()), tok.GetLocation())
}

func (compiler *Compiler) forStatement() ast.Statement {
	parser := compiler.parser
	// next token's type must be FOR_ID
	tok, _ := parser.Next()
	forStatement := ast.NewForStatement(tok.GetLocation())
	compiler.createExpressionForForStatement(forStatement)
	forStatement.SetBlock(ast.NewForBlock(compiler.statementListForBlock()))
	return forStatement
}

func (compiler *Compiler) createExpressionForForStatement(forStatement *ast.ForStatement) {
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
			"For statement's expression should start with "+
				token.GetDescription(tok.GetType()), tok.GetLocation()))
	}

	expressions := []ast.Expression{}

	for {
		expressions = append(expressions, compiler.expression())

		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.RSP_ID {
			break
		} else if tok.GetType() != token.SEMICOLON_ID {
			logger.CompileError(errors.NewSyntaxError(
				fmt.Sprintf("Can't use %s in for statement's expressions",
					token.GetDescription(tok.GetType())), tok.GetLocation()))
		}
	}

	if len(expressions) != 3 {
		logger.CompileError(errors.NewSyntaxError(
			fmt.Sprintf("Wrong for statement's expression size: wanted 3, got %d",
				len(expressions)), tok.GetLocation()))
	}

	forStatement.SetInit(expressions[0])
	forStatement.SetCondition(expressions[1])
	forStatement.SetPost(expressions[1])
}

func (compiler *Compiler) whileStatement() ast.Statement {
	parser := compiler.parser

	// next token's type must be WHILE_ID
	tok, _ := parser.Next()
	return ast.NewWhileStatement(compiler.conditionExpression(),
		ast.NewWhileBlock(compiler.statementListForBlock()), tok.GetLocation())
}

func (compiler *Compiler) continueStatement() ast.Statement {
	// next token's type must be CONTINUE_ID
	tok, _ := compiler.parser.Next()
	return ast.NewContinueStatement(tok.GetLocation())
}

func (compiler *Compiler) breakStatement() ast.Statement {
	// next token's type must be BREAK_ID
	tok, _ := compiler.parser.Next()
	return ast.NewBreakStatement(tok.GetLocation())
}

func (compiler *Compiler) returnStatement() ast.Statement {
	// next token's type must be RETURN_ID
	tok, _ := compiler.parser.Next()
	return ast.NewReturnStatement(compiler.expression(), tok.GetLocation())
}

func (compiler *Compiler) declarationStatement() ast.Statement {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	var typ ast.Type
	var identifier *ast.Identifier

	typ, err = compiler.typeSpecifier()
	if err != nil {
		logger.CompileError(err)
	}

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType != token.IDENTIFIER_ID {
		logger.CompileError(errors.NewSyntaxError(
			fmt.Sprintf("Can't use %s in declaration statement",
				token.GetDescription(tok.GetType())), tok.GetLocation()))
	}
	identifier = ast.NewIdentifier(tok.GetValue(), tok.GetLocation())

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.ASSIGN_ID {
		parser.RollBack(1)
		return ast.NewDeclarationStatement(ast.NewDeclaration(
			typ, identifier, nil, tok.GetLocation()))
	} else {
		return ast.NewDeclarationStatement(ast.NewDeclaration(
			typ, identifier, compiler.expression(), tok.GetLocation()))
	}
}

func (compiler *Compiler) expressionStatement() ast.Statement {
	return ast.NewExpressionStatement(compiler.expression())
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
			"Condition expression should start with "+
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
			"Condition expression should end with "+
				token.GetDescription(tok.GetType()),
			tok.GetLocation()))
	}

	return expression
}

func (compiler *Compiler) expression() ast.Expression {
	parser := compiler.parser

	var result ast.Expression
	var err errors.Error

	cursor := parser.GetCursor()
	result, err = compiler.assignExpression()
	if err == nil {
		return result
	}

	parser.Seek(cursor)
	return compiler.logicalOrExpression()
}

// assign expression
func (compiler *Compiler) assignExpression() (ast.Expression, errors.Error) {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	var left ast.Expression

	left, err = compiler.primaryExpression()
	if err != nil {
		return nil, err
	}

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if !token.IsAssignOperator(tok.GetType()) {
		return nil, errors.NewSyntaxError(fmt.Sprintf("Can't use %s in assign expression",
			token.GetDescription(tok.GetType())), tok.GetLocation())
	}

	return ast.NewAssignExpression(tok.GetType(), left, compiler.expression())
}

func (compiler *Compiler) logicalOrExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.logicalAndExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.OR_ID {
			parser.RollBack(1)
			break
		}
		result = ast.NewLogicalOrExpression(result, compiler.logicalAndExpression())
	}

	return result
}

func (compiler *Compiler) logicalAndExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.equalityExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.AND_ID {
			parser.RollBack(1)
			break
		}
		result = ast.NewLogicalAndExpression(result, compiler.equalityExpression())
	}

	return result
}

func (compiler *Compiler) equalityExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.relationExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.EQUAL_ID {
			result = ast.NewEqualExpression(result, compiler.relationExpression())
		} else if tok.GetType() == token.UNEQUAL_ID {
			result = ast.NewNotEqualExpression(result, compiler.relationExpression())
		} else {
			parser.RollBack(1)
			break
		}
	}

	return result
}

func (compiler *Compiler) relationExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.additiveExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		switch tok.GetType() {
		case token.GT_ID:
			result = ast.NewGreaterThanExpression(result, compiler.additiveExpression())
		case token.GTE_ID:
			result = ast.NewGreaterThanAndEqualExpression(result, compiler.additiveExpression())
		case token.LT_ID:
			result = ast.NewLessThanExpression(result, compiler.additiveExpression())
		case token.LTE_ID:
			result = ast.NewLessThanAndEqualExpression(result, compiler.additiveExpression())
		default:
			break
		}
	}

	return result
}

func (compiler *Compiler) additiveExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.multiplicativeExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.ADD_ID {
			result = ast.NewAddExpression(result, compiler.multiplicativeExpression())
		} else if tok.GetType() == token.SUBTRACT_ID {
			result = ast.NewSubtractExpression(result, compiler.multiplicativeExpression())
		} else {
			break
		}
	}

	return result
}

func (compiler *Compiler) multiplicativeExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	result := compiler.unaryExpression()
	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		switch tok.GetType() {
		case token.MULTIPLY_ID:
			result = ast.NewMultiplyExpression(result, compiler.unaryExpression())
		case token.DIVIDE_ID:
			result = ast.NewDivideExpression(result, compiler.unaryExpression())
		case token.MOD_ID:
			result = ast.NewModExpression(result, compiler.unaryExpression())
		default:
			parser.RollBack(1)
			break
		}
	}

	return result
}

func (compiler *Compiler) unaryExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	tok, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() == token.SUBTRACT_ID {
		return ast.NewMinusExpression(compiler.unaryExpression())
	} else if tok.GetType() == token.NOT_ID {
		return ast.NewLogicalNotExpression(compiler.unaryExpression())
	} else {
		return compiler.primaryExpression()
	}
}

func (compiler *Compiler) primaryExpression() ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	tok, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() == token.NEW_ID {
		return compiler.arrayCreationExpression(tok)
	} else {
		parser.RollBack(1)

		result := compiler.primaryExpressionWithoutArrayCreation(tok)

		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}

		if tok.GetType() == token.LMP_ID {
			// array index expression
			for {
				result = ast.NewIndexExpression(result, compiler.dimensionExpression(tok))

				tok, err = parser.Next()
				if err != nil {
					logger.CompileError(err)
				}
				if tok.GetType() != token.LMP_ID {
					// needn't roll back, go on process array index syntax
					break
				}
			}
		}

		if tok.GetType() == token.INCREMENT_ID {
			result = ast.NewIncrementExpression(result)
		} else if tok.GetType() == token.DECREMENT_ID {
			result = ast.NewDecrementExpression(result)
		} else {
			parser.RollBack(1)
		}

		return result
	}
}

func (compiler *Compiler) primaryExpressionWithoutArrayCreation(first *token.Token) ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	switch first.GetType() {
	case token.LSP_ID:
		return compiler.subExpression(first)
	case token.NULL_ID:
		return ast.NewNullExpression(first.GetLocation())
	case token.TRUE_ID:
		return ast.NewBoolExpression(true, first.GetLocation())
	case token.FALSE_ID:
		return ast.NewBoolExpression(false, first.GetLocation())
	case token.INTEGER_VALUE_ID:
		return ast.NewIntegerExpression(first.GetValue().(int64), first.GetLocation())
	case token.FLOAT_VALUE_ID:
		return ast.NewFloatExpression(first.GetValue().(float64), first.GetLocation())
	case token.STRING_VALUE_ID:
		return ast.NewStringExpression(first.GetValue().(string), first.GetLocation())
	case token.LLP_ID:
		return compiler.arrayLiteralExpression(first)
	}

	if first.GetType() != token.IDENTIFIER_ID {
		logger.CompileError(errors.NewSyntaxError("Can't parse"+
			token.GetDescription(first.GetType()), first.GetLocation()))
	}

	identifier := ast.NewIdentifier(first.GetValue(), first.GetLocation())
	second, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if second.GetType() == token.LSP_ID {
		// function call
		return ast.NewFunctionCallExpression(identifier, compiler.argumentList(second))
	} else {
		parser.RollBack(1)
		return ast.NewIdentifierExpression(identifier)
	}
}

func (compiler *Compiler) subExpression(leftParentheses *token.Token) ast.Expression {
	expression := compiler.expression()

	tok, err := compiler.parser.Next()
	if err != nil {
		compiler.logger.CompileError(err)
	}
	if tok.GetType() != token.RSP_ID {
		compiler.logger.CompileError(errors.NewParenthesesNotMatchedError(
			token.GetDescription(leftParentheses.GetType()),
			token.GetDescription(tok.GetType()),
			leftParentheses.GetLocation(), tok.GetLocation()))
	}

	return expression
}

func (compiler *Compiler) arrayLiteralExpression(
	leftParentheses *token.Token) *ast.ArrayLiteralExpression {
	var tok *token.Token
	var err errors.Error

	values := []ast.Expression{}

	for {
		values = append(values, compiler.expression())

		tok, err = compiler.parser.Next()
		if err != nil {
			compiler.logger.CompileError(err)
		}
		if tok.GetType() == token.RLP_ID {
			break
		} else if tok.GetType() != token.COMMA_ID {
			compiler.logger.CompileError(errors.NewSyntaxError(fmt.Sprintf(
				"Can't use %s in array literal expression", token.GetDescription(tok.GetType())),
				tok.GetLocation()))
		}
	}

	return ast.NewArrayLiteralExpression(values, leftParentheses.GetLocation())
}

func (compiler *Compiler) argumentList(leftParentheses *token.Token) []*ast.Argument {
	var tok *token.Token
	var err errors.Error

	arguments := []*ast.Argument{}

	for {
		arguments = append(arguments, ast.NewArgument(compiler.expression()))

		tok, err = compiler.parser.Next()
		if err != nil {
			compiler.logger.CompileError(err)
		}

		if tok.GetType() == token.LSP_ID {
			break
		} else if tok.GetType() != token.COMMA_ID {
			compiler.logger.CompileError(errors.NewSyntaxError(fmt.Sprintf(
				"Can't use %s in argument list", token.GetDescription(tok.GetType())), tok.GetType()))
		}
	}

	return arguments
}

func (compiler *Compiler) arrayCreationExpression(newTok *token.Token) ast.Expression {
	parser := compiler.parser
	logger := compiler.logger

	var tok *token.Token
	var err errors.Error

	var typ ast.Type

	typ, err = compiler.basicTypeSpecifier()
	if err != nil {
		logger.CompileError(err)
	}

	dimensions := []ast.Expression{}

	for {
		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}

		if tok.GetType() != token.LMP_ID {
			parser.RollBack(1)
		}

		dimensions = append(dimensions, compiler.dimensionExpression(tok))
	}

	if len(dimensions) == 0 {
		logger.CompileError(errors.NewSyntaxError(
			"Can't use `new` on basic type", newTok.GetLocation()))
	}

	return ast.NewArrayCreationExpression(typ, dimensions)
}

func (compiler *Compiler) dimensionExpression(leftParentheses *token.Token) ast.Expression {
	expression := compiler.expression()

	tok, err := compiler.parser.Next()
	if err != nil {
		compiler.logger.CompileError(err)
	}

	if tok.GetType() != token.RMP_ID {
		compiler.logger.CompileError(errors.NewParenthesesNotMatchedError(
			token.GetDescription(leftParentheses.GetType()),
			token.GetDescription(tok.GetType()),
			leftParentheses.GetLocation(), tok.GetLocation()))
	}

	return expression
}