package ast

import (
	"github.com/mlmhl/compiler/common"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

type Statement interface {
	Execute(env *Environment) (*StatementResult, gerror.Error)
}

//
// Expression statement
//

type ExpressionStatement struct {
	expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{
		expression: expression,
	}
}

func (statement *ExpressionStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	value, err := statement.expression.Evaluate(env)
	if err != nil {
		return nil, err
	}
	return NewStatementResult(NORMAL_STATEMENT_RESULT, value), nil
}

//
// Global statement
//

type GlobalStatement struct {
	identifiers []*types.Identifier
	location    *common.Location // location of keyword 'global'
}

func NewGlobalStatement(location *common.Location) *GlobalStatement {
	return &GlobalStatement{
		identifiers: nil,
		location:    location,
	}
}

func (statement *GlobalStatement) SetIdentifiers(identifiers []*types.Identifier) {
	statement.identifiers = identifiers
}

func (statement *GlobalStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	if env.IsGlobal() {
		return nil, gerror.NewGlobalStatementInTopLevelError(statement.location)
	}

	result := NewStatementResult(NORMAL_STATEMENT_RESULT, nil)

	if statement.identifiers == nil {
		return result, nil
	}

	variables := []*types.Variable{}
	for _, id := range statement.identifiers {
		variable := env.GetGlobalVariable(id)
		if variable == nil {
			return nil, gerror.NewVariableNotFoundError(
				id.GetName(), statement.location)
		}

		localVariable := env.GetLocalVariable(id)
		if localVariable != nil {
			return nil, gerror.NewVariableDuplicateDefinitionError(id.GetName(),
				localVariable.GetLocation(), variable.GetLocation())
		}

		variables = append(variables, variable)
	}

	for _, variable := range variables {
		env.AddLocalVariable(variable)
	}

	return result, nil
}

//
// If statement
//

type elifStatement struct {
	block     *Block
	condition Expression

	location *common.Location // location of 'elif' keyword
}

type elseStatement struct {
	block    *Block
	location *common.Location // location of 'else' keyword
}

type IfStatement struct {
	condition Expression
	ifBlock   *Block
	location  *common.Location // location of 'if' keyword

	elifBlocks []*elifStatement
	elseBlock  *elseStatement
}

func NewIfStatement(location *common.Location) *IfStatement {
	return &IfStatement{
		location:   location,
		elifBlocks: []*elifStatement{},
	}
}

func (statement *IfStatement) SetCondition(condition Expression) {
	statement.condition = condition
}

func (statement *IfStatement) SetIfBlock(block *Block) {
	statement.ifBlock = block
}

func (statement *IfStatement) SetElseBlock(block *Block,
	location *common.Location) {
	statement.elseBlock = &elseStatement{
		block:    block,
		location: location,
	}
}

func (statement *IfStatement) AddElifBlock(
	condition Expression, block *Block, location *common.Location) {
	statement.elifBlocks = append(statement.elifBlocks, &elifStatement{
		block:     block,
		condition: condition,

		location: location,
	})
}

func (statement *IfStatement) Execute(env *Environment) (
	*StatementResult, gerror.Error) {
	var goon types.Value
	var err gerror.Error

	goon, err = statement.condition.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if goon.GetType() != types.BOOL_TYPE {
		return nil, gerror.NewNotBoolExpressionError("if",
			statement.location)
	}
	if goon.GetValue().(bool) {
		// matched
		return statement.ifBlock.Execute(env)
	}

	if statement.elifBlocks != nil {
		for _, elifStatement := range statement.elifBlocks {
			goon, err := elifStatement.condition.Evaluate(env)
			if err != nil {
				return nil, err
			}
			if goon.GetType() != types.BOOL_TYPE {
				return nil, gerror.NewNotBoolExpressionError("elif",
					elifStatement.location)
			}
			if goon.GetValue().(bool) {
				// matched
				return elifStatement.block.Execute(env)
			}
		}
	}

	// no matched, execute the else block if exist
	if statement.elseBlock.block != nil {
		return statement.elseBlock.block.Execute(env)
	} else {
		return NewStatementResult(NORMAL_STATEMENT_RESULT, nil), nil
	}
}

type WhileStatement struct {
	block     *Block
	condition Expression

	location *common.Location // location for 'while' keyword
}

func NewWhileStatement(location *common.Location,
	condition Expression, block *Block) *WhileStatement {
	return &WhileStatement{
		block:     block,
		condition: condition,

		location: location,
	}
}

func (statement *WhileStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	var goon types.Value
	var err gerror.Error
	var result *StatementResult
	for {
		goon, err = statement.condition.Evaluate(env)
		if err != nil {
			break
		}
		if goon.GetType() != types.BOOL_TYPE {
			err = gerror.NewNotBoolExpressionError("while", statement.location)
			break
		}
		if !goon.GetValue().(bool) {
			result = NewStatementResult(NORMAL_STATEMENT_RESULT, nil)
			break
		} else {
			result, err = statement.block.Execute(env)
			if err != nil {
				break
			}
			if result.GetType() == RETURN_STATEMENT_RESULT {
				break
			}
			if result.GetType() == BREAK_STATEMENT_RESULT {
				result.SetType(NORMAL_STATEMENT_RESULT)
				break
			}
			// CONTINUE_STATEMENT_RESULT or NORMAL_STATEMENT_RESULT just go on
		}
	}

	return result, err
}

type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *Block

	location *common.Location // location for 'for' keyword
}

