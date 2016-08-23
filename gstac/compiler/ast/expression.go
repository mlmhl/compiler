package ast

import (
	"reflect"
	"strconv"
	"strings"

	"fmt"
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/executable"
	"github.com/mlmhl/compiler/gstac/token"
	"github.com/mlmhl/goutil/encoding"
)

type Expression interface {
	Fix(context *Context) (Expression, errors.Error)
	CastTo(destType Type, context *Context) (Expression, errors.Error)
	Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error)

	getType(context *Context) (Type, errors.Error)
	getLocation() *common.Location
}

//
// base expression
//

type baseExpression struct {
	this Expression
}

func (expression *baseExpression) Fix(context *Context) (Expression, errors.Error) {
	// No-OP
	return expression.this, nil
}

func (expression *baseExpression) CastTo(destType Type, context *Context) (Expression, errors.Error) {
	srcType, err := expression.this.getType(context)
	if err != nil {
		return expression.this, err
	}
	return typeCast(srcType, destType, expression.this)
}

func (expression *baseExpression) Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error) {
	// NO-OP
	return nil, nil
}

func (expression *baseExpression) getType(context *Context) (Type, errors.Error) {
	panic("Can't invoke `getType` on " + reflect.TypeOf(expression).Elem().Name())
}

//
// value expression
//

type valueExpression interface {
	getValue() interface{}

	// support for code byte generation

	// return according operation's code byte
	// Normally return the constant in executable.code
	operatorCode() byte

	// return the value's code byte
	valueEncode(exe *executable.Executable) []byte
}

type baseValueExpression struct {
	baseExpression
	typ      Type
	location *common.Location
}

func (expression *baseValueExpression) Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error) {
	buffer := []byte{}
	valueExpression := expression.this.(valueExpression)

	buffer = append(buffer, expression.location.Encode()...)
	buffer = append(buffer, valueExpression.operatorCode())
	buffer = append(buffer, valueExpression.valueEncode(exe)...)
	return buffer, nil
}

// Default implement for all value expressions. Just store the
// value in constant pool and encode the according index to code byte.
func (expression *baseValueExpression) valueEncode(exe *executable.Executable) []byte {
	return encoding.DefaultEncoder.Int(
		exe.AddConstantValue(expression.this.(valueExpression).getValue()))
}

func (expression *baseValueExpression) getType(context *Context) (Type, errors.Error) {
	return expression.typ, nil
}

func (expression *baseValueExpression) getLocation() *common.Location {
	return expression.location
}

type NullExpression struct {
	baseValueExpression
}

