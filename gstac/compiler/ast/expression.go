package ast

import (
	"reflect"
	"strings"

	"fmt"
	"github.com/mlmhl/compiler/common"
	"github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/token"
	"strconv"
)

type Expression interface {
	Fix(context *Context) (Expression, errors.Error)
	CastTo(destType Type, context *Context) (Expression, errors.Error)
	Generate(context *Context)

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
	return expression, nil
}

func (expression *baseExpression) CastTo(destType Type, context *Context) (Expression, errors.Error) {
	srcType, err := expression.this.getType(context)
	if err != nil {
		return expression, err
	}
	return typeCast(srcType, destType, expression)
}

func (expression *baseExpression) Generate(context *Context) {
	// NO-OP
}

func (expression *baseExpression) getType(context *Context) (Type, errors.Error) {
	panic("Can't invoke `getType` on " + reflect.TypeOf(expression).Elem().Name())
}

//
// value expression
//

type valueExpression interface {
	getValue() interface{}
}

type baseValueExpression struct {
	baseExpression
	typ      Type
	location *common.Location
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

type IntegerExpression struct {
	baseValueExpression
	value int
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
	return expression
}

func (expression *StringExpression) getValue() interface{} {
	return expression.value
}

// Literal array like int[] a = `{1, 2, 3}`
type ArrayLiteralExpression struct {
	baseExpression
	values []Expression

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

func (expression *ArrayLiteralExpression) CastTo(destType Type, context *Context) Expression {
	var err errors.Error
	for i, exp := range expression.values {
		expression.values[i], err = exp.CastTo(destType, context)
		if err != nil {
			return expression, err
		}
	}
	return expression
}

func (expression *ArrayLiteralExpression) getLocation() *common.Location {
	return expression.location
}

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

func (expression *IdentifierExpression) getIdentifier() *Identifier {
	return expression.identifier
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
	assignExpression := assignExpression{
		left:    left,
		operand: operand,
	}
	var expression *assignExpression

	switch typ {
	case token.ASSIGN_ID:
		expression = &NormalAssignExpression{
			assignExpression: assignExpression,
		}
	case token.ADD_ASSIGN_ID:
		expression = &AddAssignExpression{
			assignExpression: assignExpression,
		}
	case token.SUBTRACT_ID:
		expression = &SubtractAssignExpression{
			assignExpression: assignExpression,
		}
	case token.MUL_ASSIGN_ID:
		expression = &MultiplyAssignExpression{
			assignExpression: assignExpression,
		}
	case token.DIV_ASSIGN_ID:
		expression = &DivideAssignExpression{
			assignExpression: assignExpression,
		}
	case token.MOD_ASSIGN_ID:
		expression = &ModAssignExpression{
			assignExpression: assignExpression,
		}
	default:
		expression = nil
	}

	if expression != nil {
		expression.this = expression
	}

	return expression
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
		expression.right, err = expression.right.CastTo(expression.left.getType(context))
	} else {
		expression.left, err = expression.left.CastTo(expression.right.getType(context))
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
		return expression, err
	}

	rightType, err := expression.right.getType(context)
	if err != nil {
		return expression, err
	}

	if leftType.isPriorityOf(rightType) {
		return leftType
	} else {
		return rightType
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
			return NewStringExpression("nil"+r, location), nil
		} else {
			return nil, errors.NewInvalidOperationError(tag, location, "null", reflect.TypeOf(right).Name())
		}
	}

	switch l := left.(type) {
	case bool:
		switch r := right.(type) {
		case string:
			if l {
				return NewStringExpression("true"+r, location), nil
			} else {
				return NewStringExpression("false"+r, location), nil
			}
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "bool", reflect.TypeOf(right).Name(), location)
		}
	case int64:
		switch r := right.(type) {
		case int64:
			return NewIntegerExpression(l+r, location), nil
		case float64:
			return NewFloatExpression(float64(l)+r, location), nil
		case string:
			return NewStringExpression(strconv.Itoa(int(l))+r, location), nil
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
			return NewStringExpression(fmt.Sprintf("%f", l)+r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right).Name())
		}
	case string:
		switch r := right.(type) {
		case bool:
			if r {
				return NewStringExpression(l+"true", location), nil
			} else {
				return NewStringExpression(l+"false", location), nil
			}
		case int64:
			return NewStringExpression(l+strconv.Itoa(int(r)), location), nil
		case float64:
			return NewStringExpression(l+fmt.Sprintf("%f", r), location), nil
		case string:
			return NewStringExpression(l+r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "string", reflect.TypeOf(right).Name())
		}
	default:
		return nil, errors.NewInvalidOperationError(tag,
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
		return nil, errors.NewInvalidOperationError(tag,
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
		return nil, errors.NewInvalidOperationError(tag,
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
		return nil, errors.NewInvalidOperationError(tag,
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
		return nil, errors.NewInvalidOperationError(tag,
			reflect.TypeOf(left).Name(), reflect.TypeOf(right).Name())
	}
}

type EqualExpression struct {
	baseBinaryExpression
}

func NewEqualExpression(left, right Expression) *EqualExpression {
	expression := EqualExpression{
		baseBinaryExpression: binaryExpression{
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
	expression := NotEqualExpression{
		baseBinaryExpression: binaryExpression{
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
	expression := GreaterThanExpression{
		baseBinaryExpression: binaryExpression{
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
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l > float64(r), location), nil
		case float64:
			return NewBoolExpression(l > r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right))
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left), reflect.TypeOf(right))
	}
}

type GreaterThanAndEqualExpression struct {
	baseBinaryExpression
}

func NewGreaterThanAndEqualExpression(left, right Expression) *GreaterThanAndEqualExpression {
	expression := GreaterThanAndEqualExpression{
		baseBinaryExpression: binaryExpression{
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
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l >= float64(r), location), nil
		case float64:
			return NewBoolExpression(l >= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right))
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left), reflect.TypeOf(right))
	}
}

