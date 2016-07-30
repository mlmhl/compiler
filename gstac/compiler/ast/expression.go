package ast

import (
	"strings"

	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/token"
)

type Expression interface {
	Fix(context *Context)
	TypeCast(destType Type) Expression
	Generate(context *Context)
}

//
// value expression
//

type valueExpression struct {
	typ Type
}

func (expression *valueExpression) Fix(context *Context) {
	// NO OP
}

type NullExpression struct {
	valueExpression
}

func NewNullExpression() *NullExpression {
	return &NullExpression{
		valueExpression: valueExpression{
			typ: NULL_TYPE,
		},
	}
}

func (expression *NullExpression) TypeCast(destType Type) Expression {

}

type BoolExpression struct {
	valueExpression
	value bool
}

func NewBoolExpression(value bool) *BoolExpression {
	return &BoolExpression{
		valueExpression: valueExpression{
			typ: BOOL_TYPE,
		},
		value: value,
	}
}

type IntegerExpression struct {
	valueExpression
	value int
}

func NewIntegerExpression(value int64) *IntegerExpression {
	return &IntegerExpression{
		valueExpression: valueExpression{
			typ: INTEGER_TYPE,
		},
		value: value,
	}
}

type FloatExpression struct {
	valueExpression
	value float64
}

func NewFloatExpression(value float64) *FloatExpression {
	return &FloatExpression{
		valueExpression: valueExpression{
			typ: FLOAT_TYPE,
		},
		value: value,
	}
}

type StringExpression struct {
	valueExpression
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
		valueExpression: valueExpression{
			typ: STRING_TYPE,
		},
		value: string(buffer),
	}
}

// Literal array like int[] a = {1, 2, 3}
type ArrayLiteralExpression struct {
	values []Expression
}

func NewArrayLiteralExpression(values []Expression) *ArrayLiteralExpression {
	return &ArrayLiteralExpression{
		values: values,
	}
}

func (expression *ArrayLiteralExpression) Fix(context *Context) {
	for _, value := range(expression.values) {
		value.Fix(context)
	}
}

type IdentifierExpression struct {
	identifier *Identifier

	// use identifier's location as expression's location
}

func NewIdentifierExpression(identifier *Identifier) *IdentifierExpression {
	return &IdentifierExpression{
		identifier: identifier,
	}
}

type assignExpression struct {
	left    Expression
	operand Expression

	// use left expression's location as the whole expression's location
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

func NewAssignExpression(typ int, left, operand Expression) Expression {
	assignExpression := assignExpression{
		left:    left,
		operand: operand,
	}
	switch typ {
	case token.ASSIGN_ID:
		return &NormalAssignExpression{
			assignExpression: assignExpression,
		}
	case token.ADD_ASSIGN_ID:
		return &AddAssignExpression{
			assignExpression: assignExpression,
		}
	case token.SUBTRACT_ID:
		return &SubtractAssignExpression{
			assignExpression: assignExpression,
		}
	case token.MUL_ASSIGN_ID:
		return &MultiplyAssignExpression{
			assignExpression: assignExpression,
		}
	case token.DIV_ASSIGN_ID:
		return &DivideAssignExpression{
			assignExpression: assignExpression,
		}
	case token.MOD_ASSIGN_ID:
		return &ModAssignExpression{
			assignExpression: assignExpression,
		}
	default:
		return nil
	}
}

type binaryExpression struct {
	left  Expression
	right Expression

	// use left expression's location as the whole expression's location
}

type AddExpression struct {
	binaryExpression
}

func NewAddExpression(left, right Expression) *AddExpression {
	return &AddExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type SubtractExpression struct {
	binaryExpression
}

func NewSubtractExpression(left, right Expression) *SubtractExpression {
	return &SubtractExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type MultiplyExpression struct {
	binaryExpression
}

func NewMultiplyExpression(left, right Expression) *MultiplyExpression {
	return &MultiplyExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type DivideExpression struct {
	binaryExpression
}

func NewDivideExpression(left, right Expression) *DivideExpression {
	return &DivideExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type ModExpression struct {
	binaryExpression
}

func NewModExpression(left, right Expression) *ModExpression {
	return &ModExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type EqualExpression struct {
	binaryExpression
}

func NewEqualExpression(left, right Expression) *EqualExpression {
	return EqualExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type NotEqualExpression struct {
	binaryExpression
}

func NewNotEqualExpression(left, right Expression) *NotEqualExpression {
	return NotEqualExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type GreaterThanExpression struct {
	binaryExpression
}

func NewGreaterThanExpression(left, right Expression) *GreaterThanExpression {
	return GreaterThanExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type GreaterThanAndEqualExpression struct {
	binaryExpression
}

func NewGreaterThanAndEqualExpression(left, right Expression) *GreaterThanAndEqualExpression {
	return GreaterThanAndEqualExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type LessThanExpression struct {
	binaryExpression
}

func NewLessThanExpression(left, right Expression) *LessThanExpression {
	return LessThanExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type LessThanAndEqualExpression struct {
	binaryExpression
}

func NewLessThanAndEqualExpression(left, right Expression) *LessThanAndEqualExpression {
	return LessThanAndEqualExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type LogicalOrExpression struct {
	binaryExpression
}

func NewLogicalOrExpression(left, right Expression) *LogicalOrExpression {
	return &LogicalOrExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type LogicalAndExpression struct {
	binaryExpression
}

func NewLogicalAndExpression(left, right Expression) *LogicalAndExpression {
	return &LogicalAndExpression{
		binaryExpression: binaryExpression{
			left:  left,
			right: right,
		},
	}
}

type unaryExpression struct {
	operand Expression
}

type LogicalNotExpression struct {
	unaryExpression
}

func NewLogicalNotExpression(operand Expression) *LogicalOrExpression {
	return &LogicalOrExpression{
		unaryExpression{
			operand: operand,
		},
	}
}

type MinusExpression struct {
	unaryExpression
}

func NewMinusExpression(operand Expression) *MinusExpression {
	return &MinusExpression{
		unaryExpression: unaryExpression{
			operand: operand,
		},
	}
}

type IncrementExpression struct {
	unaryExpression
}

func NewIncrementExpression(operand Expression) *IncrementExpression {
	return &IncrementExpression{
		unaryExpression: unaryExpression{
			operand: operand,
		},
	}
}

type DecrementExpression struct {
	unaryExpression
}

func NewDecrementExpression(operand Expression) *DecrementExpression {
	return &DecrementExpression{
		unaryExpression: unaryExpression{
			operand: operand,
		},
	}
}

type FunctionCallExpression struct {
	identifier *Identifier
	arguments  []*Argument
}

func NewFunctionCallExpression(identifier *Identifier,
	arguments []*Argument) *FunctionCallExpression {
	return &FunctionCallExpression{
		identifier: identifier,
		arguments: arguments,
	}
}

type IndexExpression struct {
	array Expression
	index Expression
}

func NewIndexExpression(array, index Expression) *IndexExpression {
	return &IndexExpression{
		array: array,
		index: index,
	}
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
	typ        Type
	dimensions []Expression
}

func NewArrayCreationExpression(typ Type, dimensions []Expression) *ArrayCreationExpression {
	return &ArrayCreationExpression{
		typ: typ,
		dimensions: dimensions,
	}
}
