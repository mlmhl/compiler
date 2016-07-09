package ast

import (
	"strings"

	"github.com/mlmhl/compiler/common"
	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

type Expression interface {
	Evaluate(env *Environment) (types.Value, gerror.Error)
}

type StringExpression struct {
	value types.Value
}

func NewStringExpression(value string) (*StringExpression, gerror.Error) {
	value = strings.Trim(value, "\"")
	buffer := []byte{}
	for i := 0; i < len(value); i++ {
		b := value[i]
		if b == '\\' {
			if i == len(value) - 1 {
				return nil, gerror.NewSyntaxError("Invalid string value: " + value, nil)
			}
			switch value[i + 1] {
			case 'n':
				b = '\n'
			case 'r':
				b = '\r'
			case 't':
				b = '\t'
			case '\\':
				b = '\\'
			default:
				return nil, gerror.NewSyntaxError("Invalid string value: " + value, nil)
			}
			i++
		}
		buffer = append(buffer, b)
	}

	return &StringExpression{types.NewValue(types.STRING_TYPE, string(buffer))}, nil
}

func (expression *StringExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return expression.value, nil
}

type IntegerExpression struct {
	value types.Value
}

func NewIntegerExpression(value int64) *IntegerExpression {
	return &IntegerExpression{types.NewValue(types.INTEGER_TYPE, value)}
}

func (expression *IntegerExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return expression.value, nil
}

type FloatExpression struct {
	value types.Value
}

func NewFloatExpression(value float64) *FloatExpression {
	return &FloatExpression{types.NewValue(types.FLOAT_TYPE, value)}
}

func (expression *FloatExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return expression.value, nil
}

type BoolExpression struct {
	value types.Value
}

func NewBoolExpression(value bool) *BoolExpression {
	return &BoolExpression{types.NewValue(types.BOOL_TYPE, value)}
}

func (expression *BoolExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return expression.value, nil
}

type NullExpression struct {
	value types.Value
}

func NewNullExpression() *NullExpression {
	return &NullExpression{types.NewValue(types.NULL_TYPE, nil)}
}

func (expression *NullExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return expression.value, nil
}

type IdentifierExpression struct {
	identifier *types.Identifier

	// Using identifier's location as IdentifierExpression's location.
}

func NewIdentifierExpression(identifier *types.Identifier) *IdentifierExpression {
	return &IdentifierExpression{identifier}
}

func (expression *IdentifierExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	var variable *types.Variable

	if env.IsGlobal() {
		variable = env.GetGlobalVariable(expression.identifier)
	} else {
		variable = env.GetLocalVariable(expression.identifier)
	}

	if variable == nil {
		return nil, gerror.NewVariableNotFoundError(
			expression.identifier.GetName(), expression.identifier.GetLocation())
	}
	return variable.GetValue(), nil
}

type AssignExpression struct {
	operand    Expression
	identifier *types.Identifier
}

func NewAssignExpression(operand Expression, identifier *types.Identifier) *AssignExpression {
	return &AssignExpression{
		operand:    operand,
		identifier: identifier,
	}
}

func (expression *AssignExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	var err gerror.Error
	var right types.Value
	var left *types.Variable

	if right, err = expression.operand.Evaluate(env); err != nil {
		return nil, err
	}

	if env.IsGlobal() {
		left = env.GetGlobalVariable(expression.identifier)
		if left == nil {
			// Create a new global variable.
			env.AddGlobalVariable(types.NewVariable(expression.identifier, right))
		} else {
			left.SetValue(right)
		}
	} else {
		left = env.GetLocalVariable(expression.identifier)
		if left == nil {
			// Create a new local variable.
			env.AddLocalVariable(types.NewVariable(expression.identifier, right))
		} else {
			left.SetValue(right)
		}
	}

	return right, nil
}

type unaryExpression struct {
	expression Expression

	// operation signal's location
	location *common.Location
}

type binaryExpression struct {
	left  Expression
	right Expression

	// operation signal's location
	location *common.Location
}

type AddExpression struct {
	binaryExpression
}

func NewAddExpression(left, right Expression, location *common.Location) *AddExpression {
	return &AddExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *AddExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return arithmeticOperation(types.ADD, expression.left,
		expression.right, expression.location, env)
}

type SubtractExpression struct {
	binaryExpression
}

func NewSubtractExpression(left, right Expression, location *common.Location) *SubtractExpression {
	return &SubtractExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *SubtractExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return arithmeticOperation(types.SUBTRACT, expression.left,
		expression.right, expression.location, env)
}

type MultiplyExpression struct {
	binaryExpression
}

func NewMultiplyExpression(left, right Expression, location *common.Location) *MultiplyExpression {
	return &MultiplyExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *MultiplyExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return arithmeticOperation(types.MULTIPLY, expression.left,
		expression.right, expression.location, env)
}

type DivideExpression struct {
	binaryExpression
}

func NewDivideExpression(left, right Expression, location *common.Location) *DivideExpression {
	return &DivideExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *DivideExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return arithmeticOperation(types.DIVIDE, expression.left,
		expression.right, expression.location, env)
}

type ModExpression struct {
	binaryExpression
}

func NewModExpression(left, right Expression, location *common.Location) *ModExpression {
	return &ModExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *ModExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return arithmeticOperation(types.MOD, expression.left,
		expression.right, expression.location, env)
}

type MinusExpression struct {
	unaryExpression
}

func NewMinusExpression(expression Expression, location *common.Location) *MinusExpression {
	return &MinusExpression{
		unaryExpression: unaryExpression{
			expression: expression,
			location:   location,
		},
	}
}

func (expression *MinusExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	value, err := expression.expression.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if value.GetType() == types.INTEGER_TYPE {
		return types.NewValue(types.INTEGER_TYPE, -value.GetValue().(int64)), nil
	}
	if value.GetType() == types.FLOAT_TYPE {
		return types.NewValue(types.FLOAT_TYPE, -value.GetValue().(int64)), nil
	}

	return nil, gerror.NewInvalidOperationError(expression.location, types.MINUS,
		value.GetType().String())
}

type EqualExpression struct {
	binaryExpression
}

func NewEqualExpression(left, right Expression, location *common.Location) *EqualExpression {
	return &EqualExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *EqualExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.EQUAL, expression.left,
		expression.right, expression.location, env)
}

