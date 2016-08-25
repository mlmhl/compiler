package common

import "github.com/mlmhl/goutil/encoding"

type Location struct {
	line     int
	position int
	fileName string
}

func NewLocation(line, position int, fileName string) *Location {
	return &Location{
		line: line,
		position: position,
		fileName: fileName,
	}
}

func (location *Location) GetLine() int {
	return location.line
}

func (location *Location) GetPosition() int {
	return location.position
}

func (location *Location) GetFileName() string {
	return location.fileName
}

// for test
func (location *Location) Equal(other *Location) bool {
	return location.line == other.line && location.position == other.position &&
		location.fileName == other.fileName
}

// Encode encode the line and position to code byte
func (location *Location) Encode() []byte {
	buffer := encoding.DefaultEncoder.Int(location.line)
	buffer = append(buffer, encoding.DefaultEncoder.Int(location.position)...)
	return buffer
}

// Decode decode a Location from code byte, invoker should provide the file name
func (location *Location) Decode(buffer []byte, fileName string) (*Location, []byte) {
	line := encoding.DefaultDecoder.Int(buffer)
	position := encoding.DefaultDecoder.Int(buffer[8:])
	return NewLocation(line, position, fileName), buffer[16:]
}