func NewNullExpression(location *common.Location) *NullExpression {
	expression := &NullExpression{
		baseValueExpression: baseValueExpression{
			typ:      NULL_TYPE,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *NullExpression) getValue() interface{} {
	return nil
}

func (expression *NullExpression) operatorCode() byte {
	return executable.PUSH_NULL
}

// Need't store null value
func (expression *NullExpression) valueEncode() []byte {
	return []byte{}
}

type BoolExpression struct {
	baseValueExpression
	value bool
}

func NewBoolExpression(value bool, location *common.Location) *BoolExpression {
	expression := &BoolExpression{
		baseValueExpression: baseValueExpression{
			typ:      BOOL_TYPE,
			location: location,
		},
		value: value,
	}
	expression.this = expression
	return expression
}

func (expression *BoolExpression) getValue() interface{} {
	return expression.value
}

func (expression *BoolExpression) operatorCode() byte {
	if expression.value {
		return executable.PUSH_BOOL_TRUE
	} else {
		return executable.PUSH_BOOL_FALSE
	}
}

// Needn't store bool value
func (expression *BoolExpression) valueEncode() []byte {
	return []byte{}
}

type IntegerExpression struct {
	baseValueExpression
	value int64
}

func NewIntegerExpression(value int64, location *common.Location) *IntegerExpression {
	expression := &IntegerExpression{
		baseValueExpression: baseValueExpression{
			typ:      INTEGER_TYPE,
			location: location,
		},
		value: value,
	}
	expression.this = expression
	return expression
}

func (expression *IntegerExpression) getValue() interface{} {
	return expression.value
}

func (expression *IntegerExpression) operatorCode() byte {
	return executable.PUSH_INT
}

type FloatExpression struct {
	baseValueExpression
	value float64
}

func NewFloatExpression(value float64, location *common.Location) *FloatExpression {
	expression := &FloatExpression{
		baseValueExpression: baseValueExpression{
			typ:      FLOAT_TYPE,
			location: location,
		},
		value: value,
	}
	expression.this = expression
	return expression
}

func (expression *FloatExpression) getValue() interface{} {
	return expression.value
}

func (expression *FloatExpression) operatorCode() []byte {
	return executable.PUSH_FLOAT
}

type StringExpression struct {
	baseValueExpression
	value string
}

func NewStringExpression(value string, location *common.Location) (*StringExpression, errors.Error) {
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

	expression := &StringExpression{
		baseValueExpression: baseValueExpression{
			typ:      STRING_TYPE,
			location: location,
		},
		value: string(buffer),
	}
	expression.this = expression
	return expression, nil
}

func (expression *StringExpression) getValue() interface{} {
	return expression.value
}

func (expression *StringExpression) operatorCode() byte {
	return executable.PUSH_STRING
}

// Literal array like int[] a = `{1, 2, 3}`
type ArrayLiteralExpression struct {
	baseExpression
	values []Expression

	// type according to array declaration
	typ Type

	// location of left large parentheses
	location *common.Location
}

func NewArrayLiteralExpression(values []Expression,
	location *common.Location) *ArrayLiteralExpression {
	expression := &ArrayLiteralExpression{
		values:   values,
		location: location,
	}
	expression.this = expression
	return expression
}

func (expression *ArrayLiteralExpression) Fix(context *Context) (Expression, errors.Error) {
	for i, value := range expression.values {
		newExpr, err := value.Fix(context)
		if err != nil {
			return expression, err
		}
		expression.values[i] = newExpr
	}
	return expression, nil
}

func (expression *ArrayLiteralExpression) CastTo(destType Type, context *Context) (Expression, errors.Error) {
	expression.typ = destType
	var err errors.Error
	for i, exp := range expression.values {
		expression.values[i], err = exp.CastTo(destType, context)
		if err != nil {
			return expression, err
		}
	}
	return expression, nil
}

// Generate code byte for each array element expression
// and generate array's size into ArrayLiteralExpression's code byte.
func (expression *ArrayLiteralExpression) Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error) {
	var buf []byte
	var buffer []byte
	var err errors.Error

	// Generate code byte of array elements
	for _, subExpr := range(expression.values) {
		buf, err = subExpr.Generate(context, exe)
		if err != nil {
			return buffer, err
		}
		buffer = append(buffer, buf...)
	}

	// Generate code byte of ArrayLiteralExpression

	buffer = append(buffer, expression.location.Encode()...)
	// Generate is called after CastTo, so typ is already set to the correct value
	buffer = append(buffer, executable.GetOperatorCode(
		executable.NEW_ARRAY_LITERAL_BOOL, expression.typ.GetOffSet()))

	return buffer, nil
}

// If typ is not set, getType will return a nil, but this doesn't matter,
// as getType won't be called up to now.
func (expression *ArrayLiteralExpression) getType(context *Context) (Type, errors.Error) {
	return expression.typ
}

func (expression *ArrayLiteralExpression) getLocation() *common.Location {
	return expression.location
}

// variable reference expression
type IdentifierExpression struct {
	baseExpression
	identifier *Identifier

	// use identifier's location as expression's location
}

func NewIdentifierExpression(identifier *Identifier) *IdentifierExpression {
	expression := &IdentifierExpression{
		identifier: identifier,
	}
	expression.this = expression
	return expression
}

func (expression *IdentifierExpression) Generate(context *Context, exe *executable.Executable) ([]byte, errors.Error) {
	buffer := []byte{}

	buffer = append(buffer, expression.getLocation().Encode()...)
	buffer = append(buffer, executable.VARIABLE_REFERENCE)
	buffer = append(buffer, encoding.DefaultEncoder.Int(context.GetSymbolIndex(expression.identifier.GetName())))

	return buffer, nil
}

func (expression *IdentifierExpression) getLocation() *common.Location {
	return expression.identifier.GetLocation()
}

type assignExpression struct {
	baseExpression
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
	expression := assignExpression{
		left:    left,
		operand: operand,
	}
	switch typ {
	case token.ASSIGN_ID:
		result := &NormalAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	case token.ADD_ASSIGN_ID:
		result := &AddAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	case token.SUBTRACT_ID:
		result := &SubtractAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	case token.MUL_ASSIGN_ID:
		result := &MultiplyAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	case token.DIV_ASSIGN_ID:
		result := &DivideAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	case token.MOD_ASSIGN_ID:
		result := &ModAssignExpression{
			assignExpression: expression,
		}
		result.this = result
		return result
	default:
		return nil
	}
}

func (expression *assignExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error
	expression.left, err = expression.left.Fix(context)
	if err != nil {
		return expression, err
	}
	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return expression, err
	}
	srcType, err := expression.left.getType(context)
	if err != nil {
		return expression, err
	}
	expression.operand, err = expression.operand.CastTo(srcType, context)
	return expression, err
}