type LessThanExpression struct {
	baseBinaryExpression
}

func NewLessThanExpression(left, right Expression) *LessThanExpression {
	expression := LessThanExpression{
		baseBinaryExpression: binaryExpression{
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
		return errors.NewInvalidOperationError(tag, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l+r, location)
		case float64:
			return NewBoolExpression(float64(l)+r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l+float64(r), location)
		case float64:
			return NewBoolExpression(l+r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "float", reflect.TypeOf(right))
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left), reflect.TypeOf(right))
	}
}

type LessThanAndEqualExpression struct {
	baseBinaryExpression
}

func NewLessThanAndEqualExpression(left, right Expression) *LessThanAndEqualExpression {
	expression := LessThanAndEqualExpression{
		baseBinaryExpression: binaryExpression{
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
		return nil, errors.NewInvalidOperationError(tag, "null")
	}

	switch l := left.(type) {
	case int64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l <= r, location)
		case float64:
			return NewBoolExpression(float64(l) <= r, location)
		default:
			return nil, errors.NewInvalidOperationError(tag, location, "int", reflect.TypeOf(right))
		}
	case float64:
		switch r := right.(type) {
		case int64:
			return NewBoolExpression(l <= float64(r), location), nil
		case float64:
			return NewBoolExpression(l <= r, location), nil
		default:
			return nil, errors.NewInvalidOperationError(tag, "float", reflect.TypeOf(right))
		}
	default:
		return nil, errors.NewInvalidOperationError(tag, location, reflect.TypeOf(left), reflect.TypeOf(right))
	}
}

type LogicalOrExpression struct {
	baseBinaryExpression
}

