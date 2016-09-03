package executable

import (
	"github.com/mlmhl/compiler/gstac/errors"
)

// Executable is a collection of a source file's code byte.
type Executable struct {
	labelPool      *LabelPool
	constantPool   *ConstantPool
	codeByteBuffer *CodeByteBuffer

	// size for current object
	size int
	// current object is the first one or not
	isFirst bool

	// Support for continue and break statement
	currentBreakLabel int
	currentContinueLabel int
}

func NewExecutable() *Executable {
	return &Executable{
		labelPool:    NewLabelPool(),
		constantPool: NewConstantPool(),

		size:    0,
		isFirst: true,

		currentBreakLabel:    -1,
		currentContinueLabel: -1,
	}
}

func (executable *Executable) Open(fileName string) errors.Error {
	codeByteBuffer, err := NewCodeByteBuffer(fileName)
	if err != nil {
		return errors.NewExecutableFileCreationError(err)
	}
	executable.codeByteBuffer = codeByteBuffer
	return nil
}

func (executable *Executable) AddSymbolList(code []byte) {
	executable.codeByteBuffer.WriteSlice(code)
}

// Start write a list of some object to file
func (executable *Executable) StartList() {
	executable.codeByteBuffer.Write('{')
}

// StartObject should be invoked before process a new object(Function variable etc.)
func (executable *Executable) StartObject() {
	if executable.isFirst {
		executable.isFirst = false
	} else {
		executable.codeByteBuffer.Write(',')
	}
}

func (executable *Executable) Append(b byte) {
	executable.codeByteBuffer.Write(b)
	executable.size += 1
}

func (executable *Executable) AppendSlice(segment []byte) {
	executable.codeByteBuffer.WriteSlice(segment)
	executable.size += len(segment)
}

func (executable *Executable) GetSize() int {
	return executable.size
}

func (executable *Executable) SetContinueLabel(label int) {
	executable.currentContinueLabel = label
}

func (executable *Executable) GetContinueLabel() int {
	return executable.currentContinueLabel
}

func (executable *Executable) ResetContinueLabel() {
	executable.currentContinueLabel = -1
}

func (executable *Executable) SetBreakLabel(label int) {
	executable.currentBreakLabel = label
}

func (executable *Executable) GetBreakLabel() int {
	return executable.currentBreakLabel
}

func (executable *Executable) ResetBreakLabel() {
	executable.currentBreakLabel = -1
}

// Can't add new code segment of the same object after EndObject invoked
func (executable *Executable) EndObject() {
	// reset current object's size back to zero
	executable.size = 0
}

// End write a list of some object to file
func (executable *Executable) EndList() {
	executable.codeByteBuffer.Write('}')
	executable.isFirst = true
}

func (executable *Executable) NewLabel() int {
	return executable.labelPool.NewLabel()
}

func (executable *Executable) SetLabel(label, address int) {
	executable.labelPool.SetLabel(label, address)
}
func (executable *Executable) AddConstantValue(value interface{}) int {
	return executable.constantPool.AddIfAbsent(value)
}

// Close write constant pool to executable file and close the executable
func (executable *Executable) Close() {
	executable.writeConstantPool()
	executable.codeByteBuffer.Sync()
	executable.codeByteBuffer.Close()
}

func (executable *Executable) writeConstantPool() {
	executable.codeByteBuffer.WriteSlice(executable.constantPool.Encode())
}
