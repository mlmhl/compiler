package regex

import (
	"testing"
)

func doTestNormalizeRegex(reText string, normalized []int, t *testing.T) {
	if res, err := normalizeRegex(reText); err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		if len(res) != len(normalized) {
			t.Fatalf("Wrong size: wanted %d, got %d", len(normalized), len(res))
		}
		for i, id := range res {
			if normalized[i] != id {
				t.Fatalf("Wrong id(%d): wanted %d, got %d", i, normalized[i], id)
			}
		}
	}
}

func TestNormalizeRegex(t *testing.T) {
	t.Log("Test: normalizeregex ...")

	reText := "(ABC)|(DEF)*"
	normalized := []int{LSP_ID, 65, 66, 67, RSP_ID, CHOICE_ID,
		LSP_ID, 68, 69, 70, RSP_ID, REPETITION_ID}
	doTestNormalizeRegex(reText, normalized, t)

	reText = "if a==b \\|\\| c == d"
	normalized = []int{105, 102, 32, 97, 61, 61, 98, 32, 124,
		124, 32, 99, 32, 61, 61, 32, 100}
	doTestNormalizeRegex(reText, normalized, t)

	if _, err := normalizeRegex("abc\\"); err == nil {
		t.Fatalf("An error is determined to become a correct expression")
	} else {
		targetMsg := "Invalid character at position: 3"
		if err.Error() != targetMsg {
			t.Fatalf("Wrong error message: wanted %s, got %s",
				targetMsg, err.Error())
		}
	}

	t.Log("Passed")
}

func matchCheck(targetPos int, targetGroups []int, pos int,
	groups []int, tag int, t *testing.T) {
	if pos != targetPos {
		t.Fatalf("%d: Wrong matched position, Wanted %d, got %d",
			tag, targetPos, pos)
	}
	if len(targetGroups) != len(groups) {
		t.Fatalf("%d: Wrong matched group size: Wanted %d, got %d",
			tag, len(targetGroups), len(groups))
	}

	groupSets := map[int]bool{}
	for _, group := range targetGroups {
		groupSets[group] = true
	}

	for _, group := range groups {
		if _, ok := groupSets[group]; !ok {
			t.Fatalf("%d: Wrong group id: %d", tag, group)
		}
	}
}

