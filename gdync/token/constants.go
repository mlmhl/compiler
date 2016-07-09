package token

const (
	NUMBER = "0|1|2|3|4|5|6|7|8|9"
	ALPHABET = "A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z" +
	"|" + "a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z"

	STRING = "(\".*\")"
	INTEGER = "(0|((1|2|3|4|5|6|7|8|9)(0|1|2|3|4|5|6|7|8|9)*))"
	FLOAT = "((" + NUMBER + ")+\\.(" + NUMBER + ")+)"

	TRUE = "true"
	FALSE = "false"

	LSP = "(\\()"
	RSP = "(\\))"
	LMP = "(\\[)"
	RMP = "(\\])"
	LLP = "(\\{)"
	RLP = "(\\})"

	COMMA = "(,)"
	SEMICOLON = "(;)"

	ADD = "(\\+)"
	SUBTRACT = "(\\-)"
	MULTIPLY = "(\\*)"
	DIVIDE = "(/)"
	MOD = "(%)"

	NOT = "(!)"
	OR = "(\\|\\|)"
	AND = "(&&)"

	EQUAL = "(==)"
	UNEQUAL = "(!=)"
	GT = "(>)"
	LT = "(<)"
	GTE = "(>=)"
	LTE = "(<=)"

	ASSIGN = "(=)"

	FOR = "(for)"
	WHILE = "(while)"
	BREAK = "(break)"
	CONTINUE = "(continue)"

	IF = "(if)"
	ELSE = "(else)"
	ELIF = "(elif)"

	FUNCTION_DEFINITION = "(def)"
	RETURN = "(return)"

	NULL = "(null)"
	GLOBAL = "(global)"

	WHITESPACE = "(( |\t|\n)+)"

	COMMENT = "//"

	IDENTIFIER = "((" + ALPHABET + ")" +
	"(" + ALPHABET + "|" + NUMBER + "|_" + ")*)"
)

const (
	FINISHED_ID = iota + 1
	UNKNOWN

	STRING_ID
	INTEGER_ID
	FLOAT_ID

	TRUE_ID
	FALSE_ID

	LSP_ID
	RSP_ID
	LMP_ID
	RMP_ID
	LLP_ID
	RLP_ID

	COMMA_ID
	SEMICOLON_ID

	ADD_ID
	SUBTRACT_ID
	MULTIPLY_ID
	DIVIDE_ID
	MOD_ID

	NOT_ID
	OR_ID
	AND_ID

	EQUAL_ID
	UNEQUAL_ID
	GT_ID
	LT_ID
	GTE_ID
	LTE_ID

	ASSIGN_ID

	FOR_ID
	WHILE_ID
	BREAK_ID
	CONTINUE_ID

	IF_ID
	ELSE_ID
	ELIF_ID

	FUNCTION_DEFINITION_ID
	RETURN_ID

	NULL_ID
	GLOBAL_ID

	WHITESPACE_ID

	IDENTIFIER_ID
)

var descriptions map[int]string = map[int]string {
	FINISHED_ID: "finished",

	STRING_ID: "string",
	INTEGER_ID: "integer",
	FLOAT_ID: "float",

	TRUE_ID: "true",
	FALSE_ID: "false",

	LSP_ID: "left small parentheses",
	RSP_ID: "right small parentheses",
	LMP_ID: "left middle parentheses",
	RMP_ID: "right middle parentheses",
	LLP_ID: "left large parentheses",
	RLP_ID: "right large parentheses",

	COMMA_ID: "comma",
	SEMICOLON_ID: "semicolon",

	ADD_ID: "add",
	SUBTRACT_ID: "subtract",
	MULTIPLY_ID: "multiply",
	DIVIDE_ID: "divide",
	MOD_ID: "mod",

	NOT_ID: "not",
	OR_ID: "or",
	AND_ID: "and",

	EQUAL_ID: "equal",
	UNEQUAL_ID: "unequal",
	GT_ID: "greate than",
	LT_ID: "less than",
	GTE_ID: "greate than and equal",
	LTE_ID: "less than and euqal",

	ASSIGN_ID: "assign",

	FOR_ID: "for",
	WHILE_ID: "while",
	BREAK_ID: "break",
	CONTINUE_ID: "continue",

	IF_ID: "if",
	ELSE_ID: "else",
	ELIF_ID: "elif",

	IDENTIFIER_ID: "identifier",

	FUNCTION_DEFINITION_ID: "function",
	RETURN_ID: "return",

	NULL_ID: "null",
	GLOBAL_ID: "global",

	WHITESPACE_ID: "white space",
}