package errors

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
