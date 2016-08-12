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
	CastTo(destType Type) (Expression, errors.Error)
	Generate(context *Context)

	getType() Type
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

func (expression *baseExpression) CastTo(destType Type) (Expression, errors.Error) {
	// NO-OP
	return typeCast(expression.this.getType(), destType, expression)
}

func (expression *baseExpression) Generate(context *Context) {
	// NO-OP
}

func (expression *baseExpression) getType() Type {
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

func (expression *baseValueExpression) getType() Type {
	return expression.typ
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

func (expression *ArrayLiteralExpression) CastTo(destType Type) Expression {
	for i, exp := range expression.values {
		expression.values[i] = exp.CastTo(destType)
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

func (expression *assignExpression) Fix(context *Context) {
	var err errors.Error
	expression.left, err = expression.left.Fix(context)
	if err != nil {
		return expression, err
	}
	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return expression, err
	}
	expression.operand, err = expression.operand.CastTo(expression.left.getType())
	return expression, err
}

func (expression *assignExpression) getType() Type {
	return expression.left.getType()
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

	if expression.left.getType().isPriorityOf(expression.right.getType()) {
		expression.right, err = expression.right.CastTo(expression.left.getType())
	} else {
		expression.left, err = expression.left.CastTo(expression.right.getType())
	}
	return expression, err
}

func (expression *baseBinaryExpression) CastTo(destType Type) (Expression, errors.Error) {
	srcType := expression.left.getType()
	if expression.right.getType().isPriorityOf(srcType) {
		srcType = expression.right.getType()
	}
	return typeCast(srcType, destType, expression)
}

func (expression *baseBinaryExpression) getType() Type {
	if expression.left.getType().isPriorityOf(expression.right.getType()) {
		return expression.left.getType()
	} else {
		return expression.right.getType()
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
	typ := expression.operand.getType()
	if typ != BOOL_TYPE {
		return nil, errors.NewInvalidOperationError("LOGICAL_OR", expression.getLocation(), typ.GetName())
	}

	var err errors.Error

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

func (expression *LogicalNotExpression) getType() Type {
	return BOOL_TYPE
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
	typ := expression.operand.getType()
	if typ != INTEGER_TYPE && typ != FLOAT_TYPE {
		return nil, errors.NewInvalidOperationError("MINUS", typ.GetName())
	}

	var err errors.Error

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

func (expression *MinusExpression) getType() Type {
	return expression.operand.getType()
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
	typ := expression.operand.getType()
	if typ != INTEGER_TYPE {
		return nil, errors.NewInvalidOperationError("INCREMENT",
			expression.operand.getLocation(), typ.GetName())
	}

	var err errors.Error

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return nil, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		return NewIntegerExpression(expr.getValue().(int64) + 1, expression.getLocation()), nil
	}

	return expression, nil
}

func (expression *IncrementExpression) getType() Type {
	return expression.operand.getType()
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
			operand: operand,
			location: location,
		},
	}
	expression.this = expression
	return expression
}

func (expression *DecrementExpression) Fix(context *Context) (Expression, errors.Error) {
	typ := expression.operand.getType()
	if typ != INTEGER_TYPE {
		return nil, errors.NewInvalidOperationError("DECREMENT",
			expression.operand.getLocation(), typ.GetName())
	}

	var err errors.Error

	expression.operand, err = expression.operand.Fix(context)
	if err != nil {
		return nil, err
	}

	expr, ok := expression.operand.(valueExpression)
	if ok {
		return NewIntegerExpression(expr.getValue().(int64) - 1, expression.getLocation()), nil
	}

	return expression, nil
}

func (expression *DecrementExpression) getType() Type {
	return expression.operand.getType()
}

func (expression *DecrementExpression) getLocation() *common.Location {
	return expression.location
}

type FunctionCallExpression struct {
	identifier *Identifier
	arguments  []*Argument
}

func NewFunctionCallExpression(identifier *Identifier,
	arguments []*Argument) *FunctionCallExpression {
	return &FunctionCallExpression{
		identifier: identifier,
		arguments:  arguments,
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

func NewIntegerToFloatCastExpression(operand Expression) *IntegerToFloatCastExpression {
	return &IntegerToFloatCastExpression{
		castExpression: castExpression{
			operand: operand,
		},
	}
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

type ArrayCreationExpression struct {
	typ        Type
	dimensions []Expression
}

func NewArrayCreationExpression(typ Type, dimensions []Expression) *ArrayCreationExpression {
	return &ArrayCreationExpression{
		typ:        typ,
		dimensions: dimensions,
	}
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
	if destType == NULL_TYPE {
		return operand, nil
	} else if destType == STRING_TYPE {
		return NewNullToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(NULL_TYPE, destType, operand.getLocation())
	}
}

func boolTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType == BOOL_TYPE {
		return operand, nil
	} else if destType == STRING_TYPE {
		return NewBoolToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(BOOL_TYPE, destType, operand.getLocation())
	}
}

func integerTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType == INTEGER_TYPE {
		return operand, nil
	} else if destType == FLOAT_TYPE {
		return NewIntegerToFloatCastExpression(operand), nil
	} else if destType == STRING_TYPE {
		return NewIntegerToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(INTEGER_TYPE, destType, operand.getLocation())
	}
}

func floatTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	if destType == INTEGER_TYPE {
		return NewFloatToIntegerCastExpression(operand), nil
	} else if destType == FLOAT_TYPE {
		return operand, nil
	} else if destType == STRING_TYPE {
		return NewFloatToStringCastExpression(operand), nil
	} else {
		return nil, errors.NewTypeCastError(FLOAT_TYPE, destType, operand.getLocation())
	}
}

func stringTypeCast(destType Type, operand Expression) (Expression, errors.Error) {
	return nil, errors.NewTypeCastError(STRING_TYPE, destType, operand.getLocation())
}
