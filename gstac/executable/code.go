package executable

const (
	// Put a constant value to VM's stack
	PUSH_NULL = byte(iota)
	PUSH_BOOL_TRUE
	PUSH_BOOL_FALSE
	PUSH_INT
	PUSH_FLOAT
	PUSH_STRING

	// Create a array with a literal value(may be bot specified)
	NEW_ARRAY
	NEW_ARRAY_LITERAL_BOOL
	NEW_ARRAY_LITERAL_INT
	NEW_ARRAY_LITERAL_DOUBLE
	NEW_ARRAY_LITERAL_OBJECT

	// Reference a variable through IdentifierExpression
	VARIABLE_REFERENCE

	// normal assign expression
	NORMAL_ASSIGN

	// Add operator
	ADD_BOOL
	ADD_INT
	ADD_FLOAT
	ADD_STRING

	// Sub operator
	SUBTRACT_INT
	SUBTRACT_FLOAT

	// Multiply operator
	MULTIPLY_INT
	MULTIPLY_FLOAT

	// Divide operator
	DIVIDE_INT
	DIVIDE_FLOAT

	// Mod operator
	MOD_INT

	// Support for continuous assignment like 'a=b=c'
	STACK_TOP_DUPLICATE

	// Pop the stack top to a variable
	POP_STACK_BOOL
	POP_STACK_INT
	POP_STACK_FLOAT
	POP_STACK_OBJECT

	// Pop the latest static value to a variable
	POP_STATIC_BOOL
	POP_STATIC_INT
	POP_STATIC_FLOAT
	POP_STATIC_OBJECT

	// The stack top is array, then index, then value, put value into array[index]
	POP_ARRAY_BOOL
	POP_ARRAY_INT
	POP_ARRAY_FLOAT
	POP_ARRAY_OBJECT
)

func GetOperatorCode(start byte, offset int) byte {
	return start + byte(offset)
}
