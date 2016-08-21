package executable

import (
	"testing"

	files "github.com/mlmhl/goutil/io/files"
)

const (
	fileName = "test"
)

var (
	buffer *CodeByteBuffer = nil
)

func TestCodeByteBuffer_Write(t *testing.T) {
	setup()

	t.Log("Test: CodeByteBuffer Write...")

	slice1 := []byte("hello world")
	for _, b := range slice1 {
		buffer.Write(b)
	}
	readFromCodeByteBuffer(buffer, slice1, t)

	t.Log("Passed...")

	cleanup()
}

func TestCodeByteBuffer_WriteSlice(t *testing.T) {
	setup()

	t.Log("Test: codeByteBuffer WriteSlice...")

	slice2 := []byte("Hello World")
	buffer.WriteSlice(slice2)
	readFromCodeByteBuffer(buffer, slice2, t)

	t.Log("Passed...")

	cleanup()
}

func setup() {
	if buffer == nil {
		buffer, _ = NewCodeByteBuffer(fileName)
	}
}

func readFromCodeByteBuffer(buffer *CodeByteBuffer, target []byte, t *testing.T) {
	content := []byte{}
	buf := make([]byte, 1024)

	buffer.Sync()
	for !buffer.IsEmpty() {
		cnt, err := buffer.Read(buf)
		if err != nil {
			t.Fatalf("Read from CodeByteBuffer error: %v", err)
		}
		content = append(content, buf[:cnt]...)
	}

	if string(content) != string(target) {
		t.Fatalf("Wrong content: Wanted `%s`, got `%s`", string(target), string(content))
	}

	cleanup()
}

func cleanup() {
	if buffer != nil {
		buffer.Close()
		files.Remove(fileName)
		buffer = nil
	}
}
