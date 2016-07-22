package ast

import (
	"strings"

	"github.com/mlmhl/compiler/gstac/errors"
)

type Expression interface {
	Fix(context *Context)
	Generate(context *Context)
}

type NullExpression struct {
	typ Type
}

func NewNullExpression() *NullExpression {
	return &NullExpression{
		typ: NULL_TYPE,
	}
}

type BoolExpression struct {
	typ   Type
	value bool
}

func NewBoolExpression(value bool) *BoolExpression {
	return &BoolExpression{
		typ:   BOOL_TYPE,
		value: value,
	}
}

type IntegerExpression struct {
	typ   Type
	value int
}

func NewIntegerExpression(value int64) *IntegerExpression {
	return &IntegerExpression{
		typ:   INTEGER_TYPE,
		value: value,
	}
}

type FloatExpression struct {
	typ   Type
	value float64
}

func NewFloatExpression(value float64) *FloatExpression {
	return &FloatExpression{
		typ:   FLOAT_TYPE,
		value: value,
	}
}

type StringExpression struct {
	typ   Type
	value string
}

func NewStringExpression(value string) (*StringExpression, errors.Error) {
	value = strings.Trim(value, "\"")
	buffer := []byte{}
	for i := 0; i < len(value); i++ {
		b := value[i]
		if b == '\\' {
			if i == len(value)-1 {
				return nil, errors.NewSyntaxError("Invalid string value: "+value, nil)
			}
			switch value[i+1] {
			case 'n':
				b = '\n'
			case 'r':
				b = '\r'
			case 't':
				b = '\t'
			case '\\':
				b = '\\'
			default:
				return nil, errors.NewSyntaxError("Invalid string value: "+value, nil)
			}
			i++
		}
		buffer = append(buffer, b)
	}

	return &StringExpression{
		typ:   STRING_TYPE,
		value: string(buffer),
	}
}

type IdentifierExpression struct {
	identifier *Identifier
}

type CommaExpression struct {
	left  Expression
	right Expression
}

type assignExpression struct {
	left    Expression
	operand Expression
}

type NormalAssignExpression struct {
	assignExpression
}

type AddAssignExpression struct {
	assignExpression
}

type SubtractAssignExpression struct {
	assignExpression
}

type MultiplyAssignExpression struct {
	assignExpression
}

type DivideAssignExpression struct {
	assignExpression
}

type ModAssignExpression struct {
	assignExpression
}

type binaryExpression struct {
	left  Expression
	right Expression
}

type AddExpression struct {
	binaryExpression
}

type SubtractExpression struct {
	binaryExpression
}

type MultiplyExpression struct {
	binaryExpression
}

type DivideExpression struct {
	binaryExpression
}

type ModExpression struct {
	binaryExpression
}

type EqualExpression struct {
	binaryExpression
}

type NotEqualExpression struct {
	binaryExpression
}

type GreaterThanExpression struct {
	binaryExpression
}

type GreaterThanAndEqualExpression struct {
	binaryExpression
}

type LessThanExpression struct {
	binaryExpression
}

type LessThanAndEqualExpression struct {
	binaryExpression
}

type LogicalOrExpression struct {
	binaryExpression
}

type LogicalAndExpression struct {
	binaryExpression
}

type unaryExpression struct {
	operand Expression
}

type LogicalNotExpression struct {
	unaryExpression
}

type MinusExpression struct {
	unaryExpression
}

type IncrementExpression struct {
	unaryExpression
}

type DecrementExpression struct {
	unaryExpression
}

type FunctionCallExpression struct {
	identifier *Identifier
	arguments  []Argument
}

type IndexExpression struct {
	array Expression
	index Expression
}

// Literal array like int[] a = {1, 2, 3}
type LiteralArrayExpression struct {
	expressions []Expression
}

type castExpression struct {
	operand Expression
}

type IntegerToFloatCastExpression struct {
	castExpression
}

type FloatToIntegerCastExpression struct {
	castExpression
}

type BoolToStringCastExpression struct {
	castExpression
}

type IntegerToStringCastExpression struct {
	castExpression
}

type FloatToStringCastExpression struct {
	castExpression
}

type ArrayCreationExpression struct {
	typ Type
	dimensions []Expression
}