func (expression *assignExpression) getType(context *Context) (Type, errors.Error) {
	return expression.left.getType(context)
}

func (expression *assignExpression) getLocation() *common.Location {
	return expression.left.getLocation()
}

type binaryExpression interface {
	execute(left, right interface{}, location *common.Location) (Expression, errors.Error)
}

type baseBinaryExpression struct {
	this  binaryExpression
	left  Expression
	right Expression

	// use left expression's location as the whole expression's location
}

func (expression *baseBinaryExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error
	expression.left, err = expression.left.Fix(context)
	if err != nil {
		return expression, err
	}
	expression.right, err = expression.right.Fix(context)
	if err != nil {
		return expression, err
	}

	// constant folding
	l, lok := expression.left.(valueExpression)
	r, rok := expression.right.(valueExpression)
	if lok && rok {
		return expression.this.execute(l.getValue(), r.getValue(), expression.left.getLocation())
	}

	var leftType Type
	leftType, err = expression.left.getType(context)
	if err != nil {
		return expression, err
	}

	var rightType Type
	rightType, err = expression.right.getType(context)
	if err != nil {
		return expression, err
	}

	if leftType.isPriorityOf(rightType) {
		expression.right, err = expression.right.CastTo(leftType, context)
	} else {
		expression.left, err = expression.left.CastTo(rightType, context)
	}
	return expression, err
}

func (expression *baseBinaryExpression) CastTo(destType Type, context *Context) (Expression, errors.Error) {
	srcType, err := expression.left.getType(context)
	if err != nil {
		return expression, err
	}

	rightType, err := expression.right.getType(context)
	if err != nil {
		return expression, err
	}

	if rightType.isPriorityOf(srcType) {
		srcType = rightType
	}
	return typeCast(srcType, destType, expression)
}

func (expression *baseBinaryExpression) getType(context *Context) (Type, errors.Error) {
	leftType, err := expression.left.getType(context)
	if err != nil {
		return leftType, err
	}

	rightType, err := expression.right.getType(context)
	if err != nil {
		return rightType, err
	}

	if leftType.isPriorityOf(rightType) {
		return leftType, nil
	} else {
		return rightType, nil
	}
}

func (expression *baseBinaryExpression) getLocation() *common.Location {
	return expression.left.getLocation()
}

type AddExpression struct {
	baseBinaryExpression
}

