package interpreter

import (
	"fmt"

	"github.com/mlmhl/compiler/common"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/compiler/gdync/interpreter/ast"
	"github.com/mlmhl/compiler/gdync/interpreter/types"
	"github.com/mlmhl/compiler/gdync/token"
)

func (interpreter *Interpreter) compileUnit() {
	for {
		tok, err := interpreter.parser.Next()
		if err != nil {
			interpreter.logger.CompileError(err)
		}
		if tok.GetType() == token.FINISHED_ID {
			break
		} else {
			interpreter.parser.RollBack(tok)
		}

		interpreter.definitionOrStatement()
	}
}

func (interpreter *Interpreter) definitionOrStatement() {
	tok, err := interpreter.parser.Next()
	if err != nil {
		interpreter.logger.CompileError(err)
	}

	typ := tok.GetType()
	interpreter.parser.RollBack(tok)

	if typ == token.FINISHED_ID {
		return
	} else if typ == token.FUNCTION_DEFINITION_ID {
		interpreter.env.AddFunction(interpreter.functionDefinition())
	} else {
		interpreter.statements = append(interpreter.statements, interpreter.statement())
	}
}

func (interpreter *Interpreter) functionDefinition() ast.Function {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	// Next token's type must be FUNCTION_DEFINITION_ID,
	// without error, just skip it.
	parser.Next()

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.IDENTIFIER_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("%s should followed by function identifier",
				token.GetDescription(token.FUNCTION_DEFINITION_ID)), tok.GetLocation()))
	}
	identifier := types.NewIdentifier(tok.GetValue().(string), tok.GetLocation())

	parameters := interpreter.parameterList()

	block := interpreter.block()

	return ast.NewCustomFunction(identifier, parameters, block)
}

func (interpreter *Interpreter) parameterList() []*ast.Parameter {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("Function parameters should start by %s",
				token.GetDescription(token.LSP_ID)), tok.GetLocation()))
	}

	parameters := []*ast.Parameter{}

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	} else {
		if tok.GetType() == token.RSP_ID {
			// no parameters
			return parameters
		} else {
			parser.RollBack(tok)
		}
	}

	for {
		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}

		// token's type won't be RSP_ID
		if tok.GetType() == token.IDENTIFIER_ID {
			parameters = append(parameters, ast.NewParameter(
				types.NewIdentifier(tok.GetValue().(string), tok.GetLocation())))
		} else {
			logger.CompileError(gerror.NewSyntaxError(
				fmt.Sprintf("%s can't be used as function parameter",
					token.GetDescription(tok.GetType())), tok.GetLocation()))
		}

		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		} else if tok.GetType() == token.RSP_ID {
			// parameters definition finished
			break
		} else if tok.GetType() != token.FINISHED_ID {
			// no RSP_ID found to finish parameters definition
			logger.CompileError(gerror.NewSyntaxError(
				fmt.Sprintf("Function parameter should ended by %s",
					token.GetDescription(tok.GetType())), tok.GetLocation(),
			))
		} else if tok.GetType() != token.COMMA_ID {
			// skip comma
			logger.CompileError(gerror.NewSyntaxError(
				fmt.Sprintf("%s can't be used as function paramter",
					token.GetDescription(tok.GetType())), tok.GetLocation()))
		}
	}

	return parameters
}

func (interpreter *Interpreter) block() *ast.Block {
	var tok *token.Token
	var err gerror.Error

	parser := interpreter.parser
	logger := interpreter.logger

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() != token.LLP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("Block should start with "+
				token.GetDescription(token.LLP_ID)), tok.GetLocation()))
	}

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() == token.RLP_ID {
		return ast.NewBlock([]ast.Statement{})
	} else {
		parser.RollBack(tok)
	}

	block := ast.NewBlock(interpreter.statementList())

	// Next token's type must be RLP_ID, skip it
	parser.Next()

	return block
}

func (interpreter *Interpreter) statementList() []ast.Statement {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	statements := []ast.Statement{}

	for {
		statements = append(statements, interpreter.statement())

		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		typ := tok.GetType()
		parser.RollBack(tok)

		if typ == token.RLP_ID {
			break
		}
	}

	return statements
}

func (interpreter *Interpreter) statement() ast.Statement {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}

	typ := tok.GetType()
	parser.RollBack(tok)

	switch typ {
	case token.GLOBAL_ID:
		return interpreter.globalStatement()
	case token.IF_ID:
		return interpreter.ifStatement()
	case token.WHILE_ID:
		return interpreter.whileStatement()
	case token.FOR_ID:
		return interpreter.forStatement()
	case token.RETURN_ID:
		return interpreter.returnStatement()
	case token.BREAK_ID:
		return interpreter.breakStatement()
	case token.CONTINUE_ID:
		return interpreter.continueStatement()
	default:
		return interpreter.expressionStatement()
	}
}

func (interpreter *Interpreter) globalStatement() *ast.GlobalStatement {
	parser := interpreter.parser

	// Next token's type must be GLOBAL_ID, skip it
	tok, _ := parser.Next()

	statement := ast.NewGlobalStatement(tok.GetLocation())

	tok, err := parser.Next()
	if err != nil {
		interpreter.logger.CompileError(err)
	}
	if tok.GetType() == token.SEMICOLON_ID {
		return statement
	}

	statement.SetIdentifiers(interpreter.identifierList())

	// If next token's type is SEMICOLON_ID, skip it
	if tok, err = parser.Next(); err != nil {
		interpreter.logger.CompileError(err)
	} else {
		if tok.GetType() != token.SEMICOLON_ID {}
		interpreter.parser.RollBack(tok)
	}

	return statement
}

