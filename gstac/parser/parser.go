package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/mlmhl/compiler/common"
	error "github.com/mlmhl/compiler/gstac/errors"
	"github.com/mlmhl/compiler/gstac/token"
	"github.com/mlmhl/compiler/regex"
)

type Parser struct {
	lineNumber int
	position   int
	line       string
	fileName   string
	regex      *regex.Regex
	scanner    *bufio.Scanner

	cursor int
	buffer []*token.Token
}

func NewParser() *Parser {
	regex := regex.NewRegex()

	regex.AddRegexExpression(token.STRING_TYPE, token.STRING_TYPE_ID)
	regex.AddRegexExpression(token.INTEGER_TYPE, token.INTEGER_TYPE_ID)
	regex.AddRegexExpression(token.FLOAT_TYPE, token.FLOAT_TYPE_ID)

	regex.AddRegexExpression(token.BOOL_TYPE, token.BOOL_TYPE_ID)
	regex.AddRegexExpression(token.STRING_VALUE, token.STRING_VALUE_ID)
	regex.AddRegexExpression(token.INTEGER_VALUE, token.INTEGER_VALUE_ID)
	regex.AddRegexExpression(token.FLOAT_VALUE, token.FLOAT_VALUE_ID)

	regex.AddRegexExpression(token.TRUE, token.TRUE_ID)
	regex.AddRegexExpression(token.FALSE, token.FALSE_ID)

	regex.AddRegexExpression(token.LSP, token.LSP_ID)
	regex.AddRegexExpression(token.RSP, token.RSP_ID)
	regex.AddRegexExpression(token.LMP, token.LMP_ID)
	regex.AddRegexExpression(token.RMP, token.RMP_ID)
	regex.AddRegexExpression(token.LLP, token.LLP_ID)
	regex.AddRegexExpression(token.RLP, token.RLP_ID)

	regex.AddRegexExpression(token.COMMA, token.COMMA_ID)
	regex.AddRegexExpression(token.SEMICOLON, token.SEMICOLON_ID)

	regex.AddRegexExpression(token.ADD, token.ADD_ID)
	regex.AddRegexExpression(token.SUBTRACT, token.SUBTRACT_ID)
	regex.AddRegexExpression(token.MULTIPLY, token.MULTIPLY_ID)
	regex.AddRegexExpression(token.DIVIDE, token.DIVIDE_ID)
	regex.AddRegexExpression(token.MOD, token.MOD_ID)

	regex.AddRegexExpression(token.NOT, token.NOT_ID)
	regex.AddRegexExpression(token.OR, token.OR_ID)
	regex.AddRegexExpression(token.AND, token.AND_ID)

	regex.AddRegexExpression(token.EQUAL, token.EQUAL_ID)
	regex.AddRegexExpression(token.UNEQUAL, token.UNEQUAL_ID)
	regex.AddRegexExpression(token.GT, token.GT_ID)
	regex.AddRegexExpression(token.LT, token.LT_ID)
	regex.AddRegexExpression(token.GTE, token.GTE_ID)
	regex.AddRegexExpression(token.GTE, token.GTE_ID)
	regex.AddRegexExpression(token.LTE, token.LTE_ID)

	regex.AddRegexExpression(token.ASSIGN, token.ASSIGN_ID)
	regex.AddRegexExpression(token.ADD_ASSIGN, token.ADD_ASSIGN_ID)
	regex.AddRegexExpression(token.SUB_ASSIGN, token.SUB_ASSIGN_ID)
	regex.AddRegexExpression(token.MUL_ASSIGN, token.MUL_ASSIGN_ID)
	regex.AddRegexExpression(token.DIV_ASSIGN, token.DIV_ASSIGN_ID)
	regex.AddRegexExpression(token.MOD_ASSIGN, token.MOD_ASSIGN_ID)

	regex.AddRegexExpression(token.INCREMENT, token.INCREMENT_ID)
	regex.AddRegexExpression(token.DECREMENT, token.DECREMENT_ID)

	regex.AddRegexExpression(token.FOR, token.FOR_ID)
	regex.AddRegexExpression(token.WHILE, token.WHILE_ID)
	regex.AddRegexExpression(token.BREAK, token.BREAK_ID)
	regex.AddRegexExpression(token.CONTINUE, token.CONTINUE_ID)

	regex.AddRegexExpression(token.IF, token.IF_ID)
	regex.AddRegexExpression(token.ELSE, token.ELSE_ID)
	regex.AddRegexExpression(token.ELIF, token.ELIF_ID)

	regex.AddRegexExpression(token.IDENTIFIER, token.IDENTIFIER_ID)

	regex.AddRegexExpression(token.FUNCTION_DEFINITION, token.FUNCTION_DEFINITION_ID)
	regex.AddRegexExpression(token.RETURN, token.RETURN_ID)

	regex.AddRegexExpression(token.NULL, token.NULL_ID)
	regex.AddRegexExpression(token.GLOBAL, token.GLOBAL_ID)

	regex.AddRegexExpression(token.NEW, token.NEW_ID)

	regex.AddRegexExpression(token.WHITESPACE, token.WHITESPACE_ID)

	regex.Compile()

	return &Parser{
		regex: regex,
	}
}