func NewAddExpression(left, right Expression) *AddExpression {
	expression := &AddExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *AddExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "ADD"

	if left != nil {
		tmp := left
		left = right
		right = tmp
	}
	if left == nil {
		if right == nil {
			return NewNullExpression(location), nil
		}
		if r, ok := right.(string); ok {
			return NewStringExpression("nil"+r, location)
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "null", reflect.TypeOf(right).Name())
		}
	}

	switch l := left.(type) {
	case bool:
		switch r := right.(type) {
		case string:
			if l {
				return NewStringExpression("true"+r, location)
			} else {
				return NewStringExpression("false"+r, location)
			}
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "bool", reflect.TypeOf(right).Name())
		}
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l+r, location), nil
		case float64:
			return NewFloatExpression(float64(l)+r, location), nil
		case string:
			return NewStringExpression(strconv.Itoa(int(l))+r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewFloatExpression(l+float64(r), location), nil
		case float64:
			return NewFloatExpression(l+r, location), nil
		case string:
			return NewStringExpression(fmt.Sprintf("%f", l)+r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	case string:
		switch r := right.(type) {
		case bool:
			if r {
				return NewStringExpression(l+"true", location)
			} else {
				return NewStringExpression(l+"false", location)
			}
		case int64:
			return NewStringExpression(l+strconv.Itoa(int(r)), location)
		case float64:
			return NewStringExpression(l+fmt.Sprintf("%f", r), location)
		case string:
			return NewStringExpression(l+r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "string", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type SubtractExpression struct {
	baseBinaryExpression
}

func NewSubtractExpression(left, right Expression) *SubtractExpression {
	expression := &SubtractExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *SubtractExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "SUBTRACT"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l-r, location), nil
		case float64:
			return NewFloatExpression(float64(l)-r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewFloatExpression(l-float64(r), location), nil
		case float64:
			return NewFloatExpression(l-r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type MultiplyExpression struct {
	baseBinaryExpression
}

func NewMultiplyExpression(left, right Expression) *MultiplyExpression {
	expression := &MultiplyExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *MultiplyExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "MULTIPLY"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l*r, location), nil
		case float64:
			return NewFloatExpression(float64(l)*r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewFloatExpression(l*float64(r), location), nil
		case float64:
			return NewFloatExpression(l*r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type DivideExpression struct {
	baseBinaryExpression
}

func NewDivideExpression(left, right Expression) *DivideExpression {
	expression := &DivideExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *DivideExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "DIVIDE"
	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l*r, location), nil
		case float64:
			return NewFloatExpression(float64(l)/r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewFloatExpression(l/float64(r), location), nil
		case float64:
			return NewFloatExpression(l/r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type ModExpression struct {
	baseBinaryExpression
}

func NewModExpression(left, right Expression) *ModExpression {
	expression := &ModExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *ModExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "MOD"
	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l%r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type EqualExpression struct {
	baseBinaryExpression
}

func NewEqualExpression(left, right Expression) *EqualExpression {
	expression := &EqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *EqualExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "EQUAL"

	if left != nil {
		tmp := left
		left = right
		right = tmp
	}

	if left == nil {
		if right == nil {
			return NewBoolExpression(true, location), nil
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "null", reflect.TypeOf(right).Name())
		}
	}

	switch l := left.(type) {
	case bool:
		if r, ok := right.(bool); ok {
			return NewBoolExpression(l == r, location), nil
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "bool", reflect.TypeOf(right).Name())
		}
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l == r, location), nil
		case float64:
			return NewBoolExpression(float64(l) == r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l == float64(r), location), nil
		case float64:
			return NewBoolExpression(l == r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type NotEqualExpression struct {
	baseBinaryExpression
}

func NewNotEqualExpression(left, right Expression) *NotEqualExpression {
	expression := &NotEqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *NotEqualExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "UNEQUAL"

	if left != nil {
		tmp := left
		left = right
		right = tmp
	}

	if left == nil {
		if right == nil {
			return NewBoolExpression(false, location), nil
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "null", reflect.TypeOf(right).Name())
		}
	}

	switch l := left.(type) {
	case bool:
		if r, ok := right.(bool); ok {
			return NewBoolExpression(l != r, location), nil
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "bool", reflect.TypeOf(right).Name())
		}
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l != r, location), nil
		case float64:
			return NewBoolExpression(float64(l) != r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l != float64(r), location), nil
		case float64:
			return NewBoolExpression(l != r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type GreaterThanExpression struct {
	baseBinaryExpression
}

func NewGreaterThanExpression(left, right Expression) *GreaterThanExpression {
	expression := &GreaterThanExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *GreaterThanExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "GREARER_THAN"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l > r, location), nil
		case float64:
			return NewBoolExpression(float64(l) > r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l > float64(r), location), nil
		case float64:
			return NewBoolExpression(l > r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type GreaterThanAndEqualExpression struct {
	baseBinaryExpression
}

func NewGreaterThanAndEqualExpression(left, right Expression) *GreaterThanAndEqualExpression {
	expression := &GreaterThanAndEqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *GreaterThanAndEqualExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "GRETER_THAN"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l >= r, location), nil
		case float64:
			return NewBoolExpression(float64(l) >= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l >= float64(r), location), nil
		case float64:
			return NewBoolExpression(l >= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type LessThanExpression struct {
	baseBinaryExpression
}

func NewLessThanExpression(left, right Expression) *LessThanExpression {
	expression := &LessThanExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LessThanExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "LESS_THAN"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l < r, location), nil
		case float64:
			return NewBoolExpression(float64(l) < r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l < float64(r), location), nil
		case float64:
			return NewBoolExpression(l < r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type LessThanAndEqualExpression struct {
	baseBinaryExpression
}

func NewLessThanAndEqualExpression(left, right Expression) *LessThanAndEqualExpression {
	expression := &LessThanAndEqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LessThanAndEqualExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	tag := "LESS_THAN"

	if left == nil || right == nil {
		return nil, errors.NewInvalidOperationError(tag, location, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l <= r, location), nil
		case float64:
			return NewBoolExpression(float64(l) <= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right).Name())
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l <= float64(r), location), nil
		case float64:
			return NewBoolExpression(l <= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type LogicalOrExpression struct {
	baseBinaryExpression
}

func NewLogicalOrExpression(left, right Expression) *LogicalOrExpression {
	expression := &LogicalOrExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LogicalOrExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	l, lok := left.(bool)
	r, rok := right.(bool)

	if lok && rok {
		return NewBoolExpression(l || r, location), nil
	} else {
		return nil, errors.NewInvalidOperationError("LOGICAL_AND", location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type LogicalAndExpression struct {
	baseBinaryExpression
}

func NewLogicalAndExpression(left, right Expression) *LogicalAndExpression {
	expression := &LogicalAndExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LogicalAndExpression) execute(left, right interface{},
	location *common.Location) (Expression, errors.Error) {
	l, lok := left.(bool)
	r, rok := right.(bool)

	if lok && rok {
		return NewBoolExpression(l && r, location), nil
	} else {
		return nil, errors.NewInvalidOperationError("LOGICAL_AND", location,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type unaryExpression struct {
	operand Expression

	// location of operator
	location *common.Location
}

type LogicalNotExpression struct {
	baseExpression
	unaryExpression
}

func NewLogicalNotExpression(operand Expression, location *common.Location) *LogicalNotExpression {
	expression := &LogicalNotExpression{
		unaryExpression: unaryExpression{
			operand:  operand,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LogicalNotExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	typ, err := expression.operand.getType(context)
	if err != nil {
		return expression, err
	}
	if typ != BOOL_TYPE {
		return nil, errors.NewInvalidOperationError("LOGICAL_OR", expression.getLocation(), typ.GetName())
	}

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return expression, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		return NewBoolExpression(!(expr.getValue().(bool)), expression.getLocation()), nil
	}

	return expression, nil
}

func (expression *LogicalNotExpression) getType(context *Context) (Type, errors.Error) {
	return BOOL_TYPE, nil
}

func (expression *LogicalNotExpression) getLocation() *common.Location {
	return expression.location
}

type MinusExpression struct {
	baseExpression
	unaryExpression
}

func NewMinusExpression(operand Expression, location *common.Location) *MinusExpression {
	expression := &MinusExpression{
		unaryExpression: unaryExpression{
			operand:  operand,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *MinusExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	typ, err := expression.operand.getType(context)
	if typ != INTEGER_TYPE && typ != FLOAT_TYPE {
		return nil, errors.NewInvalidOperationError("MINUS", expression.location, typ.GetName())
	}

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return expression, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		if typ == INTEGER_TYPE {
			return NewIntegerExpression(-(expr.getValue().(int64)), expression.getLocation()), nil
		} else {
			return NewFloatExpression(-(expr.getValue().(float64)), expression.getLocation()), nil
		}
	}

	return expression, nil
}

func (expression *MinusExpression) getType(context *Context) (Type, errors.Error) {
	return expression.operand.getType(context)
}

func (expression *MinusExpression) getLocation() *common.Location {
	return expression.location
}

type IncrementExpression struct {
	baseExpression
	unaryExpression
}

func NewIncrementExpression(operand Expression, location *common.Location) *IncrementExpression {
	expression := &IncrementExpression{
		unaryExpression: unaryExpression{
			operand:  operand,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *IncrementExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	typ, err := expression.operand.getType(context)
	if err != nil {
		return expression, err
	}
	if typ != INTEGER_TYPE {
		return nil, errors.NewInvalidOperationError("INCREMENT",
			expression.operand.getLocation(), typ.GetName())
	}

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return nil, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		return NewIntegerExpression(expr.getValue().(int64)+1, expression.getLocation()), nil
	}

	return expression, nil
}

func (expression *IncrementExpression) getType(context *Context) (Type, errors.Error) {
	return expression.operand.getType(context)
}

func (expression *IncrementExpression) getLocation() *common.Location {
	return expression.location
}

type DecrementExpression struct {
	baseExpression
	unaryExpression
}

func NewDecrementExpression(operand Expression, location *common.Location) *DecrementExpression {
	expression := &DecrementExpression{
		unaryExpression: unaryExpression{
			operand:  operand,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *DecrementExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	typ, err := expression.operand.getType(context)
	if err != nil {
		return expression, err
	}
	if typ != INTEGER_TYPE {
		return nil, errors.NewInvalidOperationError("DECREMENT",
			expression.operand.getLocation(), typ.GetName())
	}

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return nil, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		return NewIntegerExpression(expr.getValue().(int64)-1, expression.getLocation()), nil
	}

	return expression, nil
}

func (expression *DecrementExpression) getType(context *Context) (Type, errors.Error) {
	return expression.operand.getType(context)
}

func (expression *DecrementExpression) getLocation() *common.Location {
	return expression.location
}

type FunctionCallExpression struct {
	baseExpression
	identifier *Identifier
	arguments  []*Argument

	function *Function

	// use identifier's location as expression's location
}

func NewFunctionCallExpression(identifier *Identifier,
	arguments []*Argument) *FunctionCallExpression {
	expression := &FunctionCallExpression{
		identifier: identifier,
		arguments:  arguments,

		function: nil,
	}
	expression.this = expression
	return expression
}

func (expression *FunctionCallExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	if expression.function == nil {
		err = expression.searchFunction(context)
		if err != nil {
			return expression, nil
		}
	}

	args := expression.arguments
	params := expression.function.GetParameterList()

	if len(args) != len(params) {
		return nil, errors.NewArgumentCountMismatchError(len(params),
			len(args), expression.identifier.GetLocation())
	}

	for i, arg := range args {
		if err = arg.Fix(context); err != nil {
			return expression, err
		}
		if err = arg.CastTo(params[i].GetType(), context); err != nil {
			return expression, err
		}
	}

	return expression, nil
}

func (expression *FunctionCallExpression) getType(context *Context) (Type, errors.Error) {
	if expression.function == nil {
		err := expression.searchFunction(context)
		if err != nil {
			return nil, err
		}
	}
	return expression.function.GetType(), nil
}

func (expression *FunctionCallExpression) getLocation() *common.Location {
	return expression.identifier.GetLocation()
}

func (expression *FunctionCallExpression) searchFunction(context *Context) errors.Error {
	if expression.function == nil {
		expression.function = context.GetFunction(expression.identifier.GetName())
	}
	if expression.function == nil {
		return errors.NewFunctionNotFoundError(expression.identifier.GetName(), expression.getLocation())
	}
	return nil
}

type IndexExpression struct {
	baseExpression
	array Expression
	index Expression
}

func NewIndexExpression(array, index Expression) *IndexExpression {
	return &IndexExpression{
		array: array,
		index: index,
	}
}

func (expression *IndexExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	expression.array, err = expression.array.Fix(context)
	if err != nil {
		return expression, err
	}

	var typ Type

	typ, err = expression.array.getType(context)
	if !typ.IsDeriveType() {
		return expression, errors.NewInvalidTypeError(typ.GetName(), "array", expression.getLocation())
	}

	expression.index, err = expression.index.Fix(context)
	if err != nil {
		return expression, err
	}
	typ, err = expression.index.getType(context)
	if typ != INTEGER_TYPE {
		return expression, errors.NewIndexNotIntError(typ.GetName(), expression.index.getLocation())
	}

	return expression, nil
}

func (expression *IndexExpression) getType(context *Context) (Type, errors.Error) {
	typ, err := expression.array.getType(context)
	if err != nil {
		return nil, err
	}
	return typ.GetBaseType(), nil
}

func (expression *IndexExpression) getLocation() *common.Location {
	return expression.array.getLocation()
}

type castExpression struct {
	operand Expression
}

func (expression *castExpression) getLocation() *common.Location {
	return expression.operand.getLocation()
}

type IntegerToFloatCastExpression struct {
	baseExpression
	castExpression
}

func NewIntegerToFloatCastExpression(operand Expression) *IntegerToFloatCastExpression {
	expression := &IntegerToFloatCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.this = expression
	return expression
}

func (expression *IntegerToFloatCastExpression) getType(context *Context) (Type, errors.Error) {
	return FLOAT_TYPE, nil
}

type FloatToIntegerCastExpression struct {
	baseExpression
	castExpression
}

func NewFloatToIntegerCastExpression(operand Expression) *FloatToIntegerCastExpression {
	expression := &FloatToIntegerCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.this = expression
	return expression
}

func (expression *FloatToIntegerCastExpression) getType(context *Context) (Type, errors.Error) {
	return INTEGER_TYPE, nil
}

type NullToStringCastExpression struct {
	baseExpression
	castExpression
}

func NewNullToStringCastExpression(operand Expression) *NullToStringCastExpression {
	expression := &NullToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.this = expression
	return expression
}

func (expression *NullToStringCastExpression) getType(context *Context) (Type, errors.Error) {
	return STRING_TYPE, nil
}

type BoolToStringCastExpression struct {
	castExpression
}

func NewBoolToStringCastExpression(operand Expression) *NullToStringCastExpression {
	return &NullToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *BoolToStringCastExpression) getType() (Type, errors.Error) {
	return STRING_TYPE, nil
}

type IntegerToStringCastExpression struct {
	baseExpression
	castExpression
}

func NewIntegerToStringCastExpression(operand Expression) *IntegerToStringCastExpression {
	expression := &IntegerToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.this = expression
	return expression
}

func (expression *IntegerToStringCastExpression) getType(context *Context) (Type, errors.Error) {
	return STRING_TYPE, nil
}

type FloatToStringCastExpression struct {
	baseExpression
	castExpression
}

func NewFloatToStringCastExpression(operand Expression) *FloatToStringCastExpression {
	expression := &FloatToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.this = expression
	return expression
}

func (expression *FloatToStringCastExpression) getType(context *Context) (Type, errors.Error) {
	return STRING_TYPE, nil
}

type ArrayCreationExpression struct {
	baseExpression
	baseType   Type
	dimensions []Expression

	// use location of keyword 'new' as the whole expression's location
	location *common.Location

	typ Type
}

func NewArrayCreationExpression(baseType Type, dimensions []Expression,
	location *common.Location) *ArrayCreationExpression {
	expression := &ArrayCreationExpression{
		baseType:   baseType,
		dimensions: dimensions,

		location: location,

		typ: nil,
	}
	expression.this = expression
	return expression
}

func (expression *ArrayCreationExpression) Fix(context *Context) (Expression, errors.Error) {
	var err errors.Error

	for i, dim := range expression.dimensions {
		expression.dimensions[i], err = dim.Fix(context)
		if err != nil {
			return expression, err
		}
		typ, err := expression.dimensions[i].getType(context)
		if err != nil {
			return expression, err
		}
		if typ != INTEGER_TYPE {
			return expression, errors.NewArraySizeNotIntError(typ.GetName(),
				expression.dimensions[i].getLocation())
		}
	}
	if expression.typ == nil {
		expression.fetchType()
	}

	return expression, nil
}

func (expression *ArrayCreationExpression) getType(context *Context) (Type, errors.Error) {
	if expression.typ == nil {
		expression.fetchType()
	}
	return expression.typ, nil
}

func (expression *ArrayCreationExpression) getLocation() *common.Location {
	return expression.location
}

func (expression *ArrayCreationExpression) fetchType() {
	deriveTags := []DeriveTag{}
	for i := 0; i < len(expression.dimensions); i++ {
		deriveTags = append(deriveTags, NewArrayDerive())
	}
	expression.typ = NewDeriveType(expression.baseType, deriveTags)
}

//
// utility functions
//

func typeCast(srcType, destType Type, operand Expression) (Expression, errors.Error) {
	switch srcType {
	case NULL_TYPE:
		return nullTypeCast(destType, operand)
	case BOOL_TYPE:
		return boolTypeCast(destType, operand)
	case INTEGER_TYPE:
		return integerTypeCast(destType, operand)
	case FLOAT_TYPE:
		return floatTypeCast(destType, operand)
	case STRING_TYPE:
		return stringTypeCast(destType, operand)
	default:
		return nil, errors.NewUnsupportedTypeError(srcType.GetName(), nil)
	}
}

func nullTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType.Equal(NULL_TYPE) {
		return operand, nil
	} else if destType.Equal(STRING_TYPE) {
		return NewNullToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(NULL_TYPE.GetName(), destType.GetName(), operand.getLocation())
	}
}

func boolTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType.Equal(BOOL_TYPE) {
		return operand, nil
	} else if destType.Equal(STRING_TYPE) {
		return NewBoolToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(BOOL_TYPE.GetName(), destType.GetName(), operand.getLocation())
	}
}

func integerTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType.Equal(INTEGER_TYPE) {
		return operand, nil
	} else if destType.Equal(FLOAT_TYPE) {
		return NewIntegerToFloatCastExpression(operand), nil
	} else if destType.Equal(STRING_TYPE) {
		return NewIntegerToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(INTEGER_TYPE.GetName(), destType.GetName(), operand.getLocation())
	}
}

func floatTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType.Equal(INTEGER_TYPE) {
		return NewFloatToIntegerCastExpression(operand), nil
	} else if destType.Equal(FLOAT_TYPE) {
		return operand, nil
	} else if destType.Equal(STRING_TYPE) {
		return NewFloatToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(FLOAT_TYPE.GetName(), destType.GetName(), operand.getLocation())
	}
}

func stringTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	return nil, errors.NewTypeCastError(STRING_TYPE.GetName(), destType.GetName(), operand.getLocation())
}
