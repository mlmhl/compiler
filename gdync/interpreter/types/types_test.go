package types

import (
	"testing"

	gerror "github.com/mlmhl/compiler/gdync/errors"
	"fmt"
)

func testArithmeticOperation(op string, left, right Value,
	target interface{}, t *testing.T) {
	t.Logf("Test: %s on %s, %s ...", op, left.GetType().String(), right.GetType().String())

	res, err := ArithmeticOperation(op, left, right)
	if value, ok := target.(Value); ok {
		if err != nil {
			t.Fatal("Unexpected error: " + err.(gerror.Error).GetMessage())
		}
		if value.GetValue() != res.GetValue() {
			t.Fatalf("Wrong result: Wanted %v, got %v", value, res)
		}
	} else {
		if err == nil {
			t.Fatal("There should be an error, but found nil")
		}
		if err.GetMessage() != target.(gerror.Error).GetMessage() {
			t.Fatalf("Wrong error message: Wanted (%s), got (%s)",
				err.GetMessage(), target.(gerror.Error))
		}
	}

	t.Log("Passed ...")
}

func TestAddOperation(t *testing.T) {
	i1 := NewValue(INTEGER_TYPE, int64(3))
	i2 := NewValue(INTEGER_TYPE, int64(10))

	f1 := NewValue(FLOAT_TYPE, float64(2.5))
	f2 := NewValue(FLOAT_TYPE, float64(7.5))

	s1 := NewValue(STRING_TYPE, "Hello")
	s2 := NewValue(STRING_TYPE, "World")

	b1 := NewValue(BOOL_TYPE, true)
	b2 := NewValue(BOOL_TYPE, false)

	n := NewValue(NULL_TYPE, nil)

	// test add for integer
	testArithmeticOperation(ADD, i1, s1, NewValue(STRING_TYPE, "3Hello"), t)
	testArithmeticOperation(ADD, i1, i2, NewValue(INTEGER_TYPE, int64(13)), t)
	testArithmeticOperation(ADD, i1, f1, NewValue(FLOAT_TYPE, float64(5.5)), t)

	// test subtract for integer
	testArithmeticOperation(SUBTRACT, i1, i2, NewValue(INTEGER_TYPE, int64(-7)), t)
	testArithmeticOperation(SUBTRACT, i1, f1, NewValue(FLOAT_TYPE, float64(0.5)), t)

	// test multiply for integer
	testArithmeticOperation(MULTIPLY, i1, s1, NewValue(STRING_TYPE, "HelloHelloHello"), t)
	testArithmeticOperation(MULTIPLY, i1, i2, NewValue(INTEGER_TYPE, int64(30)), t)
	testArithmeticOperation(MULTIPLY, i1, f1, NewValue(FLOAT_TYPE, float64(7.5)), t)

	// test divide for integer
	testArithmeticOperation(DIVIDE, i2, i1, NewValue(INTEGER_TYPE, int64(3)), t)
	testArithmeticOperation(DIVIDE, i2, f1, NewValue(FLOAT_TYPE, float64(4.0)), t)

	// test mod for integer
	testArithmeticOperation(MOD, i2, i1, NewValue(INTEGER_TYPE, int64(1)), t)

	// test add for float
	testArithmeticOperation(ADD, f1, s1, NewValue(STRING_TYPE, fmt.Sprintf("%f",2.5) + "Hello"), t)
	testArithmeticOperation(ADD, f1, i1, NewValue(FLOAT_TYPE, float64(5.5)), t)
	testArithmeticOperation(ADD, f1, f2, NewValue(FLOAT_TYPE, float64(10.0)), t)

	// test subtract for float
	testArithmeticOperation(SUBTRACT, f1, i1, NewValue(FLOAT_TYPE, float64(-0.5)), t)
	testArithmeticOperation(SUBTRACT, f2, f1, NewValue(FLOAT_TYPE, float64(5.0)), t)

	// test multiply for float
	testArithmeticOperation(MULTIPLY, f1, i1, NewValue(FLOAT_TYPE, float64(7.5)), t)
	testArithmeticOperation(MULTIPLY, f1, f2, NewValue(FLOAT_TYPE, float64(2.5 * 7.5)), t)

	// test divide for float
	testArithmeticOperation(DIVIDE, f2, i1, NewValue(FLOAT_TYPE, float64(2.5)), t)
	testArithmeticOperation(DIVIDE, f2, f1, NewValue(FLOAT_TYPE, float64(3.0)), t)

	// test add for string
	testArithmeticOperation(ADD, s1, s2, NewValue(STRING_TYPE, "HelloWorld"), t)
	testArithmeticOperation(ADD, s1, i1, NewValue(STRING_TYPE, "Hello3"), t)
	testArithmeticOperation(ADD, s1, f1, NewValue(STRING_TYPE, "Hello" + fmt.Sprintf("%f", 2.5)), t)
	testArithmeticOperation(ADD, s1, b1, NewValue(STRING_TYPE, "Hellotrue"), t)
	testArithmeticOperation(ADD, s1, n, NewValue(STRING_TYPE,"Hellonull"), t)

	// test multiply for string
	testArithmeticOperation(MULTIPLY, s1, i1, NewValue(STRING_TYPE, "HelloHelloHello"), t)

	// test add for bool
	testArithmeticOperation(ADD, b2, s2, NewValue(STRING_TYPE, "falseWorld"), t)

	// test for unsupported operation
	testArithmeticOperation(MULTIPLY, s1, s2, gerror.NewInvalidOperationError(nil,
		MULTIPLY, s1.GetType().String(), s2.GetType().String()), t)
}
