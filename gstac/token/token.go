package token

import (
	"github.com/mlmhl/compiler/common"
)

func GetDescription(typ int) string {
	return descriptions[typ]
}

func IsAssignOperator(typ int) bool {
	return typ == ASSIGN_ID || typ == ADD_ASSIGN_ID ||
		typ == MUL_ASSIGN_ID || typ == DIV_ASSIGN_ID || typ == MOD_ASSIGN_ID
}

type Token struct {
	typ      int
	value    interface{}
	location *common.Location
}

func NewToken(location *common.Location) *Token {
	return &Token{
		typ:      UNKNOWN,
		value:    nil,
		location: location,
	}
}

func (token *Token) GetType() int {
	return token.typ
}

func (token *Token) GetValue() interface{} {
	return token.value
}

func (token *Token) GetLocation() *common.Location {
	return token.location
}

func (token *Token) SetType(typ int) *Token {
	token.typ = typ
	return token
}

func (token *Token) SetValue(value interface{}) *Token {
	token.value = value
	return token
}

// for test
func (token *Token) Equal(other *Token) bool {
	return token.typ == other.typ && token.value == other.value &&
		token.location.Equal(other.location)
}