func (parser *Parser) Parse(fileName string) error.Error {
	if file, err := os.Open(fileName); err != nil {
		return error.NewInternalError(err.Error())
	} else {
		parser.scanner = bufio.NewScanner(file)
		parser.fileName = fileName
		parser.line = ""
		parser.lineNumber = 0
		parser.position = 0

		parser.cursor = -1
		parser.buffer = []*token.Token{}

		return nil
	}
}

// Get next token
func (parser *Parser) Next() (*token.Token, error.Error) {
	if !parser.HasNext() {
		if len(parser.buffer) > 0 &&
			parser.buffer[len(parser.buffer)-1].GetType() == token.FINISHED_ID {
			return parser.buffer[len(parser.buffer)-1], nil
		}

		tok := token.NewToken(common.NewLocation(
			-1, -1, parser.fileName)).SetType(token.FINISHED_ID)
		parser.appendToBuffer(tok)
		return tok, nil
	}

	if parser.cursor < len(parser.buffer)-1 {
		parser.cursor++
		tok := parser.buffer[parser.cursor]
		return tok, nil
	}

	for {
		length, types := parser.regex.Match(parser.line[parser.position:])
		pos := parser.position + length
		if len(types) == 0 {
			return nil, error.NewSyntaxError("Unsupported syntax",
				common.NewLocation(parser.lineNumber, parser.position, parser.fileName))
		} else {
			if types[0] == token.WHITESPACE_ID {
				// skip white space
				parser.position = pos
				continue
			}

			typ := types[0]
			value := parser.line[parser.position:pos]
			tok := token.NewToken(common.NewLocation(parser.lineNumber,
				parser.position, parser.fileName)).SetType(typ)

			switch typ {
			case token.STRING_VALUE_ID:
				tok.SetValue(value)

			case token.INTEGER_VALUE_ID:
				if v, err := strconv.Atoi(value); err != nil {
					return nil, error.NewSyntaxError("Unsupported integer syntax",
						common.NewLocation(parser.lineNumber, parser.position, parser.fileName))
				} else {
					tok.SetValue(int64(v))
				}
			case token.FLOAT_VALUE_ID:
				if v, err := strconv.ParseFloat(value, 64); err != nil {
					return nil, error.NewSyntaxError("Unsupported float synatx",
						common.NewLocation(parser.lineNumber, parser.position, parser.fileName))
				} else {
					tok.SetValue(v)
				}

			case token.TRUE_ID:
				tok.SetValue(true)
			case token.FALSE_ID:
				tok.SetValue(false)

			case token.NULL_ID:
				tok.SetValue(nil)

			case token.IDENTIFIER_ID:
				// identifier's value is variable name.
				tok.SetValue(parser.line[parser.position:pos])
			}

			parser.position = pos
			parser.appendToBuffer(tok)

			return tok, nil
		}
	}
}

func (parser *Parser) HasNext() bool {
	if parser.cursor < len(parser.buffer)-1 {
		return true
	}

	if len(parser.line) == parser.position {
		// Current line finished, fetch next line.
		for {
			if !parser.scanner.Scan() {
				return false
			}
			parser.line = parser.scanner.Text()
			parser.lineNumber++

			// skip empty line and comment
			if len(parser.line) > 0 && !strings.HasPrefix(parser.line, token.COMMENT) {
				parser.position = 0
				return true
			}
		}
	} else {
		return true
	}
}

func (parser *Parser) GetCursor() int {
	return parser.cursor
}

func (parser *Parser) Seek(cursor int) {
	if cursor < parser.cursor {
		if cursor < -1 {
			cursor = -1
		}
		parser.cursor = cursor
	}
}

func (parser *Parser) RollBack(size int) {
	parser.cursor -= size
	if parser.cursor < -1 {
		parser.cursor = -1
	}
}

func (parser *Parser) Commit() {
	parser.cursor = -1
	parser.buffer = []*token.Token{}
}

func (parser *Parser) appendToBuffer(tok *token.Token) {
	parser.cursor++
	parser.buffer = append(parser.buffer, tok)
}
