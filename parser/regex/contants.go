package regex

// Mata symbol
const (
	EPSILON = 0

	CHOICE = '|'
	REPETITION = '*'
	ARBITRARY = '.'
	ZERO_OR_ONE = '?'
	ONE_OR_MORE = '+'

	// left small parentheses
	LSP = '('
	// right small parentheses
	RSP = ')'
	// left middle parentheses
	LMP = '['
	// right middle parentheses
	RMP = ']'
	// left larger parentheses
	LLP = '{'
	// right larger parentheses
	RLP = '}'
)

// Mata symbol id
const (
	EPSILON_ID = 0

	CHOICE_ID = 301
	REPETITION_ID = 302
	ARBITRARY_ID = 303
	ZERO_OR_ONE_ID = 304
	ONE_OR_MORE_ID = 305

	// left small parentheses
	LSP_ID = 401
	// right small parentheses
	RSP_ID = 402
	// left middle parentheses
	LMP_ID = 403
	// right middle parentheses
	RMP_ID = 404
	// left larger parentheses
	LLP_ID = 405
	// right larger parentheses
	RLP_ID = 406
)

var mataSymbolId map[int]int = map[int]int{
	EPSILON: EPSILON_ID,

	CHOICE: CHOICE_ID,
	REPETITION: REPETITION_ID,
	ARBITRARY: ARBITRARY_ID,
	ZERO_OR_ONE: ZERO_OR_ONE_ID,
	ONE_OR_MORE: ONE_OR_MORE_ID,

	LSP: LSP_ID,
	RSP: RSP_ID,
	LMP: LMP_ID,
	RMP: RMP_ID,
	LLP: LLP_ID,
	RLP: RLP_ID,
}

var mataSymbolSet map[int]bool = map[int]bool{
	EPSILON_ID: true,

	CHOICE_ID: true,
	REPETITION_ID: true,
	ARBITRARY_ID: true,
	ZERO_OR_ONE_ID: true,
	ONE_OR_MORE_ID: true,

	LSP_ID: true,
	RSP_ID: true,
	LMP_ID: true,
	RMP_ID: true,
	LLP_ID: true,
	RLP_ID: true,
}
