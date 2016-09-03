package ast

import (
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/executable"
	"github.com/mlmhl/goutil/encoding"
)

type Statement interface {
	Fix(context *Context) errors.Error
	Generate(context *Context, exe *executable.Executable) errors.Error
}

//
// block
//

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

func (block *baseBlock) Generate(context *Context, exe *executable.Executable) errors.Error {
	for _, statement := range block.statements {
		if err := statement.Generate(context, exe); err != nil {
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

func (elifStatement *ElifStatement) Generate(endLabel int, context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(elifStatement.location.Encode())

	if err = elifStatement.condition.Generate(context, exe); err != nil {
		return err
	}

	ifFalseLabel := exe.NewLabel()
	exe.Append(executable.JUMP_IF_FALSE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(ifFalseLabel))

	if err = elifStatement.block.Generate(context, exe); err != nil {
		return err
	}
	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(endLabel))

	exe.SetLabel(ifFalseLabel, exe.GetSize())

	return nil
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

func (elseStatement *ElseStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(elseStatement.location.Encode())
	return elseStatement.block.Generate(context, exe)
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

func (ifStatement *IfStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(ifStatement.location.Encode())

	if err = ifStatement.condition.Generate(context, exe); err != nil {
		return err
	}

	ifFalseLabel := exe.NewLabel()
	exe.Append(executable.JUMP_IF_FALSE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(ifFalseLabel))

	if err = ifStatement.ifBlock.Generate(context, exe); err != nil {
		return err
	}

	endLabel := exe.NewLabel()

	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(endLabel))

	exe.SetLabel(ifFalseLabel, exe.GetSize())

	for _, elifStatement := range(ifStatement.elifStatements) {
		if err = elifStatement.Generate(endLabel, context, exe); err != nil {
			return err
		}
	}

	ifStatement.elseStatement.Generate(context, exe)
	exe.SetLabel(endLabel, exe.GetSize())

	return nil
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

func (forStatement *ForStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(forStatement.location.Encode())

	if forStatement.init != nil {
		if err = forStatement.init.Generate(context, exe); err != nil {
			return err
		}
	}

	startLabel := exe.NewLabel()
	exe.SetLabel(startLabel, exe.GetSize())

	// endLabel is also the breakLabel
	endLabel := exe.NewLabel()
	exe.SetBreakLabel(endLabel)
	defer exe.ResetBreakLabel()

	if forStatement.condition != nil {
		if err = forStatement.condition.Generate(context, exe); err != nil {
			return err
		}
		exe.Append(executable.JUMP_IF_FALSE)
		exe.AppendSlice(encoding.DefaultEncoder.Int(endLabel))
	}

	continueLabel := exe.NewLabel()
	exe.SetContinueLabel(continueLabel)
	defer exe.ResetContinueLabel()

	if forStatement.block != nil {
		if err = forStatement.block.Generate(context, exe); err != nil {
			return err
		}
	}

	exe.SetLabel(continueLabel, exe.GetSize())

	if forStatement.post != nil {
		if err = forStatement.post.Generate(context, exe); err != nil {
			return err
		}
	}

	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(startLabel))

	exe.SetLabel(endLabel, exe.GetSize())

	return nil
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

func (whileStatement *WhileStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(whileStatement.location.Encode())

	startLabel := exe.NewLabel()
	exe.SetLabel(startLabel, exe.GetSize())

	// endLabel is also the breakLabel
	endLabel := exe.NewLabel()
	exe.SetBreakLabel(endLabel)
	defer exe.ResetBreakLabel()

	if err = whileStatement.condition.Generate(context, exe); err != nil {
		return err
	}
	exe.Append(executable.JUMP_IF_FALSE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(endLabel))

	continueLabel := exe.NewLabel()
	exe.SetContinueLabel(continueLabel)
	defer exe.ResetBreakLabel()

	if err = whileStatement.block.Generate(context, exe); err != nil {
		return err
	}

	exe.SetLabel(continueLabel, exe.GetSize())

	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(startLabel))

	exe.SetLabel(endLabel, exe.GetSize())

	return nil
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

func (continueStatement *ContinueStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(continueStatement.location.Encode())
	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(exe.GetContinueLabel()))
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

func (breakStatement *BreakStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(breakStatement.location.Encode())
	exe.Append(executable.JUMP)
	exe.AppendSlice(encoding.DefaultEncoder.Int(exe.GetBreakLabel()))
	return nil
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

func (returnStatement *ReturnStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(returnStatement.location.Encode())
	if err := returnStatement.returnValue.Generate(context, exe); err != nil {
		return err
	}
	exe.Append(executable.RETURN)
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

func (statement *DeclarationStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	return statement.declaration.Generate(context, exe)
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

func (statement *ExpressionStatement) Generate(context *Context, exe *executable.Executable) errors.Error {
	if err := statement.expression.Generate(context, exe); err != nil {
		return err
	}
	exe.Append(executable.STACK_POP)
	return nil
}