func (interpreter *Interpreter) identifierList() []*types.Identifier {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	identifiers := []*types.Identifier{}

	for {
		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.IDENTIFIER_ID {
			identifiers = append(identifiers, types.NewIdentifier(
				tok.GetValue().(string), tok.GetLocation()))
		} else {
			logger.CompileError(gerror.NewSyntaxError(
				fmt.Sprintf("In global statement, should be %s, not %s",
					token.GetDescription(token.IDENTIFIER_ID),
					token.GetDescription(tok.GetType())), tok.GetLocation()))
		}

		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.COMMA_ID {
			parser.RollBack(tok)
			logger.CompileError(gerror.NewSyntaxError(
				"Invalid syntax in global statement", tok.GetLocation()))
		}
	}

	return identifiers
}

func (interpreter *Interpreter) ifStatement() *ast.IfStatement {
	parser := interpreter.parser

	// Next token's type must be IF_ID
	tok, _ := parser.Next()
	statement := ast.NewIfStatement(tok.GetLocation())

	statement.SetCondition(interpreter.conditionExpression())
	statement.SetIfBlock(interpreter.block())

	interpreter.elifStatement(statement)

	statement.SetElseBlock(interpreter.elseStatement())

	return statement
}

func (interpreter *Interpreter) elifStatement(statement *ast.IfStatement) {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	for {
		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.ELIF_ID {
			parser.RollBack(tok)
			break
		}
		condition := interpreter.conditionExpression()
		block := interpreter.block()
		statement.AddElifBlock(condition, block, tok.GetLocation())
	}
}

func (interpreter *Interpreter) elseStatement() (*ast.Block, *common.Location) {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() == token.ELSE_ID {
		return interpreter.block(), tok.GetLocation()
	} else {
		return nil, nil
	}
}

func (interpreter *Interpreter) whileStatement() *ast.WhileStatement {
	//Next token's type must be WHILE_ID
	tok, _ := interpreter.parser.Next()
	condition := interpreter.conditionExpression()
	block := interpreter.block()
	return ast.NewWhileStatement(tok.GetLocation(), condition, block)
}

func (interpreter *Interpreter) conditionExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("Condition expression should start with %s, not %s",
				token.GetDescription(token.LSP_ID), token.GetDescription(tok.GetType())),
			tok.GetLocation()))
	}

	expression := interpreter.expression()

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.RSP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("Condition expression should stopped with %s, not %s",
				token.GetDescription(token.RSP_ID), token.GetDescription(tok.GetType())),
			tok.GetLocation()))
	}

	return expression
}

func (interpreter *Interpreter) forStatement() *ast.ForStatement {
	// Next token's type must be FOR_ID
	tok, _ := interpreter.parser.Next()
	statement := ast.NewForStatement(tok.GetLocation())

	interpreter.forExpression(statement)

	statement.SetBlock(interpreter.block())

	return statement
}

func (interpreter *Interpreter) forExpression(statement *ast.ForStatement) {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.LSP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("For expression  must started by %s, bot %s",
				token.GetDescription(token.LSP_ID), token.GetDescription(tok.GetType())),
			tok.GetLocation()))
	}

	expressions := []ast.Expression{}

	for {
		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}

		if tok.GetType() == token.RSP_ID {
			if len(expressions) != 3 {
				logger.CompileError(gerror.NewSyntaxError(
					fmt.Sprintf("For expression need init, condition and post, " +
					"only found %d expressions", len(expressions)), tok.GetLocation()))
			}
			break
		} else {
			if len(expressions) >= 3 {
				logger.CompileError(gerror.NewSyntaxError(
					fmt.Sprintf("For expression only need init, condition and post, " +
					"but found %d expressions", len(expressions) + 1), tok.GetLocation()))
			}
			parser.RollBack(tok)
			expressions = append(expressions, interpreter.expression())
		}
	}

	statement.SetInit(expressions[0])
	statement.SetCondition(expressions[1])
	statement.SetPost(expressions[2])
}

func (interpreter *Interpreter) returnStatement() *ast.ReturnStatement {
	// Next token's type must be RETURN_ID
	tok, _ := interpreter.parser.Next()
	return ast.NewReturnStatement(interpreter.expression(), tok.GetLocation())
}

func (interpreter *Interpreter) breakStatement() *ast.BreakStatement {
	return ast.NewBreakStatement(interpreter.jumpStatement())
}

func (interpreter *Interpreter) continueStatement() *ast.ContinueStatement {
	return ast.NewContinueStatement(interpreter.jumpStatement())
}

func (interpreter *Interpreter) jumpStatement() *common.Location {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	var location *common.Location

	// Next token's type must be BREAK_ID
	tok, _ = parser.Next()
	location = tok.GetLocation()

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.SEMICOLON_ID {
		// skip the semicolon if exist
		parser.RollBack(tok)
	}

	return location
 }

func (interpreter *Interpreter) expressionStatement() *ast.ExpressionStatement {
	return ast.NewExpressionStatement(interpreter.expression())
}
