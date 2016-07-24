package ast

import "github.com/mlmhl/compiler/common"

type Statement interface {
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

type ExpressionStatement struct {
}

type DeclarationStatement struct {
	// use declaration's location as statement's location
	declaration *Declaration
}

//
// If statement
//

type ElifStatement struct {
	condition Expression
	block     *IfBlock

	// location for 'elif' keyword
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

	// location for 'else' keyword
	location *common.Location
}

type IfStatement struct {
	condition      Expression
	ifBlock        *IfBlock
	elifStatements []*ElifStatement
	elseBlock      *ElseStatement

	// location for 'if' keyword
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

func (ifStatement *IfStatement) SetElifStatements(statements *ElifStatement) {
	ifStatement.elifStatements = statements
}

func (ifStatement *IfStatement) SetElseBlock(statement *ElseStatement) {
	ifStatement.elseBlock = statement
}

//
// while statement
//

type WhileStatement struct {
	condition *Expression
	block     *WhileBlock
}

//
// for statement
//

type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *ForBlock
}

//
// return statement
//

type ReturnStatement struct {
	returnValue Expression
	location    *common.Location
}

//
// break statement
//

type BreakStatement struct {
	location *common.Location
}

//
// continue statement
//

type ContinueStatement struct {
	location *common.Location
}
