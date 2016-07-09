package ast

import (
	"reflect"
	"testing"

	"github.com/mlmhl/compiler/gdync/interpreter/types"
)

func functionTest(function Function, arguments []types.Value,
	env *Environment, targetRes types.Value, t *testing.T) {
	t.Logf("Test: %s ...", reflect.TypeOf(function).Elem().Name())

	if res, err := function.Evaluate(arguments, env); err != nil {
		if err.GetLocation() == nil {
			t.Fatalf("Execute error: " + err.GetMessage())
		} else {
			t.Fatalf("Execute error: %s, at %s, %s, %d", err.GetMessage(),
				err.GetLocation().GetFileName(), err.GetLocation().GetLine(),
				err.GetLocation().GetPosition())
		}
	} else {
		if targetRes == nil {
			if res != nil {
				t.Fatalf("Wrong result, Wanted nil, got %v", res)
			}
		} else {
			if res.GetType() != targetRes.GetType() {
				t.Fatalf("Wrong result type: Wanted %s, got %s",
					targetRes.GetType().String(), res.GetType().String())
			} else {
				if res.GetValue() != targetRes.GetValue() {
					t.Fatalf("Wrong result value, Wanted %v, got %v",
						targetRes.GetValue(), res.GetValue())
				}
			}
		}
	}

	t.Logf("Passed")
}

func TestNewPrintfFunction(t *testing.T) {
	function := NewPrintfFunction()
	arguments := []types.Value{}
	arguments = append(arguments, types.NewValue(types.STRING_TYPE, "My name is %s, I am %d years old"))
	arguments = append(arguments, types.NewValue(types.STRING_TYPE, "gdync"))
	arguments = append(arguments, types.NewValue(types.INTEGER_TYPE, 0))
	functionTest(function, arguments, nil, nil, t)
}
