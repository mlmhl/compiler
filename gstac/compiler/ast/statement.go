package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/executable"
)

type Statement interface {
	Fix(context *Context) errors.Error
	Generate(executable *executable.Executable) ([]byte, errors.Error)
}

//
// block
//

type UndefinedBlock struct {
}

func NewUndefinedBlock() *UndefinedBlock {
	return &UndefinedBlock{}
}

type baseBlock struct {
	statements []Statement
}

func (block *baseBlock) Fix(context *Context) errors.Error {
	for _, statement := range block.statements {
		if err := statement.Fix(context); err != nil {
			return err
		}
	}
	return nil
}

type IfBlock struct {
	baseBlock
}

func NewIfBlock(statementList []Statement) *IfBlock {
	return &IfBlock{
		baseBlock: baseBlock{
			statements: statementList,
		},
	}
}

type ForBlock struct {
	baseBlock
}

func NewForBlock(statementList []Statement) *ForBlock {
	return &ForBlock{
		baseBlock: baseBlock{
			statements: statementList,
		},
	}
}

type WhileBlock struct {
	baseBlock
}

func NewWhileBlock(statementList []Statement) *WhileBlock {
	return &WhileBlock{
		baseBlock: baseBlock{
			statements: statementList,
		},
	}
}

type FunctionBlock struct {
	baseBlock
}

func NewFunctionBlock(statementList []Statement) *FunctionBlock {
	return &FunctionBlock{
		baseBlock: baseBlock{
			statements: statementList,
		},
	}
}

//
// If statement
//

type ElifStatement struct {
	condition Expression
	block     *IfBlock

	// location of `elif` keyword
	location *common.Location
}

func NewElifStatement(location *common.Location) *ElifStatement {
	return &ElifStatement{
		location: location,
	}
}

func (elifStatement *ElifStatement) SetCondition(condition Expression) {
	elifStatement.condition = condition
}

func (elifStatement *ElifStatement) SetBlock(block *IfBlock) {
	elifStatement.block = block
}

func (elifStatement *ElifStatement) Fix(context *Context) errors.Error {
	var err errors.Error
	elifStatement.condition, err = elifStatement.condition.Fix(context)
	if err != nil {
		return err
	}
	return elifStatement.block.Fix(NewContext(context.GetSymbolList(),
		context, context.GetOutFunctionDefinition()))
}

type ElseStatement struct {
	block *IfBlock

	// location of `else` keyword
	location *common.Location
}

func NewElseStatement(block *IfBlock, location *common.Location) *ElseStatement {
	return &ElseStatement{
		block:    block,
		location: location,
	}
}

func (elseStatement *ElseStatement) Fix(context *Context) errors.Error {
	return elseStatement.block.Fix(NewContext(context.GetSymbolList(),
		context, context.GetOutFunctionDefinition()))
}

type IfStatement struct {
	condition      Expression
	ifBlock        *IfBlock
	elifStatements []*ElifStatement
	elseStatement  *ElseStatement

	// location of `if` keyword
	location *common.Location
}

func NewIfStatement(location *common.Location) *IfStatement {
	return &IfStatement{
		location: location,
	}
}

func (ifStatement *IfStatement) SetCondition(condition Expression) {
	ifStatement.condition = condition
}

func (ifStatement *IfStatement) SetIfBlock(block *IfBlock) {
	ifStatement.ifBlock = block
}

func (ifStatement *IfStatement) SetElifStatements(statements []*ElifStatement) {
	ifStatement.elifStatements = statements
}

func (ifStatement *IfStatement) SetElseBlock(statement *ElseStatement) {
	ifStatement.elseStatement = statement
}

func (ifStatement *IfStatement) Fix(context *Context) errors.Error {
	var err errors.Error

	ifStatement.condition, err = ifStatement.condition.Fix(context)
	if err != nil {
		return err
	}
	if err = ifStatement.ifBlock.Fix(NewContext(context.GetSymbolList(), context,
		context.GetOutFunctionDefinition())); err != nil {
		return err
	}

	for _, elifStatement := range ifStatement.elifStatements {
		if err = elifStatement.Fix(context); err != nil {
			return err
		}
	}

	return ifStatement.elseStatement.Fix(context)
}

