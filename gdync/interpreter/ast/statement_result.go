package ast

import (
	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

const (
	NORMAL_STATEMENT_RESULT = iota
	RETURN_STATEMENT_RESULT
	BREAK_STATEMENT_RESULT
	CONTINUE_STATEMENT_RESULT
)

type StatementResult struct {
	typ int
	value types.Value
}

func NewStatementResult(typ int, value types.Value) *StatementResult {
	return &StatementResult{
		typ: typ,
		value: value,
	}
}

func (result *StatementResult) GetType() int {
	return result.typ
}

func (result *StatementResult) GetValue() types.Value {
	return result.value
}

func (result *StatementResult) SetType(typ int) {
	result.typ = typ
}