type NotEqualExpression struct {
	binaryExpression
}

func NewNotEqualExpression(left, right Expression, location *common.Location) *NotEqualExpression {
	return &NotEqualExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *NotEqualExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.NOT_EQUAL, expression.left,
		expression.right, expression.location, env)
}

type GTExpression struct {
	binaryExpression
}

func NewGTExpression(left, right Expression, location *common.Location) *GTExpression {
	return &GTExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *GTExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.GT, expression.left,
		expression.right, expression.location, env)
}

type LTExpression struct {
	binaryExpression
}

func NewLTExpression(left, right Expression, location *common.Location) *LTExpression {
	return &LTExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *LTExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.LT, expression.left,
		expression.right, expression.location, env)
}

type GTEExpression struct {
	binaryExpression
}

func NewGTEExpression(left, right Expression, location *common.Location) *GTEExpression {
	return &GTEExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *GTEExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.GTE, expression.left,
		expression.right, expression.location, env)
}

type LTEExpression struct {
	binaryExpression
}

func NewLTEExpression(left, right Expression, location *common.Location) *LTEExpression {
	return &LTEExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *LTEExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return relationalOperation(types.LTE, expression.left,
		expression.right, expression.location, env)
}

type AndExpression struct {
	binaryExpression
}

func NewAndExpression(left, right Expression, location *common.Location) *AndExpression {
	return &AndExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *AndExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return logicalOperation(types.AND, expression.left, expression.right, expression.location, env)
}