//
// for statement
//

type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *ForBlock

	// location of `for` keyword
	location *common.Location
}

func NewForStatement(location *common.Location) *ForStatement {
	return &ForStatement{
		location: location,
	}
}

func (forStatement *ForStatement) SetInit(init Expression) {
	forStatement.init = init
}

func (forStatement *ForStatement) SetCondition(condition Expression) {
	forStatement.condition = condition
}

func (forStatement *ForStatement) SetPost(post Expression) {
	forStatement.post = post
}

func (forStatement *ForStatement) SetBlock(block *ForBlock) {
	forStatement.block = block
}

func (forStatement *ForStatement) Fix(context *Context) errors.Error {
	var err errors.Error

	forStatement.init, err = forStatement.init.Fix(context)
	if err != nil {
		return err
	}

	forStatement.condition, err = forStatement.condition.Fix(context)
	if err != nil {
		return err
	}

	forStatement.post, err = forStatement.post.Fix(context)
	if err != nil {
		return err
	}

	return forStatement.block.Fix(NewContext(context.GetSymbolList(),
		context, context.GetOutFunctionDefinition()))
}

//
// while statement
//

type WhileStatement struct {
	condition Expression
	block     *WhileBlock

	// location of `while` keyword
	location *common.Location
}

func NewWhileStatement(condition Expression, block *WhileBlock,
	location *common.Location) *WhileStatement {
	return &WhileStatement{
		condition: condition,
		block:     block,
		location:  location,
	}
}

func (whileStatement *WhileStatement) Fix(context *Context) errors.Error {
	var err errors.Error

	whileStatement.condition, err = whileStatement.condition.Fix(context)
	if err != nil {
		return err
	}

	return whileStatement.block.Fix(NewContext(context.GetSymbolList(),
		context, context.GetOutFunctionDefinition()))
}

//
// continue statement
//

type ContinueStatement struct {
	location *common.Location
}

func NewContinueStatement(location *common.Location) *ContinueStatement {
	return &ContinueStatement{
		location: location,
	}
}

func (continueStatement *ContinueStatement) Fix(context *Context) errors.Error {
	return nil
}

//
// break statement
//

type BreakStatement struct {
	location *common.Location
}

func NewBreakStatement(location *common.Location) *BreakStatement {
	return &BreakStatement{
		location: location,
	}
}

func (breakStatement *BreakStatement) Fix(context *Context) errors.Error {
	return nil
}

//
// return statement
//

type ReturnStatement struct {
	returnValue Expression
	location    *common.Location
}

func NewReturnStatement(value Expression, location *common.Location) *ReturnStatement {
	return &ReturnStatement{
		returnValue: value,
		location:    location,
	}
}

func (returnStatement *ReturnStatement) Fix(context *Context) errors.Error {
	function := context.GetOutFunctionDefinition()
	if function == nil {
		return errors.NewSyntaxError("Can't use return statement in global scope", returnStatement.location)
	}

	var err errors.Error
	returnStatement.returnValue, err = returnStatement.returnValue.CastTo(function.GetType(), context)
	if err != nil {
		return err
	}

	return nil
}

//
// declaration statement
//

type DeclarationStatement struct {
	// use declaration's location as statement's location
	declaration *Declaration
}

func NewDeclarationStatement(declaration *Declaration) *DeclarationStatement {
	return &DeclarationStatement{
		declaration: declaration,
	}
}

func (statement *DeclarationStatement) Fix(context *Context) {
	context.AddVariable(statement.declaration.GetName(), statement.declaration)
	statement.declaration.Fix(context)
}

//
// raw expression statement
//

type ExpressionStatement struct {
	// use expression's location as statement's location
	expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{
		expression: expression,
	}
}

func (statement *ExpressionStatement) Fix(context *Context) errors.Error {
	var err errors.Error
	statement.expression, err = statement.expression.Fix(context)
	return err
}
