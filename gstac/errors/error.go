package errors

import (
	"fmt"

	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/compiler/ast"
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

type DuplicateDeclarationError struct {
	baseError
}

func NewDuplicateDeclarationError(name string, firstLocation,
	secondLocation *common.Location) *DuplicateDeclarationError {
	return &DuplicateDeclarationError{
		baseError: baseError{
			message: fmt.Sprintf("Duplicated declaration %s, has been declared at %s, %d, %d",
				name, firstLocation.GetFileName(), firstLocation.GetLine(), firstLocation.GetPosition()),
			location: secondLocation,
		},
	}
}

type DuplicateFunctionDefinitionError struct {
	baseError
}

func NewDuplicateFunctionDefinitionError(name string, firstLocation,
	secondLocation *common.Location) *DuplicateFunctionDefinitionError {
	return &DuplicateFunctionDefinitionError{
		baseError: baseError{
			message: fmt.Sprintf("Duplicated function definition %s, has been definited at %s, %d, %d",
				name, firstLocation.GetFileName(), firstLocation.GetLine(), firstLocation.GetPosition()),
			location: secondLocation,
		},
	}
}

type UnsupportedTypeError struct {
	baseError
}

func NewUnsupportedTypeError(typeName string, location *common.Location) *UnsupportedTypeError {
	return &UnsupportedTypeError{
		baseError: baseError{
			message:  fmt.Sprintf("Can't use %s as a type keyword", typeName),
			location: location,
		},
	}
}

type ParenthesesNotMatchedError struct {
	baseError
}

func NewParenthesesNotMatchedError(leftDesc, rightDesc string,
	leftLoc, rightLoc *common.Location) *ParenthesesNotMatchedError {
	return &ParenthesesNotMatchedError{
		baseError: baseError{
			message: fmt.Sprintf("Can't find %s to match %s at %s, %d, %d", rightDesc,
				leftDesc, leftLoc.GetFileName(), leftLoc.GetLine(), leftLoc.GetPosition()),
			location: rightLoc,
		},
	}
}

type TypeCastError struct {
	baseError
}

func NewTypeCastError(srcType, destType ast.Type, location *common.Location) *TypeCastError {
	return &TypeCastError{
		baseError: baseError{
			message:  fmt.Sprintf("can't cast %s to %s", srcType.GetName(), destType.GetName()),
			location: location,
		},
	}
}

type InvalidOperationError struct {
	baseError
}

func NewInvalidOperationError(op string, location *common.Location, types ...string) *InvalidOperationError {
	return &InvalidOperationError{
		baseError: baseError{
			message:  fmt.Sprintf("Can't invoke %s operation on %v", op, types),
			location: location,
		},
	}
}

type FunctionNotFoundError struct {
	baseError
}

func NewFunctionNotFoundError(name string, location *common.Location) *FunctionNotFoundError {
	return &FunctionNotFoundError{
		baseError: baseError{
			message:  "Can't find function " + name,
			location: location,
		},
	}
}

type ArgumentCountMismatchError struct {
	baseError
}

func NewArgumentCountMismatchError(paramCount, argCount int,
	location *common.Location) *ArgumentCountMismatchError {
	return &ArgumentCountMismatchError{
		baseError: baseError{
			message: fmt.Sprintf("Target parameter count is %d, but argument size is %d instead",
				argCount, paramCount),
			location: location,
		},
	}
}

type InvalidTypeError struct {
	baseError
}

func NewInvalidTypeError(typ, target string, location *common.Location) *InvalidTypeError {
	return &InvalidTypeError{
		baseError: baseError{
			message:  fmt.Sprintf("Can't use %s as %s", typ, target),
			location: location,
		},
	}
}

type IndexNotIntError struct {
	baseError
}

func NewIndexNotIntError(typ string, location *common.Location) *IndexNotIntError {
	return &IndexNotIntError{
		baseError: baseError{
			message:  fmt.Sprintf("Can't use %s as array index", typ),
			location: location,
		},
	}
}

type ArraySizeNotIntError struct {
	baseError
}

func NewArraySizeNotIntError(typ string, location *common.Location) *ArraySizeNotIntError {
	return &ArraySizeNotIntError{
		baseError: &baseError{
			message:  fmt.Sprintf("Can't use %s as array size", typ),
			location: location,
		},
	}
}

type ParameterDuplicatedDefinitionError struct {
	baseError
}

func NewParameterDuplicatedDefinitionError(name string, firstLoc,
	secondLoc *common.Location) *ParameterDuplicatedDefinitionError {
	return &ParameterDuplicatedDefinitionError{
		baseError: baseError{
			message: fmt.Sprintf("Parameter %s has been defined at %s, %d, %d", name,
			firstLoc.GetFileName(), firstLoc.GetLine(), firstLoc.GetPosition()),
			location: secondLoc,
		},
	}
}
