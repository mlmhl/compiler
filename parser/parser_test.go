package parser

import (
	"testing"

	"github.com/mlmhl/compiler/token"
	"github.com/mlmhl/compiler/common"
)

func TestParser(t *testing.T) {
	t.Log("Test: Parser ...")

	fileName := "test"
	parser := NewParser()
	parser.Parse(fileName)

	tokens := []*token.Token{
		token.NewToken(common.NewLocation(1, 0,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(1, 2,
			fileName)).SetType(token.ASSIGN_ID),
		token.NewToken(common.NewLocation(1, 4,
			fileName)).SetType(token.INTEGER_ID).SetValue(int64(5)),
		token.NewToken(common.NewLocation(3, 0,
			fileName)).SetType(token.WHILE_ID),
		token.NewToken(common.NewLocation(3, 6,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(3, 8,
			fileName)).SetType(token.GT_ID),
		token.NewToken(common.NewLocation(3, 10,
			fileName)).SetType(token.INTEGER_ID).SetValue(int64(0)),
		token.NewToken(common.NewLocation(3, 12,
			fileName)).SetType(token.LLP_ID),
		token.NewToken(common.NewLocation(4, 4,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(4, 6,
			fileName)).SetType(token.ASSIGN_ID),
		token.NewToken(common.NewLocation(4, 8,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(4, 10,
			fileName)).SetType(token.SUBTRACT_ID),
		token.NewToken(common.NewLocation(4, 12,
			fileName)).SetType(token.INTEGER_ID).SetValue(int64(1)),
		token.NewToken(common.NewLocation(5, 0,
			fileName)).SetType(token.RLP_ID),
		token.NewToken(common.NewLocation(-1, -1,
			fileName)).SetType(token.FINISHED_ID),
	}

	for i, target := range tokens {
		if tok, err := parser.Next(); err != nil {
			t.Fatalf("Parser error: " + err.GetMessage())
		} else {
			if !tok.Equal(target) {
				t.Fatalf("Wrong token(%d), Wanted %v, got %v", i, target, tok)
			}
		}
	}

	t.Log("Passed")
}
