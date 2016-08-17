package codebyte

import (
	"io"
	"os"

	files "github.com/mlmhl/goutil/io/files"
)

const (
	THRESHOLD = 1024 * 1024 // 1M
)

type CodeByteBuffer struct {
	buffer []byte
	file   *os.File
	empty  bool
}

func NewCodeByteBuffer(fileName string) (*CodeByteBuffer, error) {
	file, err := files.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &CodeByteBuffer{
		buffer: []byte{},
		file:   file,
		empty:  false,
	}, nil
}

func (buffer *CodeByteBuffer) IsEmpty() bool {
	return buffer.empty
}

func (buffer *CodeByteBuffer) Write(b byte) {
	buffer.WriteSlice([]byte{b})

}

func (buffer *CodeByteBuffer) WriteSlice(slice []byte) {
	buffer.buffer = append(buffer.buffer, slice...)
	if len(buffer.buffer) >= THRESHOLD {
		buffer.Sync()
	}
}

func (buffer *CodeByteBuffer) Read(b []byte) (int, error) {
	if buffer.IsEmpty() {
		return 0, nil
	}
	cnt, err := buffer.file.Read(b)
	if err != nil {
		if err == io.EOF {
			buffer.empty = true
		} else {
			return cnt, err
		}
	}
	return cnt, nil
}

func (buffer *CodeByteBuffer) Sync() error {
	buffer.flush()
	return buffer.file.Sync()
}

func (buffer *CodeByteBuffer) Close() error {
	return buffer.file.Close()
}

func (buffer *CodeByteBuffer) flush() {
	for {
		cnt, _ := buffer.file.Write(buffer.buffer)
		if cnt == len(buffer.buffer) {
			break
		}
		buffer.buffer = buffer.buffer[cnt:]
	}
	buffer.buffer = []byte{}
}
