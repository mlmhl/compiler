package types

import (
	"fmt"

	gerror "github.com/mlmhl/compiler/gdync/errors"
	"github.com/mlmhl/goutil/math"
)

//
// value type flyweight
//

const (
	STRING_TYPE = stringType("String")
	INTEGER_TYPE = integerType("Integer")
	FLOAT_TYPE = floatType("Float")
	BOOL_TYPE = boolType("Bool")
	NULL_TYPE = nullType("Null")
)

//
// value types
//

type ValueType interface {
	String() string
}

type stringType string

func (typ stringType) String() string {
	return string(typ)
}

type integerType string

func (typ integerType) String() string {
	return string(typ)
}

type floatType string

func (typ floatType) String() string {
	return string(typ)
}

type boolType string

func (typ boolType) String() string {
	return string(typ)
}

type nullType string

func (typ nullType) String() string {
	return string(typ)
}

//
// value
//

type Value interface {
	GetType() ValueType
	GetValue() interface{}

	SetValue(value interface{})
}

func NewValue(typ ValueType, value interface{}) Value {
	base := baseValue{
		typ: typ,
		value: value,
	}

	if typ == STRING_TYPE {
		return &stringValue{base}
	}
	if typ == INTEGER_TYPE {
		return &integerValue{base}
	}
	if typ == FLOAT_TYPE {
		return &floatValue{base}
	}
	if typ == BOOL_TYPE {
		return &boolValue{base}
	}
	if typ == NULL_TYPE {
		return &nullValue{
			baseValue: baseValue{
				typ: typ,
				value: nil,
			},
		}
	}
	panic("Invalid value type: " + typ.String())
}

type baseValue struct {
	typ ValueType
	value interface{}
}

func (value *baseValue) GetType() ValueType {
	return value.typ
}

func (value *baseValue) GetValue() interface{} {
	return value.value
}

func (value *baseValue) SetValue(v interface{}) {
	value.value = v
}

//
// default implements
//

func defaultOperation(op string, left, right Value) (Value, gerror.Error) {
	return nil, gerror.NewInvalidOperationError(nil, op,
		left.GetType().String(), right.GetType().String())
}

func (value *baseValue) AddString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(ADD, value, other)
}

func (value *baseValue) AddInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(ADD, value, other)
}

func (value *baseValue) AddFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(ADD, value, other)
}

func (value *baseValue) AddBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(ADD, value, other)
}

func (value *baseValue) AddNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(ADD, value, other)
}

func (value *baseValue) SubtractString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(SUBTRACT, value, other)
}

func (value *baseValue) SubtractInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(SUBTRACT, value, other)
}

func (value *baseValue) SubtractFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(SUBTRACT, value, other)
}

func (value *baseValue) SubtractBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(SUBTRACT, value, other)
}

func (value *baseValue) SubtractNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(SUBTRACT, value, other)
}

func (value *baseValue) MultiplyString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(MULTIPLY, value, other)
}

func (value *baseValue) MultiplyInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(MULTIPLY, value, other)
}

func (value *baseValue) MultiplyFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(MULTIPLY, value, other)
}

func (value *baseValue) MultiplyBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(MULTIPLY, value, other)
}

func (value *baseValue) MultiplyNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(MULTIPLY, value, other)
}

func (value *baseValue) DivideString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(DIVIDE, value, other)
}

func (value *baseValue) DivideInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(DIVIDE, value, other)
}

func (value *baseValue) DivideFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(DIVIDE, value, other)
}

func (value *baseValue) DivideBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(DIVIDE, value, other)
}

func (value *baseValue) DivideNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(DIVIDE, value, other)
}

func (value *baseValue) ModString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(MOD, value, other)
}

func (value *baseValue) ModInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(MOD, value, other)
}

func (value *baseValue) ModFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(MOD, value, other)
}

func (value *baseValue) ModBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(MOD, value, other)
}

func (value *baseValue) ModNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(MOD, value, other)
}

func (value *baseValue) EqualString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) EqualInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) EqualFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) EqualBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) EqualNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) NotEqualString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) NotEqualInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) NotEqualFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) NotEqualBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) NotEqualNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanOrEqualString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanOrEqualBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) GreaterThanOrEqualNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanOrEqualString(other *stringValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanOrEqualBool(other *boolValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

func (value *baseValue) LessThanOrEqualNull(other *nullValue) (Value, gerror.Error) {
	return defaultOperation(GT, value, other)
}

type stringValue struct {
	baseValue
}

//
// Add operation for stringValue
//

func (value *stringValue) AddString(other *stringValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: value.value.(string) + other.value.(string),
		},
	}, nil
}

func (value *stringValue) AddInteger(other *integerValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: value.value.(string) + fmt.Sprintf("%d", other.value.(int64)),
		},
	}, nil
}

func (value *stringValue) AddFloat(other *floatValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: value.value.(string) + fmt.Sprintf("%f", other.value.(float64)),
		},
	}, nil
}