func NewLogicalOrExpression(left, right Expression) *LogicalOrExpression {
	expression := &LogicalOrExpression{
		baseBinaryExpression: binaryExpression{
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
		return nil, errors.NewInvalidOperationError("LOGICAL_AND", reflect.TypeOf(left), reflect.TypeOf(right))
	}
}

type LogicalAndExpression struct {
	baseBinaryExpression
}

func NewLogicalAndExpression(left, right Expression) *LogicalAndExpression {
	expression := &LogicalAndExpression{
		baseBinaryExpression: binaryExpression{
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
		return nil, errors.NewInvalidOperationError("LOGICAL_AND", reflect.TypeOf(left), reflect.TypeOf(right))
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

func NewLogicalNotExpression(operand Expression, location *common.Location) *LogicalOrExpression {
	expression := &LogicalOrExpression{
		unaryExpression{
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
		return nil, errors.NewInvalidOperationError("MINUS", typ.GetName())
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

func (expression *DecrementExpression) getType(context *Context) Type {
	return expression.operand.getType(context)
}

func (expression *DecrementExpression) getLocation() *common.Location {
	return expression.location
}

type FunctionCallExpression struct {
	identifier *Identifier
	arguments  []*Argument

	function *Function

	// use identifier's location as expression's location
}

func NewFunctionCallExpression(identifier *Identifier,
	arguments []*Argument) *FunctionCallExpression {
	return &FunctionCallExpression{
		identifier: identifier,
		arguments:  arguments,

		function: nil,
	}
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
		if err = arg.CastTo(params[i].GetType()); err != nil {
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
	return expression.function.GetType().GetName(), nil
}

func (expression *FunctionCallExpression) GetLocation() *common.Location {
	return expression.identifier.GetLocation()
}

func (expression *FunctionCallExpression) searchFunction(context *Context) errors.Error {
	if expression.function == nil {
		expression.function = context.GetFunction(expression.identifier.GetName())
	}
	if expression.function == nil {
		return errors.NewFunctionNotFoundError(expression.identifier.GetName(), expression.GetLocation())
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

	typ, err = expression.array.getType()
	if !typ.IsDeriveType() {
		return expression, errors.NewInvalidTypeError(typ.GetName(), "array", expression.getLocation())
	}

	expression.index, err = expression.index.Fix(context)
	if err != nil {
		return expression, err
	}
	typ, err = expression.index.getType()
	if typ != INTEGER_TYPE {
		return expression, errors.NewIndexNotIntError(typ.GetName(), expression.index.getLocation())
	}

	return expression, nil
}

func (expression *IndexExpression) getType() (Type, errors.Error) {
	typ, err := expression.array.getType()
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
	return castExpression.operand.getLocation()
}

type IntegerToFloatCastExpression struct {
	castExpression
}

func NewIntegerToFloatCastExpression(operand Expression) *IntegerToFloatCastExpression {
	return &IntegerToFloatCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *IntegerToFloatCastExpression) getType() (Type, errors.Error) {
	return FLOAT_TYPE, nil
}

type FloatToIntegerCastExpression struct {
	castExpression
}

func NewFloatToIntegerCastExpression(operand Expression) *FloatToIntegerCastExpression {
	return &FloatToIntegerCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *FloatToIntegerCastExpression) getType() (Type, errors.Error) {
	return INTEGER_TYPE, nil
}

type NullToStringCastExpression struct {
	castExpression
}

func NewNullToStringCastExpression(operand Expression) *NullToStringCastExpression {
	return &NullToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *NullToStringCastExpression) getType() (Type, errors.Error) {
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
	castExpression
}

func NewIntegerToStringCastExpression(operand Expression) *IntegerToStringCastExpression {
	return &IntegerToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *IntegerToStringCastExpression) getType() (Type, errors.Error) {
	return STRING_TYPE, nil
}

type FloatToStringCastExpression struct {
	castExpression
}

func NewFloatToStringCastExpression(operand Expression) *FloatToStringCastExpression {
	return &FloatToStringCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
}

func (expression *FloatToStringCastExpression) getType() (Type, errors.Error) {
	return STRING_TYPE, nil
}

type ArrayCreationExpression struct {
	baseExpression
	baseType   Type
	dimensions []Expression

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
		return nil, errors.NewTypeCastError(NULL_TYPE, destType, operand.getLocation())
	}
}

func boolTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType.Equal(BOOL_TYPE) {
		return operand, nil
	} else if destType.Equal(STRING_TYPE) {
		return NewBoolToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(BOOL_TYPE, destType, operand.getLocation())
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
		return nil, errors.NewTypeCastError(INTEGER_TYPE, destType, operand.getLocation())
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
		return nil, errors.NewTypeCastError(FLOAT_TYPE, destType, operand.getLocation())
	}
}

func stringTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	return nil, errors.NewTypeCastError(STRING_TYPE, destType, operand.getLocation())
}
