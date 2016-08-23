package executable

type OperatorCode byte

var (
	// Put a constant value to VM's stack
	PUSH_NULL       byte = 0
	PUSH_BOOL_TRUE  byte = 1
	PUSH_BOOL_FALSE byte = 2
	PUSH_INT        byte = 3
	PUSH_FLOAT      byte = 4
	PUSH_STRING     byte = 5

	// Create a array with a literal value(may be bot specified)
	NEW_ARRAY                byte
	NEW_ARRAY_LITERAL_BOOL   byte
	NEW_ARRAY_LITERAL_INT    byte
	NEW_ARRAY_LITERAL_DOUBLE byte
	NEW_ARRAY_LITERAL_OBJECT byte

	// Reference a variable through IdentiferExpression
	VARIABLE_REFERENCE byte

	// Add operator
	ADD_INT byte
	ADD_FLOAT byte
	ADD_STRING byte

	// Sub operator
	SUBTRACT_INT byte
	SUBTRACT_FLOAT byte
)

func GetOperatorCode(start OperatorCode, offset int) {
	return start + byte(offset)
}