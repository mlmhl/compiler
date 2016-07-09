package error

import (
	"fmt"

	"github.com/mlmhl/compiler/common"
)

type Error interface {
	GetMessage() string
	GetLocation() *common.Location
	SetLocation(location *common.Location)
}

type baseError struct {
	message  string
	location *common.Location
}

func (error *baseError) GetMessage() string {
	return error.message
}

func (error *baseError) GetLocation() *common.Location {
	return error.location
}

func (error *baseError) SetLocation(location *common.Location) {
	error.location = location
}

//
// compile error
//

type SyntaxError struct {
	baseError
}

func NewSyntaxError(message string, location *common.Location) *SyntaxError {
	return &SyntaxError{
		baseError: baseError{
			message:  message,
			location: location,
		},
	}
}

type VariableNotFoundError struct {
	baseError
}

func NewVariableNotFoundError(name string,
	location *common.Location) *VariableNotFoundError {
	return &VariableNotFoundError{
		baseError: baseError{
			message:  fmt.Sprintf("Undefined variable %s", name),
			location: location,
		},
	}
}

type VariableDuplicateDefinitionError struct {
	baseError
}

func NewVariableDuplicateDefinitionError(name string, firstLoc,
	secondLoc *common.Location) *VariableDuplicateDefinitionError {
	format := "Duplicated variable definition %s, " +
		"has been defined at %s, %d, %d"
	return &VariableDuplicateDefinitionError{
		baseError: baseError{
			message: fmt.Sprintf(format, name, firstLoc.GetFileName(),
				firstLoc.GetLine(), firstLoc.GetPosition()),
			location: secondLoc,
		},
	}
}

type FunctionNotFoundError struct {
	baseError
}

func NewFunctionNotFoundError(name string,
	location *common.Location) *FunctionNotFoundError {
	return &FunctionNotFoundError{
		baseError: baseError{
			message:  fmt.Sprintf("Undefined function %s", name),
			location: location,
		},
	}
}

type FunctionDuplicateDefinitionError struct {
	baseError
}

func NewFunctionDuplicateDefinitionError(name string, firstLoc,
	secondLoc *common.Location) *FunctionDuplicateDefinitionError {
	format := "Duplicated function definition %s, at %s, %d, %d"
	return &FunctionDuplicateDefinitionError{
		baseError: baseError{
			message: fmt.Sprintf(format, name, firstLoc.GetFileName(),
				firstLoc.GetLine(), firstLoc.GetPosition()),
			location: secondLoc,
		},
	}
}

//
// Runtime error
//

type ArgumentTooFewError struct {
	baseError
}

func NewArgumentTooFewError(name string, target, size int,
	location *common.Location) *ArgumentTooFewError {
	return &ArgumentTooFewError{
		baseError: baseError{
			message: fmt.Sprintf("Not enough arguments in call to %s, need %d,"+
				"but found %d instead", name, target, size),
			location: location,
		},
	}
}

type ArgumentTooManyError struct {
	message  string
	location *common.Location
}

func NewArgumentTooManyError(name string, target, size int,
	location *common.Location) *ArgumentTooFewError {
	return &ArgumentTooFewError{
		baseError: baseError{
			message: fmt.Sprintf("Too many arguments in call to %s, "+
				"need %d, but found %d instead", name, target, size),
			location: location,
		},
	}
}

type TypeMismatchError struct {
	baseError
}

func NewTypeMismatchError(target, typ string, value interface{},
	location *common.Location) *TypeMismatchError {
	return &TypeMismatchError{
		baseError: baseError{
			message: fmt.Sprintf("Can't use %v(type %s) as type %s",
				value, typ, target),
			location: location,
		},
	}
}

type InvalidOperationError struct {
	baseError
}

func NewInvalidOperationError(location *common.Location, op string,
	types ...string) *InvalidOperationError {
	return &InvalidOperationError{
		baseError: baseError{
			message:  fmt.Sprintf("Can't invoke %s operation on %v", op, types),
			location: location,
		},
	}
}

type GlobalStatementInTopLevelError struct {
	baseError
}

func NewGlobalStatementInTopLevelError(
	location *common.Location) *GlobalStatementInTopLevelError {
	return &GlobalStatementInTopLevelError{
		baseError: baseError{
			message:  "Can't use 'global' in global scope",
			location: location,
		},
	}
}

type NotBoolExpressionError struct {
	baseError
}

func NewNotBoolExpressionError(typ string,
	location *common.Location) *NotBoolExpressionError {
	return &NotBoolExpressionError{
		baseError: baseError{
			message: fmt.Sprintf("Condition expression for %s " +
			"must be a bool expression", typ),
			location: location,
		},
	}
}

//
// internal error
//

type InternalError struct {
	message string
}

func NewInternalError(message string) *InternalError {
	return &InternalError{
		message: message,
	}
}

func (error *InternalError) GetMessage() string {
	return error.message
}

func (error *InternalError) GetLocation() *common.Location {
	// NO-OP...
	panic("Can't invoke GetLocation on InternalError")
}

func (error *InternalError) SetLocation(location *common.Location) {
	// NO_OP...
	panic("Can't invoke SetLocation on InternalError")
}