func (value *stringValue) AddBool(other *boolValue) (Value, gerror.Error) {
	var str string
	if other.value.(bool) {
		str = "true"
	} else {
		str = "false"
	}

	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: value.value.(string) + str,
		},
	}, nil
}

func (value *stringValue) AddNull(other *nullValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: value.value.(string) + "null",
		},
	}, nil
}

//
// Multiply operation for stringValue
//

func (value *stringValue) MultiplyInteger(other *integerValue) (Value, gerror.Error) {
	cnt := int(other.value.(int64))
	if cnt < 0 {
		return nil, gerror.NewInvalidOperationError(nil, MULTIPLY,
			value.GetType().String(), "negtive integer")
	}

	buffer := []byte{}
	unit := []byte(value.value.(string))

	for i := 0; i < cnt; i++ {
		buffer = append(buffer, unit...)
	}

	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: (string)(buffer),
		},
	}, nil
}

type integerValue struct {
	baseValue
}

//
// Equal operation for stringValue
//

func (value *stringValue) EqualString(other *stringValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(string) == other.value.(string),
		},
	}, nil
}

//
// Not Equal operation for stringValue
//

func (value *stringValue) NotEqualString(other *stringValue) (Value, gerror.Error) {
	result, err := value.EqualString(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Greater than operation for stringValue
//

func (value *stringValue) GreaterThanString(other *stringValue) (Value, gerror.Error) {
	left := value.value.(string)
	right := value.value.(string)

	for i := 0; i < math.MinInt(len(left), len(right)); i++ {
		if left[i] < right[i] {
			return &boolValue{
				baseValue: baseValue{
					typ: BOOL_TYPE,
					value: false,
				},
			}, nil
		} else if left[i] > right[i] {
			return &boolValue{
				baseValue: baseValue{
					typ: BOOL_TYPE,
					value: true,
				},
			}, nil
		}
	}

	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: len(left) > len(right),
		},
	}, nil
}

//
// Greater than or equal operation for stringValue
//

func (value *stringValue) GreaterThanOrEqualString(other *stringValue) (Value, gerror.Error) {
	left := value.value.(string)
	right := value.value.(string)

	for i := 0; i < math.MinInt(len(left), len(right)); i++ {
		if left[i] < right[i] {
			return &boolValue{
				baseValue: baseValue{
					typ: BOOL_TYPE,
					value: false,
				},
			}, nil
		} else if left[i] > right[i] {
			return &boolValue{
				baseValue: baseValue{
					typ: BOOL_TYPE,
					value: true,
				},
			}, nil
		}
	}

	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: len(left) >= len(right),
		},
	}, nil
}

//
// Less operation for stringValue
//

func (value *stringValue) LessThanString(other *stringValue) (Value, gerror.Error) {
	result, err := value.GreaterThanOrEqualString(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Less than or equal operation for stringValue
//

func (value *stringValue) LessThanOrEqualString(other *stringValue) (Value, gerror.Error) {
	result, err := value.GreaterThanString(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Add operation for integerValue
//

func (value *integerValue) AddString(other *stringValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: fmt.Sprintf("%d", value.value.(int64)) + other.value.(string),
		},
	}, nil
}

func (value *integerValue) AddInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: INTEGER_TYPE,
			value: value.value.(int64) + other.value.(int64),
		},
	}, nil
}

func (value *integerValue) AddFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: float64(value.value.(int64)) + other.value.(float64),
		},
	}, nil
}

//
// Subtract operation for integerValue
//

func (value *integerValue) SubtractInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: INTEGER_TYPE,
			value: value.value.(int64) - other.value.(int64),
		},
	}, nil
}

func (value *integerValue) SubtractFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: float64(value.value.(int64)) - other.value.(float64),
		},
	}, nil
}

//
// Multiply operation for integerValue
//

func (value *integerValue) MultiplyString(other *stringValue) (Value, gerror.Error) {
	cnt := int(value.value.(int64))
	if cnt < 0 {
		return nil, gerror.NewInvalidOperationError(nil, MULTIPLY,
			"negtive integer", value.GetType().String())
	}

	buffer := []byte{}
	unit := ([]byte)(other.value.(string))

	for i := 0; i < cnt; i++ {
		buffer = append(buffer, unit...)
	}

	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: string(buffer),
		},
	}, nil
}

func (value *integerValue) MultiplyInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: INTEGER_TYPE,
			value: value.value.(int64) * other.value.(int64),
		},
	}, nil
}

func (value *integerValue) MultiplyFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: float64(value.value.(int64)) * other.value.(float64),
		},
	}, nil
}

//
// Divide operation for integerValue
//

func (value *integerValue) DivideInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: INTEGER_TYPE,
			value: value.value.(int64) / other.value.(int64),
		},
	}, nil
}

func (value *integerValue) DivideFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: float64(value.value.(int64)) / other.value.(float64),
		},
	}, nil
}

//
// Mod operation for integerValue
//

func (value *integerValue) ModInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: INTEGER_TYPE,
			value: value.value.(int64) % other.value.(int64),
		},
	}, nil
}

