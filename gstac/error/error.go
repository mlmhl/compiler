package error

import (
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