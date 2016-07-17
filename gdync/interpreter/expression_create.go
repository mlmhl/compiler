package interpreter

import (
	"fmt"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/compiler/gdync/interpreter/ast"
	"github.com/mlmhl/compiler/gdync/interpreter/types"
	"github.com/mlmhl/compiler/gdync/token"
)

//
// Construction expression based on Operator-Precedence Parsing
//
// Operator precedence as follows:
// -(minus) !
// * / %
// + -
// > >= < <=
// == !=
// &&
// ||
// =
//

func (interpreter *Interpreter) expression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	defer func() {
		// Skip semicolon if exist
		if tok, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.SEMICOLON_ID {
			parser.RollBack(tok)
		}
	}()

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() == token.FINISHED_ID {
		parser.RollBack(tok)
		return nil
	}

	if tok.GetType() == token.IDENTIFIER_ID {
		var nToken *token.Token
		if nToken, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}
		if nToken.GetType() == token.ASSIGN_ID {
			// create a assign expression
			return ast.NewAssignExpression(interpreter.expression(),
				types.NewIdentifier(tok.GetValue().(string), tok.GetLocation()))
		} else {
			parser.RollBack(nToken)
			parser.RollBack(tok)
		}
	} else {
		parser.RollBack(tok)
	}

	return interpreter.logicalOrExpression()

}

func (interpreter *Interpreter) logicalOrExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.logicalAndExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.OR_ID {
			parser.RollBack(tok)
			break
		}
		expression := interpreter.logicalAndExpression()

		// Constant folding
		if _, ok := result.(*ast.BoolExpression); ok {
			if _, ok = expression.(*ast.BoolExpression); ok {
				left, _ := result.Evaluate(nil)
				right, _ := expression.Evaluate(nil)
				value, _ := types.LogicalOperation(types.Or, left, right)
				result = ast.NewBoolExpression(value.GetValue().(bool))
				continue
			}
		}

		result = ast.NewOrExpression(result, expression, tok.GetLocation())
	}

	return result
}

func (interpreter *Interpreter) logicalAndExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.equalityExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.AND_ID {
			parser.RollBack(tok)
			break
		}
		expression := interpreter.equalityExpression()

		// Constant folding
		if _, ok := result.(*ast.BoolExpression); ok {
			if _, ok = expression.(*ast.BoolExpression); ok {
				left, _ := result.Evaluate(nil)
				right, _ := expression.Evaluate(nil)
				value, _ := types.LogicalOperation(types.AND, left, right)
				result = ast.NewBoolExpression(value.GetValue().(bool))
				continue
			}
		}

		result = ast.NewOrExpression(result, expression, tok.GetLocation())
	}

	return result
}

func (interpreter *Interpreter) equalityExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.relationalExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.EQUAL_ID &&
			tok.GetType() != token.UNEQUAL_ID {
			parser.RollBack(tok)
			break
		}

		expression := interpreter.relationalExpression()

		if tok.GetType() == token.EQUAL_ID {
			result = ast.NewEqualExpression(result, expression, tok.GetLocation())
		} else {
			result = ast.NewNotEqualExpression(result, expression, tok.GetLocation())
		}
	}

	return result
}

func (interpreter *Interpreter) relationalExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.additiveExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.GT_ID &&
			tok.GetType() != token.LT_ID &&
			tok.GetType() != token.GTE_ID &&
			tok.GetType() != token.LTE_ID {
			parser.RollBack(tok)
			break
		}

		expression := interpreter.additiveExpression()

		if tok.GetType() == token.GT_ID {
			result = ast.NewGTExpression(result, expression, tok.GetLocation())
		} else if tok.GetType() == token.LT_ID {
			result = ast.NewLTExpression(result, expression, tok.GetLocation())
		} else if tok.GetType() == token.GTE_ID {
			result = ast.NewGTEExpression(result, expression, tok.GetLocation())
		} else {
			result = ast.NewLTEExpression(result, expression, tok.GetLocation())
		}
	}

	return result
}

func (interpreter *Interpreter) additiveExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.multiplicativeExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.ADD_ID &&
			tok.GetType() != token.SUBTRACT_ID {
			parser.RollBack(tok)
			break
		}

		expression := interpreter.multiplicativeExpression()

		if tok.GetType() == token.ADD_ID {
			result = ast.NewAddExpression(result, expression, tok.GetLocation())
		} else {
			result = ast.NewSubtractExpression(result, expression, tok.GetLocation())
		}
	}

	return result
}

