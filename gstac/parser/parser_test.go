package parser

import (
	"testing"

	"github.com/mlmhl/compiler/gstac/token"
	"github.com/mlmhl/compiler/common"
)

func TestParser(t *testing.T) {
	t.Log("Test: Parser ...")

	fileName := "test"
	parser := NewParser()
	parser.Parse(fileName)

	tokens := []*token.Token{
		token.NewToken(common.NewLocation(1, 0,
		fileName)).SetType(token.INTEGER_TYPE_ID),
		token.NewToken(common.NewLocation(1, 4,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(1, 6,
			fileName)).SetType(token.ASSIGN_ID),
		token.NewToken(common.NewLocation(1, 8,
			fileName)).SetType(token.INTEGER_VALUE_ID).SetValue(int64(5)),
		token.NewToken(common.NewLocation(5, 0,
			fileName)).SetType(token.WHILE_ID),
		token.NewToken(common.NewLocation(5, 6,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(5, 8,
			fileName)).SetType(token.GT_ID),
		token.NewToken(common.NewLocation(5, 10,
			fileName)).SetType(token.INTEGER_VALUE_ID).SetValue(int64(0)),
		token.NewToken(common.NewLocation(5, 12,
			fileName)).SetType(token.LLP_ID),
		token.NewToken(common.NewLocation(6, 4,
			fileName)).SetType(token.IDENTIFIER_ID).SetValue("i"),
		token.NewToken(common.NewLocation(6, 5,
			fileName)).SetType(token.DECREMENT_ID),
		token.NewToken(common.NewLocation(7, 0,
			fileName)).SetType(token.RLP_ID),
		token.NewToken(common.NewLocation(-1, -1,
			fileName)).SetType(token.FINISHED_ID),
	}

	for i, target := range tokens {
		if tok, err := parser.Next(); err != nil {
			t.Fatalf("Parser error: " + err.GetMessage())
		} else {
			if !tok.Equal(target) {
				t.Fatalf("Wrong token(%d), Wanted %s(%v), got %s(%v)", i,
					token.GetDescription(target.GetType()), target.GetLocation(),
					token.GetDescription(tok.GetType()), tok.GetLocation())
			}
		}
	}

	// test for cursor
	backup := 3
	parser.RollBack(backup)
	for i := len(tokens) - backup; i < len(tokens); i++ {
		target := tokens[i]
		if tok, err := parser.Next(); err != nil {
			t.Fatalf("Parser error: " + err.GetMessage())
		} else {
			if !tok.Equal(target) {
				t.Fatalf("Wrong token(%d), Wanted %s(%v), got %s(%v)", i,
					token.GetDescription(target.GetType()), target.GetLocation(),
					token.GetDescription(tok.GetType()), tok.GetLocation())
			}
		}
	}

	parser.Commit()
	cursor := parser.GetCursor()
	if cursor != -1 {
		t.Fatalf("Wrong curosr: Wanted -1, got %d", cursor)
	}

	t.Log("Passed")
}