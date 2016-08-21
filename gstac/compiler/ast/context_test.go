package ast

import (
	"testing"
	"strings"
)

func TestSymbolLIst(t *testing.T) {
	symbolList := newSymbolList()

	symbols := []string{
		"Hello",
		"World",
		"Compiler",
		"parser",
		"analyzer",
	}

	t.Log("Test: SymbolList...")

	for _, symbol := range(symbols) {
		symbolList.Put(symbol)
	}
	for target, symbol := range(symbols) {
		i := symbolList.Get(symbol)
		if i != target {
			t.Fatalf("Wrong index(%s): Wanted %d, got %d", symbol, target, i)
		}
	}

	targetCodeByte:= "{Hello:0,World:1,Compiler:2,parser:3,analyzer:4}"
	codeByte := string(symbolList.Encode())
	{
		symbols := strings.Split(codeByte, ",")
		newSymbols := make([]string, len(symbols))
		for _, symbol := range(symbols) {
			switch {
			case strings.Contains(symbol, "0"):
				newSymbols[0] = symbol
			case strings.Contains(symbol, "1"):
				newSymbols[1] = symbol
			case strings.Contains(symbol, "2"):
				newSymbols[2] = symbol
			case strings.Contains(symbol, "3"):
				newSymbols[3] = symbol
			case strings.Contains(symbol, "4"):
				newSymbols[4] = symbol
			}
		}
		codeByte = strings.Join(newSymbols, ",")
	}

	if targetCodeByte != codeByte {
		t.Fatalf("Wrong encode: Wanted %s, got %s", targetCodeByte, codeByte)
	}

	t.Log("Passed...")
}