//
// Equal operation for integerValue
//

func (value *integerValue) EqualInteger(other *integerValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(int64) == other.value.(int64),
		},
	}, nil
}

//
// Not Equal operation for integerValue
//

func (value *integerValue) NotEqualInteger(other *integerValue) (Value, gerror.Error) {
	result, err := value.EqualInteger(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Greater than operation for integerValue
//

func (value *integerValue) GreaterThanInteger(other *integerValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(int64) > other.value.(int64),
		},
	}, nil
}

func (value *integerValue) GreaterThanFloat(other *floatValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: float64(value.value.(int64)) > other.value.(float64),
		},
	}, nil
}

//
// Greater than or equal operation for integerValue
//

func (value *integerValue) GreaterThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(int64) >= other.value.(int64),
		},
	}, nil
}

func (value *integerValue) GreaterThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: float64(value.value.(int64)) >= other.value.(float64),
		},
	}, nil
}

//
// Less than operation for integerValue
//

func (value *integerValue) LessThanInteger(other *integerValue) (Value, gerror.Error) {
	result, err := value.GreaterThanOrEqualInteger(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

func (value *integerValue) LessThanFloat(other *floatValue) (Value, gerror.Error) {
	result, err := value.GreaterThanOrEqualFloat(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Less than or equal operation for integerValue
//

func (value *integerValue) LessThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	result, err := value.GreaterThanInteger(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

func (value *integerValue) LessThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	result, err := value.GreaterThanFloat(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

type floatValue struct {
	baseValue
}

//
// Add operation for floatValue
//

func (value *floatValue) AddString(other *stringValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: fmt.Sprintf("%f", value.value.(float64)) + other.value.(string),
		},
	}, nil
}

func (value *floatValue) AddInteger(other *integerValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) + float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) AddFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) + other.value.(float64),
		},
	}, nil
}

//
// Subtract operation for floatValue
//

func (value *floatValue) SubtractInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) - float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) SubtractFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) - other.value.(float64),
		},
	}, nil
}

//
// Multiply operation for integerValue
//

func (value *floatValue) MultiplyInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) * float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) MultiplyFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) * other.value.(float64),
		},
	}, nil
}

//
// Divide operation for integerValue
//

func (value *floatValue) DivideInteger(other *integerValue) (Value, gerror.Error) {
	return &integerValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) / float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) DivideFloat(other *floatValue) (Value, gerror.Error) {
	return &floatValue{
		baseValue: baseValue{
			typ: FLOAT_TYPE,
			value: value.value.(float64) / other.value.(float64),
		},
	}, nil
}

//
// Equal operation for floatValue
//

func (value *floatValue) EqualFloat(other *floatValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(float64) == other.value.(float64),
		},
	}, nil
}

//
// Not Equal operation for integerValue
//

func (value *floatValue) NotEqualFloat(other *floatValue) (Value, gerror.Error) {
	result, err := value.EqualFloat(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Greater than operation for integerValue
//

func (value *floatValue) GreaterThanInteger(other *integerValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(float64) > float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) GreaterThanFloat(other *floatValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(float64) > other.value.(float64),
		},
	}, nil
}

//
// Greater than or equal operation for integerValue
//

func (value *floatValue) GreaterThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(float64) >= float64(other.value.(int64)),
		},
	}, nil
}

func (value *floatValue) GreaterThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	return &boolValue{
		baseValue: baseValue{
			typ: BOOL_TYPE,
			value: value.value.(float64) >= other.value.(float64),
		},
	}, nil
}

//
// Less than operation for integerValue
//

func (value *floatValue) LessThanInteger(other *integerValue) (Value, gerror.Error) {
	result, err := value.GreaterThanOrEqualInteger(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

func (value *floatValue) LessThanFloat(other *floatValue) (Value, gerror.Error) {
	result, err := value.GreaterThanOrEqualFloat(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

//
// Less than or equal operation for integerValue
//

func (value *floatValue) LessThanOrEqualInteger(other *integerValue) (Value, gerror.Error) {
	result, err := value.GreaterThanInteger(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

func (value *floatValue) LessThanOrEqualFloat(other *floatValue) (Value, gerror.Error) {
	result, err := value.GreaterThanFloat(other)
	if err != nil {
		return nil, err
	}
	result.SetValue(!result.GetValue().(bool))
	return result, nil
}

type boolValue struct {
	baseValue
}

//
// Add operation for boolValue
//

func (value *boolValue) AddString(other *stringValue) (Value, gerror.Error) {
	var str string
	if value.value.(bool) {
		str = "true"
	} else {
		str = "false"
	}

	return &floatValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: str + other.value.(string),
		},
	}, nil
}

type nullValue struct {
	baseValue
}

//
// Add operation for null value
//

func (value *nullValue) AddString(other *stringValue) (Value, gerror.Error) {
	return &stringValue{
		baseValue: baseValue{
			typ: STRING_TYPE,
			value: "null" + other.value.(string),
		},
	}, nil
}