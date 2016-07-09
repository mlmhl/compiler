package ast

import (
	gerror "github.com/mlmhl/compiler/gdync/errors"
)

//
// Block is a code segment surround of "{" and "}".
//

type Block struct {
	statements []Statement
}

func NewBlock(statements []Statement) *Block {
	return &Block{
		statements: statements,
	}
}

func (block *Block) Execute(env *Environment) (
	*StatementResult, gerror.Error) {
	var err gerror.Error
	var result *StatementResult

	for _, statement := range block.statements {
		result, err = statement.Execute(env)
		if err != nil {
			return nil, err
		}
		if result.GetType() != NORMAL_STATEMENT_RESULT {
			return result, nil
		}
	}

	if result == nil {
		result = NewStatementResult(NORMAL_STATEMENT_RESULT, nil)
	}
	return result, nil
}
