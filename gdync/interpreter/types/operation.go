package types

import (
	"reflect"

	gerror "github.com/mlmhl/compiler/gdync/errors"
)

//
// operations
//
const (
	ADD = "Add"
	SUBTRACT = "Subtract"
	MULTIPLY = "Multiply"
	DIVIDE = "Divide"
	MOD = "Mod"
	MINUS = "Minus"

	EQUAL = "Equal"
	NOT_EQUAL = "NotEqual"
	GT = "GreaterThan"
	LT = "LessThan"
	GTE = "GreateThanOrEqual"
	LTE = "LessThanOrEqual"

	AND = "And"
	Or = "Or"
	NOT = "Not"
)

func ArithmeticOperation(op string, left, right Value) (Value, gerror.Error) {
	return defaultBinaryOperation(op, left, right)
}

func RelationalOperation(op string, left, right Value) (Value, gerror.Error) {
	return defaultBinaryOperation(op, left, right)
}

type logicalOperation func(left, right Value) bool

var logicalOperations map[string]logicalOperation = map[string]logicalOperation{
	AND: func(left, right Value) bool {
		return left.GetValue().(bool) && right.GetValue().(bool)
	},
	Or: func(left, right Value) bool {
		return left.GetValue().(bool) || right.GetValue().(bool)
	},
}

func LogicalOperation(op string, left, right Value) (Value, gerror.Error) {
	if left.GetType() != BOOL_TYPE {
		return nil, gerror.NewInvalidOperationError(
			nil, op, left.GetType().String())
	}

	if right.GetType() != BOOL_TYPE {
		return nil, gerror.NewInvalidOperationError(
			nil, op, right.GetType().String())
	}

	return NewValue(BOOL_TYPE, logicalOperations[op](left, right)), nil
}

func defaultBinaryOperation(op string, left, right Value) (Value, gerror.Error) {
	args := []reflect.Value{reflect.ValueOf(right)}
	method := reflect.ValueOf(left).MethodByName(op + right.GetType().String())

	res := method.Call(args)
	if res == nil || len(res) != 2 {
		panic("Invalid operation result of " + op)
	}
	if res[1].Interface() != nil {
		return nil, res[1].Interface().(gerror.Error)
	} else {
		return res[0].Interface().(Value), nil
	}
}