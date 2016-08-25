package executable

import (
	"github.com/mlmhl/compiler/gstac/errors"
)

type Executable struct {
	labelPool      *LabelPool
	constantPool   *ConstantPool
	codeByteBuffer *CodeByteBuffer

	isFirstFunc bool
	isFirstDeclaration bool
}

func NewExecutable() *Executable {
	return &Executable{
		labelPool:    NewLabelPool(),
		constantPool: NewConstantPool(),

		isFirstFunc: true,
		isFirstDeclaration: true,
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

// start write a list of some object to file
func (executable *Executable) BeginList() {
	executable.codeByteBuffer.Write('{')
}

// Add an element to the file, this method should be called after BeginList,
// and can't be called an more after EndList.
func (executable *Executable) AddElement(code []byte) {
	if executable.isFirstDeclaration {
		executable.isFirstDeclaration = false
	} else {
		executable.codeByteBuffer.Write(',')
	}
	executable.codeByteBuffer.WriteSlice(code)
}

// End write a list of some object to file
func (executable *Executable) EndList() {
	executable.codeByteBuffer.Write('}')
}

func (executable *Executable) AddGlobalCode(code []byte) {
	executable.codeByteBuffer.WriteSlice(code)
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