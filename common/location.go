package common

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