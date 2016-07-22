package ast

import "github.com/mlmhl/compiler/common"

type Statement interface {
}

type Block struct {
	statements []Statement
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

type elifStatement struct {
	condition Expression
	block     *Block

	// location for 'elif' keyword
	location *common.Location
}

type elseStatement struct {
	block *Block

	// location for 'else' keyword
	location *common.Location
}

type ElifStatement struct {
	condition      Expression
	ifBlock        *Block
	elifStatements []*elifStatement
	elseBlock      *elseStatement

	// location for 'if' keyword
	location *common.Location
}

//
// while statement
//

type WhileStatement struct {
	condition *Expression
	block     *Block
}

//
// for statement
//

type ForStatement struct {
	init      Expression
	condition Expression
	post      Expression
	block     *Block
}

//
// return statement
//

type ReturnStatement struct {
	returnValue Expression
	location *common.Location
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