type OrExpression struct {
	binaryExpression
}

func NewOrExpression(left, right Expression, location *common.Location) *OrExpression {
	return &OrExpression{
		binaryExpression: binaryExpression{
			left:     left,
			right:    right,
			location: location,
		},
	}
}

func (expression *OrExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	return logicalOperation(types.Or, expression.left, expression.right, expression.location, env)
}

type NotExpression struct {
	unaryExpression
}

func NewNotExpression(expression Expression, location *common.Location) *NotExpression {
	return &NotExpression{
		unaryExpression: unaryExpression{
			expression: expression,
			location:   location,
		},
	}
}

func (expression *NotExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	value, err := expression.expression.Evaluate(env)
	if err != nil {
		return nil, err
	}

	if value.GetType() != types.BOOL_TYPE {
		return nil, gerror.NewInvalidOperationError(expression.location, types.NOT,
			value.GetType().String())
	}

	return types.NewValue(types.BOOL_TYPE, !value.GetValue().(bool)), nil
}

type FunctionCallExpression struct {
	arguments  []*Argument
	identifier *types.Identifier // function's identifier

	location *common.Location // function call point location
}

func NewFunctionCallExpression(arguments []*Argument, identifier *types.Identifier,
	location *common.Location) *FunctionCallExpression {
	return &FunctionCallExpression{
		arguments: arguments,
		identifier: identifier,

		location: location,
	}
}

func (expression *FunctionCallExpression) Evaluate(env *Environment) (types.Value, gerror.Error) {
	function := env.GetFunction(expression.identifier)
	if function == nil {
		return nil, gerror.NewFunctionNotFoundError(
			expression.identifier.GetName(), expression.location)
	}

	localEnv := NewEnvironment(env.GetGlobalVariables(), env.GetFunctions(), false)

	values := []types.Value{}
	for _, argument := range expression.arguments {
		value, err := argument.expression.Evaluate(env)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	value, err := function.Evaluate(values, localEnv)
	if err != nil && err.GetLocation() == nil {
		err.SetLocation(expression.location)
	}
	return value, err
}

//
// arithmetic operators
//

func arithmeticOperation(op string, leftExpression, rightExpression Expression,
	location *common.Location, env *Environment) (types.Value, gerror.Error) {
	left, right, err := evaluateOperand(leftExpression, rightExpression, env)
	if err != nil {
		return nil, err
	}

	if value, err := types.ArithmeticOperation(op, left, right); err != nil {
		err.SetLocation(location)
		return nil, err
	} else {
		return value, nil
	}
}

//
// relational operators
//

func relationalOperation(op string, leftExpression, rightExpression Expression,
	location *common.Location, env *Environment) (types.Value, gerror.Error) {

	left, right, err := evaluateOperand(leftExpression, rightExpression, env)
	if err != nil {
		return nil, err
	}

	if value, err := types.RelationalOperation(op, left, right); err != nil {
		err.SetLocation(location)
		return nil, err
	} else {
		return value, nil
	}
}

//
// logic operations
//

func logicalOperation(op string, leftExpression, rightExpression Expression,
	location *common.Location, env *Environment) (types.Value, gerror.Error) {
	left, right, err := evaluateOperand(leftExpression, rightExpression, env)
	if err != nil {
		return nil, err
	}

	value, err := types.LogicalOperation(op, left, right)
	if err != nil {
		err.SetLocation(location)
		return nil, err
	} else {
		return value, nil
	}
}

func evaluateOperand(leftExpression, rightExpression Expression,
	env *Environment) (types.Value, types.Value, gerror.Error) {
	var err gerror.Error
	var left types.Value
	var right types.Value

	if left, err = leftExpression.Evaluate(env); err != nil {
		return nil, nil, err
	}
	if right, err = rightExpression.Evaluate(env); err != nil {
		return nil, nil, err
	}

	return left, right, nil
}
