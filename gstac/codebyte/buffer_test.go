package codebyte

import (
	"testing"
)

func TestCodeByteBuffer(t *testing.T) {
	buffer := NewCodeByteBuffer()

	t.Log("Test: CodeByteBuffer Append...")

	slice1 := []byte("hello world")
	for _, b := range(slice1) {
		buffer.Append(b)
	}
	length := len(slice1)
	if buffer.Len() != length {
		t.Fatalf("Wrong size: Wanted %d, got %d", length, buffer.Len())
	}

	t.Log("Passed...")

	t.Log("Test: codeByteBuffer AppendSlice...")

	slice2 := []byte("Hello World")
	buffer.AppendSlice(slice2)
	length += len(slice1)
	if buffer.Len() != length {
		t.Fatalf("Wrong size: Wanted %d, got %d", length, buffer.Len())
	}

	t.Log("Passed...")

	t.Log("Test: CodeByteBuffer Pop...")

	bytes := []byte{}
	for !buffer.IsEmpty() {
		bytes = append(bytes, buffer.Pop())
	}
	target := string(slice1) + string(slice2)
	if string(bytes) != target {
		t.Fatalf("Wrong content: Wanted %s, got %s", target, string(bytes))
	}

	t.Log("Passed...")
}
