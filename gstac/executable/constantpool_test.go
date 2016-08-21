package executable

import (
	"testing"
)

func TestConstantPool(t *testing.T) {
	pool := NewConstantPool()
	values := []interface{}{
		1,
		3.5,
		"hello",
		108.801,
		"world",
		37,
	}

	t.Log("Test: ConstantPool first append...")

	for target, value := range(values) {
		index := pool.AddIfAbsent(value)
		if index != target {
			t.Fatalf("Wrong index(%v): Wanted %d, got %d", value, target, index)
		}
	}

	t.Log("Passed...")

	t.Log("Test: ConstantPool exist append...")

	for target, value := range(values) {
		index := pool.AddIfAbsent(value)
		if index != target {
			t.Fatalf("Wrong index(%v): Wanted %d, got %d", value, target, index)
		}
	}


	t.Log("Passed...")
}