func (interpreter *Interpreter) multiplicativeExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	result := interpreter.unaryExpression()

	for {
		tok, err := parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() != token.MULTIPLY_ID &&
			tok.GetType() != token.DIVIDE_ID &&
			tok.GetType() != token.MOD_ID {
			parser.RollBack(tok)
			break
		}

		expression := interpreter.unaryExpression()

		if tok.GetType() == token.MULTIPLY_ID {
			result = ast.NewMultiplyExpression(result, expression, tok.GetLocation())
		} else if tok.GetType() == token.DIVIDE_ID {
			result = ast.NewDivideExpression(result, expression, tok.GetLocation())
		} else {
			result = ast.NewModExpression(result, expression, tok.GetLocation())
		}
	}

	return result
}

func (interpreter *Interpreter) unaryExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error
	var expression ast.Expression

	var result ast.Expression

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() == token.SUBTRACT_ID {
		expression = interpreter.expression()
		result = ast.NewMinusExpression(expression, tok.GetLocation())
	} else if tok.GetType() == token.NOT_ID {
		expression = interpreter.expression()
		result = ast.NewNotExpression(expression, tok.GetLocation())
	} else {
		parser.RollBack(tok)
		result = interpreter.primaryExpression()
	}

	return result
}

func (interpreter *Interpreter) primaryExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	tok, err := parser.Next()
	if err != nil {
		logger.CompileError(err)
	}

	switch tok.GetType() {
	case token.IDENTIFIER_ID:
		var nToken *token.Token
		if nToken, err = parser.Next(); err != nil {
			logger.CompileError(err)
		}

		if nToken.GetType() == token.LSP_ID {
			// function call expression
			parser.RollBack(nToken)
			arguments := interpreter.argumentList()
			return ast.NewFunctionCallExpression(arguments, types.NewIdentifier(
				tok.GetValue().(string), tok.GetLocation()), tok.GetLocation())
		} else {
			parser.RollBack(nToken)
			// identifier expression
			return ast.NewIdentifierExpression(types.NewIdentifier(
				tok.GetValue().(string), tok.GetLocation()))
		}

	case token.LSP_ID:
		parser.RollBack(tok)
		return interpreter.embedExpression()

	case token.INTEGER_ID:
		return ast.NewIntegerExpression(tok.GetValue().(int64))

	case token.FLOAT_ID:
		return ast.NewFloatExpression(tok.GetValue().(float64))

	case token.STRING_ID:
		value, err := ast.NewStringExpression(tok.GetValue().(string))
		if err != nil {
			err.SetLocation(tok.GetLocation())
			logger.CompileError(err)
		}
		return value

	case token.TRUE_ID:
		fallthrough
	case token.FALSE_ID:
		return ast.NewBoolExpression(tok.GetValue().(bool))

	case token.NULL_ID:
		return ast.NewNullExpression()

	default:
		parser.RollBack(tok)
		return nil
	}
}

func (interpreter *Interpreter) argumentList() []*ast.Argument {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error

	// Next token's type must be LSP_ID, skip it
	parser.Next()

	arguments := []*ast.Argument{}

	if tok, err = parser.Next(); err != nil {
		logger.CompileError(err)
	}

	if tok.GetType() == token.RSP_ID {
		// no arguments
		return arguments
	} else {
		parser.RollBack(tok)
	}

	for {
		arguments = append(arguments, ast.NewArgument(interpreter.expression()))

		tok, err = parser.Next()
		if err != nil {
			logger.CompileError(err)
		}
		if tok.GetType() == token.RSP_ID {
			// arguments list finished
			break
		}  else if tok.GetType() != token.COMMA_ID {
			logger.CompileError(gerror.NewSyntaxError(
				fmt.Sprintf("Can't use %s in function arguments list",
					token.GetDescription(tok.GetType())), tok.GetLocation()))
		}
	}

	return arguments
}

func (interpreter *Interpreter) embedExpression() ast.Expression {
	parser := interpreter.parser
	logger := interpreter.logger

	var tok *token.Token
	var err gerror.Error
	var expression ast.Expression

	// Next token's type must be LSP_ID, skip it
	parser.Next()

	expression = interpreter.expression()

	tok, err = parser.Next()
	if err != nil {
		logger.CompileError(err)
	}
	if tok.GetType() != token.RSP_ID {
		logger.CompileError(gerror.NewSyntaxError(
			fmt.Sprintf("Embed expression should ended with %s, not %s",
			token.GetDescription(token.RSP_ID), token.GetDescription(tok.GetType())),
			tok.GetLocation()))
	}

	return expression
}
