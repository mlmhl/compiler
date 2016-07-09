package regex

import (
	"fmt"
)

// Convert the regex expression to a int array, each number in
// this array represent the unique id of corresponding character
// in the original regex expression.
func normalizeRegex(reText string) (normalized []int, err error) {
	if len(reText) == 0 {
		return nil, nil
	}
	if reText[0] == '*' || reText[0] == '|' {
		return nil, fmt.Errorf("Invalid character at position: 0")
	}


	normalized = make([]int, 0, len(reText))

	pos := 0
	lspCnt := 0
	length := len(reText)
	for pos < length {
		c := reText[pos]
		if c == '\\' {
			if pos == len(reText) - 1 {
				return nil, fmt.Errorf("Invalid character at position: %d", pos)
			}
			normalized = append(normalized, (int)(reText[pos + 1]))
			pos += 2
		} else {
			if id, ok := mataSymbolId[(int)(c)]; ok {
				if id == LSP_ID {
					lspCnt++
				} else if id == RSP_ID {
					if lspCnt == 0 {
						return nil, fmt.Errorf("Mismatched right parentheses: %d", pos)
					}
					lspCnt--
				}
				normalized = append(normalized, id)
			} else {
				normalized = append(normalized, (int)(c))
			}
			pos++
		}
	}

	if lspCnt > 0 {
		return nil, fmt.Errorf("Mismatched left parentheses: 0")
	} else {
		return normalized, nil
	}
}

func isMataSymbol(c int) bool {
	_, ok := mataSymbolSet[c]
	return ok
}