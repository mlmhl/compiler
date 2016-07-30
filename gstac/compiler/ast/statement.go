package ast

import "github.com/mlmhl/compiler/common"

type Statement interface {
	Fix(context *Context)
	Generate()
}

//
// block
//

type UndefinedBlock struct {
}

func NewUndefinedBlock() *UndefinedBlock {
	return &IfBlock{}
}

type IfBlock struct {
	statements []Statement
}

func NewIfBlock(statementList []Statement) *IfBlock {
	return &IfBlock{
		statements: statementList,
	}
}

type ForBlock struct {
	statements []Statement
}

func NewForBlock(statementList []Statement) *ForBlock {
	return &ForBlock{
		statements: statementList,
	}
}

type WhileBlock struct {
	statements []Statement
}

func NewWhileBlock(statementList []Statement) *WhileBlock {
	return &WhileBlock{
		statements: statementList,
	}
}

type FunctionBlock struct {
	statements []Statement
}

func NewFunctionBlock(statementList []Statement) *FunctionBlock {
	return &FunctionBlock{
		statements: statementList,
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

type IfStatement struct {
	condition      Expression
	ifBlock        *IfBlock
	elifStatements []*ElifStatement
	elseBlock      *ElseStatement

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
	ifStatement.elseBlock = statement
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
		block: block,
		location: location,
	}
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
		location: location,
	}
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
	statement.declaration.TypeCast()
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

func (statement *ExpressionStatement) Fix(context *Context) {
	statement.expression.Fix(context)
}