func TestNfaGraph(t *testing.T) {
	t.Log("Test: nfaGraph ...")

	cnt := 0

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(ab|cd)", 1)
		nfa.addRegexExpression("(ef|gh)", 2)

		cnt++ // 1
		pos, groupIds := nfa.match("ab")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 2
		pos, groupIds = nfa.match("gh")
		matchCheck(2, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 3
		pos, groupIds = nfa.match("abc")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(ab*|cd)", 1)
		nfa.addRegexExpression("(ef|g*h)", 2)

		cnt++ // 4
		pos, groupIds := nfa.match("a")
		matchCheck(1, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 5
		pos, groupIds = nfa.match("h")
		matchCheck(1, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 6
		pos, groupIds = nfa.match("abb")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 7
		pos, groupIds = nfa.match("ggh")
		matchCheck(3, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 8
		pos, groupIds = nfa.match("abc")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 9
		pos, groupIds = nfa.match("gih")
		matchCheck(1, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("((1|2|3|4|5|6|7|8|9)(0|1|2|3|4|5|6|7|8|9)*)", 1)
		nfa.addRegexExpression("(\"(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z" +
			"|" + "a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)*\")", 2)

		cnt++ // 10
		pos, groupIds := nfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 11
		pos, groupIds = nfa.match("0123")
		matchCheck(0, []int{}, pos, groupIds, cnt, t)

		cnt++ // 12
		pos, groupIds = nfa.match("\"helloWORLD\"")
		matchCheck(12, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 13
		pos, groupIds = nfa.match("\"hello WORLD\"")
		matchCheck(6, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(\".*\")", 1)
		nfa.addRegexExpression("(123.*)", 2)

		cnt++ // 14
		pos, groupIds := nfa.match("12345abc")
		matchCheck(8, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 15
		pos, groupIds = nfa.match("\"Hello Gdync!\"")
		matchCheck(14, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(1(2)?3)", 1)

		cnt++ // 16
		pos, groupIds := nfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 17
		pos, groupIds = nfa.match("13")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 18
		pos, groupIds = nfa.match("1223")
		matchCheck(2, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(12+3)", 1)

		cnt++ // 19
		pos, groupIds := nfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 20
		pos, groupIds = nfa.match("1223")
		matchCheck(4, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 21
		pos, groupIds = nfa.match("13")
		matchCheck(1, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(\\{)", 1)

		cnt++ // 22
		pos, groupIds := nfa.match("{")
		matchCheck(1, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("((1|2|3|4|5|6|7|8|9)(0|1|2|3|4|5|6|7|8|9)*)", 1)
		nfa.addRegexExpression("((0|1|2|3|4|5|6|7|8|9)+\\.(0|1|2|3|4|5|6|7|8|9)+)", 2)

		cnt++ // 23
		pos, groupIds := nfa.match("12345")
		matchCheck(5, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 24
		pos, groupIds = nfa.match("0.0137")
		matchCheck(6, []int{2}, pos, groupIds, cnt, t)
	}

	t.Log("Passed")
}

func TestDfaGraph(t *testing.T) {
	t.Log("Test: dfaGraph ...")

	cnt := 0

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(ab|cd)", 1)
		nfa.addRegexExpression("(ef|gh)", 2)
		dfa := nfa.toDfa()

		cnt++ // 1
		pos, groupIds := dfa.match("ab")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 2
		pos, groupIds = dfa.match("gh")
		matchCheck(2, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 3
		pos, groupIds = dfa.match("abc")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(ab*|cd)", 1)
		nfa.addRegexExpression("(ef|g*h)", 2)
		dfa := nfa.toDfa()

		cnt++ // 4
		pos, groupIds := dfa.match("a")
		matchCheck(1, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 5
		pos, groupIds = dfa.match("h")
		matchCheck(1, []int{2}, pos, groupIds, cnt, t)


		cnt++ // 6
		pos, groupIds = dfa.match("abb")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 7
		pos, groupIds = dfa.match("ggh")
		matchCheck(3, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 8
		pos, groupIds = dfa.match("abc")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 9
		pos, groupIds = dfa.match("gih")
		matchCheck(1, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("((1|2|3|4|5|6|7|8|9)(0|1|2|3|4|5|6|7|8|9)*)", 1)
		nfa.addRegexExpression("(\"(A|B|C|D|E|F|G|H|I|J|K|L|M|N|O|P|Q|R|S|T|U|V|W|X|Y|Z" +
		"|" + "a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)*\")", 2)
		dfa := nfa.toDfa()

		cnt++ // 10
		pos, groupIds := dfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 11
		pos, groupIds = dfa.match("0123")
		matchCheck(0, []int{}, pos, groupIds, cnt, t)

		cnt++ // 12
		pos, groupIds = dfa.match("\"helloWORLD\"")
		matchCheck(12, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 13
		pos, groupIds = dfa.match("\"hello WORLD\"")
		matchCheck(6, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(\".*\")", 1)
		nfa.addRegexExpression("(123.*)", 2)
		dfa := nfa.toDfa()

		cnt++ // 14
		pos, groupIds := dfa.match("12345abc")
		matchCheck(8, []int{2}, pos, groupIds, cnt, t)

		cnt++ // 15
		pos, groupIds = dfa.match("\"Hello Gdync!\"")
		matchCheck(14, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(1(2)?3)", 1)
		dfa := nfa.toDfa()

		cnt++ // 16
		pos, groupIds := dfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 17
		pos, groupIds = dfa.match("13")
		matchCheck(2, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 18
		pos, groupIds = dfa.match("1223")
		matchCheck(2, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(12+3)", 1)
		dfa := nfa.toDfa()

		cnt++ // 19
		pos, groupIds := dfa.match("123")
		matchCheck(3, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 20
		pos, groupIds = dfa.match("1223")
		matchCheck(4, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 21
		pos, groupIds = dfa.match("13")
		matchCheck(1, []int{}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("(\\{)", 1)
		dfa := nfa.toDfa()

		cnt++ // 22
		pos, groupIds := dfa.match("{")
		matchCheck(1, []int{1}, pos, groupIds, cnt, t)
	}

	{
		nfa := newNfaGraph()
		nfa.addRegexExpression("((1|2|3|4|5|6|7|8|9)(0|1|2|3|4|5|6|7|8|9)*)", 1)
		nfa.addRegexExpression("((0|1|2|3|4|5|6|7|8|9)+\\.(0|1|2|3|4|5|6|7|8|9)+)", 2)
		dfa := nfa.toDfa()

		cnt++ // 23
		pos, groupIds := dfa.match("12345")
		matchCheck(5, []int{1}, pos, groupIds, cnt, t)

		cnt++ // 24
		pos, groupIds = dfa.match("0.0137")
		matchCheck(6, []int{2}, pos, groupIds, cnt, t)
	}

	t.Log("Passed")
}
