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
	Generate(context *Context, exe *executable.Executable) errors.Error

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

func (expression *baseExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	// NO-OP
	return nil
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
	getOperatorCode() byte

	// return the value's code byte
	valueEncode(exe *executable.Executable) []byte
}

type baseValueExpression struct {
	baseExpression
	typ      Type
	location *common.Location
}

func (expression *baseValueExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	valueExpression := expression.this.(valueExpression)

	exe.AppendSlice(expression.location.Encode())
	exe.Append(valueExpression.getOperatorCode())
	exe.AppendSlice(valueExpression.valueEncode(exe))

	return nil
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

func (expression *FloatExpression) operatorCode() byte {
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
func (expression *ArrayLiteralExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	// Generate expression's location
	exe.AppendSlice(expression.location.Encode())

	// Generate code byte of array elements
	for _, subExpr := range expression.values {
		err = subExpr.Generate(context, exe)
		if err != nil {
			return err
		}
	}

	// Generate code byte of ArrayLiteralExpression

	exe.AppendSlice(expression.location.Encode())
	// Generate is called after CastTo, so typ is already set to the correct value
	exe.Append(executable.GetOperatorCode(
		executable.NEW_ARRAY_LITERAL_BOOL, expression.typ.GetOffset()))

	return nil
}

// If typ is not set, getType will return a nil, but this doesn't matter,
// as getType won't be called up to now.
func (expression *ArrayLiteralExpression) getType(context *Context) (Type, errors.Error) {
	return expression.typ, nil
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

func (expression *IdentifierExpression) Generate(context *Context,
	exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.getLocation().Encode())
	exe.Append(executable.VARIABLE_REFERENCE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(
		context.GetSymbolIndex(expression.identifier.GetName())))

	return nil
}

func (expression *IdentifierExpression) getType(context *Context) (Type, errors.Error) {
	return context.GetVariable(expression.identifier.GetName()).GetType(), nil
}

func (expression *IdentifierExpression) getLocation() *common.Location {
	return expression.identifier.GetLocation()
}

type assignExpressionInterface interface {
	// If return true, push left operand's value to stack
	isLeftNeedPush() bool

	// Return the operator code according to content assign expression
	getOperatorCode(typ Type) byte
}

type baseAssignExpression struct {
	baseExpression
	left    Expression
	operand Expression

	// use left expression's location as the whole expression's location
}

func (expression *baseAssignExpression) isLeftNeedPush() bool {
	return true
}

type NormalAssignExpression struct {
	baseAssignExpression
}

func (expression *NormalAssignExpression) isLeftNeedPush() bool {
	return false
}

func (expression *NormalAssignExpression) getOperatorCode(context *Context) (byte, errors.Error) {
	return executable.NORMAL_ASSIGN, nil
}

type AddAssignExpression struct {
	baseAssignExpression
}

func (expression *AddAssignExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.ADD_BOOL, typ.GetOffset())
}

type SubtractAssignExpression struct {
	baseAssignExpression
}

func (expression *SubtractAssignExpression) getOperatorCode(typ Type) byte {
	// Add operation doesn't support 'bool', so 'int' is the start type.
	return executable.GetOperatorCode(executable.SUBTRACT_INT, typ.GetOffset()-1)
}

type MultiplyAssignExpression struct {
	baseAssignExpression
}

func (expression *MultiplyAssignExpression) getOperatorCode(typ Type) byte {
	// Add operation doesn't support 'bool', so 'int' is the start type.
	return executable.GetOperatorCode(executable.MULTIPLY_INT, typ.GetOffset()-1)
}

type DivideAssignExpression struct {
	baseAssignExpression
}

func (expression *DivideAssignExpression) getOperatorCode(typ Type) byte {
	// Add operation doesn't support and 'bool', so 'int' is the start type.
	return executable.GetOperatorCode(executable.DIVIDE_INT, typ.GetOffset()-1)
}

type ModAssignExpression struct {
	baseAssignExpression
}

func (expression *ModAssignExpression) getOperatorCode(typ Type) byte {
	// Add operation doesn't support 'null' and 'bool', so 'int' is the start type.
	return executable.MOD_INT
}

func NewAssignExpression(typ int, left, operand Expression) Expression {
	expression := baseAssignExpression{
		left:    left,
		operand: operand,
	}
	switch typ {
	case token.ASSIGN_ID:
		result := &NormalAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	case token.ADD_ASSIGN_ID:
		result := &AddAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	case token.SUBTRACT_ID:
		result := &SubtractAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	case token.MUL_ASSIGN_ID:
		result := &MultiplyAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	case token.DIV_ASSIGN_ID:
		result := &DivideAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	case token.MOD_ASSIGN_ID:
		result := &ModAssignExpression{
			baseAssignExpression: expression,
		}
		result.this = result
		return result
	default:
		return nil
	}
}

func (expression *baseAssignExpression) Fix(context *Context) (Expression, errors.Error) {
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

func (expression *baseAssignExpression) Generate(context *Context,
	exe *executable.Executable) errors.Error {
	var err errors.Error

	// Generate expression's location
	exe.AppendSlice(expression.getLocation().Encode())

	expr := expression.this.(assignExpressionInterface)

	// put left operand's code if needed
	if expr.isLeftNeedPush() {
		err = expression.left.Generate(context, exe)
		if err != nil {
			return err
		}
	}

	// put right operand's code
	err = expression.operand.Generate(context, exe)
	if err != nil {
		return err
	}

	// write this assign expression
	exe.AppendSlice(expression.getLocation().Encode())
	typ, err := expression.getType(context)
	if err != nil {
		return err
	}
	exe.Append(expr.getOperatorCode(typ))

	// support for "a=b=c"
	exe.Append(executable.STACK_TOP_DUPLICATE)

	return popToLeftValue(expression.left, context, exe)
}

func (expression *baseAssignExpression) getType(context *Context) (Type, errors.Error) {
	return expression.left.getType(context)
}

func (expression *baseAssignExpression) getLocation() *common.Location {
	return expression.left.getLocation()
}

type binaryExpression interface {
	execute(left, right interface{}, location *common.Location) (Expression, errors.Error)
	getOperatorCode(typ Type) byte
}

type baseBinaryExpression struct {
	this  Expression
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
		return expression.this.(binaryExpression).execute(l.getValue(), r.getValue(), expression.left.getLocation())
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

func (expression *baseBinaryExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	var (
		typ  Type
		err  errors.Error
	)

	exe.AppendSlice(expression.getLocation().Encode())
	if err = expression.left.Generate(context, exe); err != nil {
		return err
	}
	if err = expression.right.Generate(context, exe); err != nil {
		return err
	}

	// Generate the operator
	if typ, err = expression.this.getType(context); err != nil {
		return err
	} else {
		exe.Append(expression.this.(binaryExpression).getOperatorCode(typ))
	}

	return nil
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

func (expression *AddExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.ADD_BOOL, typ.GetOffset())
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

func (expression *SubtractExpression) getOperatorCode(typ Type) byte {
	// For subtract operator, Int is the start type
	return executable.GetOperatorCode(executable.SUBTRACT_INT, typ.GetOffset()-1)
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

func (expression *MultiplyExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.MULTIPLY_INT, typ.GetOffset()-1)
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

func (expression *DivideExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.DIVIDE_INT, typ.GetOffset()-1)
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

func (expression *ModExpression) getOperatorCode(typ Type) byte {
	return executable.MOD_INT
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

func (expression *EqualExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.EQUAL_BOOL, typ.GetOffset())
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

func (expression *NotEqualExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.NOT_EQUAL_BOOL, typ.GetOffset())
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

func (expression *GreaterThanExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.GREATER_THAN_INT, typ.GetOffset()-1)
}

type GreaterThanOrEqualExpression struct {
	baseBinaryExpression
}

func NewGreaterThanAndEqualExpression(left, right Expression) *GreaterThanOrEqualExpression {
	expression := &GreaterThanOrEqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *GreaterThanOrEqualExpression) execute(left, right interface{},
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

func (expression *GreaterThanOrEqualExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.GREATER_THAN_OR_EQUAL_INT, typ.GetOffset()-1)
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

func (expression *LessThanExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.LESS_THAN_INT, typ.GetOffset()-1)
}

type LessThanOrEqualExpression struct {
	baseBinaryExpression
}

func NewLessThanAndEqualExpression(left, right Expression) *LessThanOrEqualExpression {
	expression := &LessThanOrEqualExpression{
		baseBinaryExpression: baseBinaryExpression{
			left:  left,
			right: right,
		},
	}
	expression.this = expression
	return expression
}

func (expression *LessThanOrEqualExpression) execute(left, right interface{},
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

func (expression *LessThanOrEqualExpression) getOperatorCode(typ Type) byte {
	return executable.GetOperatorCode(executable.LESS_THAN_OR_EQUAL_INT, typ.GetOffset()-1)
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

func (expression *LogicalOrExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(expression.getLocation().Encode())

	if err = expression.left.Generate(context, exe); err != nil {
		return err
	}

	label := exe.NewLabel()
	exe.Append(executable.STACK_TOP_DUPLICATE)
	jumpStatementCodeByte(executable.JUMP_IF_TRUE, label, exe)

	if err = expression.right.Generate(context, exe); err != nil {
		return err
	}

	exe.Append(executable.LOGICAL_OR)

	exe.SetLabel(label, exe.GetSize())

	return nil
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

func (expression *LogicalAndExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	var err errors.Error

	exe.AppendSlice(expression.getLocation().Encode())

	if err = expression.left.Generate(context, exe); err != nil {
		return err
	}

	// Use two level jump to support 'if statement jump'.
	// If left == false, first jump to the beginning of 'then condition',
	// and then jump to the beginning of 'else condition'. So we need to
	// duplicate the stack top value.
	label := exe.NewLabel()
	exe.Append(executable.STACK_TOP_DUPLICATE)
	jumpStatementCodeByte(executable.JUMP_IF_FALSE, label, exe)

	if err = expression.right.Generate(context, exe); err != nil {
		return err
	}

	exe.Append(executable.LOGICAL_AND)

	// use code byte offset as address
	exe.SetLabel(label, exe.GetSize())

	return nil
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

func (expression *LogicalNotExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.location.Encode())

	if err := expression.operand.Generate(context, exe); err != nil {
		return err
	}

	exe.Append(executable.LOGICAL_NOT)

	return nil
}

func (expression *LogicalNotExpression) getType(context *Context) (Type, errors.Error) {
	return BOOL_TYPE, nil
}

func (expression *LogicalNotExpression) getLocation() *common.Location {
	return expression.location
}

func (expression *LogicalNotExpression) getOperatorCode(typ Type) byte {
	return executable.LOGICAL_NOT
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

func (expression *MinusExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.location.Encode())

	if err := expression.operand.Generate(context, exe); err != nil {
		return err
	}

	if typ, err := expression.getType(context); err != nil {
		return err
	} else {
		exe.Append(executable.GetOperatorCode(executable.MINUS_INT, typ.GetOffset()-1))
	}

	return nil
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

func (expression *IncrementExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.location.Encode())
	if err := expression.operand.Generate(context, exe); err != nil {
		return err
	}
	// Should increment the value first, and then duplicated it.
	exe.Append(executable.INCREMENT)
	if !context.IsGlobal() {
		exe.Append(executable.STACK_TOP_DUPLICATE)
	}
	return popToLeftValue(expression.operand, context, exe)
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

func (expression *DecrementExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.location.Encode())
	if err := expression.operand.Generate(context, exe); err != nil {
		return err
	}
	// Should increment the value first, and then duplicated it.
	exe.Append(executable.DECREMENT)
	if !context.IsGlobal() {
		exe.Append(executable.STACK_TOP_DUPLICATE)
	}
	return popToLeftValue(expression.operand, context, exe)
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

func (expression *FunctionCallExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.getLocation().Encode())

	for _, arg := range(expression.arguments) {
		if err := arg.Generate(context, exe); err != nil {
			return err
		}
	}

	// Argument count needn't to save, as we can know it from function's definition.
	exe.Append(executable.FUNCTION_INVOKE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(context.GetSymbolIndex(expression.identifier.GetName())))

	return nil
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

func (expression *IndexExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.getLocation().Encode())
	if err := expression.array.Generate(context, exe); err != nil {
		return err
	}
	if err := expression.index.Generate(context, exe); err != nil {
		return err
	}

	if typ, err := expression.getType(context); err != nil {
		return err
	} else {
		exe.Append(executable.GetOperatorCode(executable.ARRAY_INDEX_BOOL, typ.GetOffset()))
		return nil
	}
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

type castExpressionInterface interface {
	getOperatorCode() byte
}

type castExpression struct {
	baseExpression
	operand Expression
}

func (expression *castExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.getLocation().Encode())
	if err := expression.operand.Generate(context, exe); err != nil {
		return err
	}
	exe.Append(expression.this.(castExpressionInterface).getOperatorCode())
	return nil
}

func (expression *castExpression) getLocation() *common.Location {
	return expression.operand.getLocation()
}

type IntegerToFloatCastExpression struct {
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

func (expression *IntegerToFloatCastExpression) getOperatorCode() byte {
	return executable.INT_TO_FLOAT
}

type FloatToIntegerCastExpression struct {
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

func (expression *FloatToIntegerCastExpression) getOperatorCode() byte {
	return executable.FLOAT_TO_INT
}

type NullToStringCastExpression struct {
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

func (expression *NullToStringCastExpression) getOperatorCode() byte {
	return executable.NULL_TO_STRING
}

type BoolToStringCastExpression struct {
	castExpression
}

func NewBoolToStringCastExpression(operand Expression) *NullToStringCastExpression {
	expression := &NullToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
	expression.castExpression.this = expression
	return expression
}

func (expression *BoolToStringCastExpression) getType() (Type, errors.Error) {
	return STRING_TYPE, nil
}

func (expression *BoolToStringCastExpression) getOperatorCode() byte {
	return executable.BOOL_TO_STRING
}

type IntegerToStringCastExpression struct {
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

func (expression *IntegerToStringCastExpression) getOperatorCode() byte {
	return executable.INT_TO_STRING
}

type FloatToStringCastExpression struct {
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

func (expression *FloatToStringCastExpression) getOperatorCode() byte {
	return executable.FLOAT_TO_STRING
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

func (expression *ArrayCreationExpression) Generate(context *Context, exe *executable.Executable) errors.Error {
	exe.AppendSlice(expression.location.Encode())

	for _, dim := range(expression.dimensions) {
		if err := dim.Generate(context, exe); err != nil {
			return nil
		}
	}

	exe.Append(executable.ARRAY_CREATE)
	exe.AppendSlice(encoding.DefaultEncoder.Int(len(expression.dimensions)))

	return nil
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

// Pop stack top value into a left value
func popToLeftValue(left Expression, context *Context,
exe *executable.Executable) errors.Error {

	var err errors.Error

	switch expr := left.(type) {
	case *IdentifierExpression:
		err = popToIdentifier(expr, context, exe)
	case *IndexExpression:
		err = popToArrayIndex(expr, context, exe)
	}

	return err
}

// Pop stack top value into a variable
func popToIdentifier(expression *IdentifierExpression, context *Context,
	exe *executable.Executable) errors.Error {
	typ, err := expression.getType(context)
	if err != nil {
		return err
	}
	var start byte
	if context.IsGlobal() {
		start = executable.POP_STATIC_BOOL
	} else {
		start = executable.POP_STACK_BOOL
	}
	// Int is the first supported type for pop stack operation
	exe.Append(executable.GetOperatorCode(start, typ.GetOffset()-1))
	return nil
}

// Pop stack top value into array index
func popToArrayIndex(expression *IndexExpression, context *Context,
	exe *executable.Executable) errors.Error {
	var err errors.Error

	if err = expression.array.Generate(context, exe); err != nil {
		return err
	}

	if err = expression.index.Generate(context, exe); err != nil {
		return err
	}

	exe.Append(executable.POP_ARRAY_BOOL)

	return nil
}

// Generate  ajump statement's code byte
func jumpStatementCodeByte(code byte, label int, exe *executable.Executable) {
	exe.Append(code)
	exe.AppendSlice(encoding.DefaultEncoder.Int(label))
}