func NewForStatement(location *common.Location) *ForStatement {
	return &ForStatement{
		location: location,
	}
}

func (statement *ForStatement) SetInit(init Expression) {
	statement.init = init
}

func (statement *ForStatement) SetCondition(condition Expression) {
	statement.condition = condition
}

func (statement *ForStatement) SetPost(post Expression) {
	statement.post = post
}

func (statement *ForStatement) SetBlock(block *Block) {
	statement.block = block
}

func (statement *ForStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	if statement.init != nil {
		if _, err := statement.init.Evaluate(env); err != nil {
			return nil, err
		}
	}

	var goon types.Value
	var err gerror.Error
	var result *StatementResult

	for {
		if statement.condition != nil {
			goon, err = statement.condition.Evaluate(env)
			if err != nil {
				return nil, err
			}
			if goon.GetType() != types.BOOL_TYPE {
				return nil, gerror.NewNotBoolExpressionError("for", statement.location)
			}
			if !goon.GetValue().(bool) {
				result = NewStatementResult(NORMAL_STATEMENT_RESULT, nil)
				break
			}
		}

		result, err = statement.block.Execute(env)
		if err != nil {
			return nil, err
		}
		if result.GetType() == RETURN_STATEMENT_RESULT {
			break
		}
		if result.GetType() == BREAK_STATEMENT_RESULT {
			result.SetType(NORMAL_STATEMENT_RESULT)
			break
		}

		if statement.post != nil {
			if _, err = statement.post.Evaluate(env); err != nil {
				break
			}
		}
	}

	return result, err
}

type ReturnStatement struct {
	returnValue Expression
	location    *common.Location // location for 'return' keyword
}

func NewReturnStatement(value Expression,
	location *common.Location) *ReturnStatement {
	return &ReturnStatement{
		returnValue: value,
		location:    location,
	}
}

func (statement *ReturnStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	if statement.returnValue == nil {
		// no return value, return a null instead
		return NewStatementResult(RETURN_STATEMENT_RESULT,
			types.NewValue(types.NULL_TYPE, nil)), nil
	} else {
		value, err := statement.returnValue.Evaluate(env)
		if err != nil {
			return nil, err
		} else {
			return NewStatementResult(RETURN_STATEMENT_RESULT, value), nil
		}
	}
}

type BreakStatement struct {
	location *common.Location // location for 'break' keyword
}

func NewBreakStatement(location *common.Location) *BreakStatement {
	return &BreakStatement{
		location: location,
	}
}

func (statement *BreakStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	return NewStatementResult(BREAK_STATEMENT_RESULT, nil), nil
}

type ContinueStatement struct {
	location *common.Location // location for 'continue' keyword
}

func NewContinueStatement(location *common.Location) *ContinueStatement {
	return &ContinueStatement{
		location: location,
	}
}

func (statement *ContinueStatement) Execute(
	env *Environment) (*StatementResult, gerror.Error) {
	return NewStatementResult(CONTINUE_STATEMENT_RESULT, nil), nil
}
