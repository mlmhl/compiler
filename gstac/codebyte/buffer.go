package codebyte

type CodeByteBuffer struct {
	buffer []byte
}

func NewCodeByteBuffer() *CodeByteBuffer {
	return &CodeByteBuffer{
		buffer: []byte{},
	}
}

func (buffer *CodeByteBuffer) IsEmpty() bool {
	return buffer.Len() == 0
}

func (buffer *CodeByteBuffer) Len() int {
	return len(buffer.buffer)
}

func (buffer *CodeByteBuffer) Append(b byte) {
	buffer.buffer = append(buffer.buffer, b)
}

func (buffer *CodeByteBuffer) AppendSlice(slice []byte) {
	buffer.buffer = append(buffer.buffer, slice...)
}

func (buffer *CodeByteBuffer) Pop() byte {
	b:= buffer.buffer[0]
	buffer.buffer = buffer.buffer[1:]
	return b
}