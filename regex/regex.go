package regex

import (
	"fmt"
)

type Regex struct {
	nfa *nfaGraph
	dfa *dfaGraph
}

func NewRegex() *Regex {
	return &Regex{
		nfa: newNfaGraph(),
		dfa: nil,
	}
}

func (regex *Regex) AddRegexExpression(reText string, groupId int) error {
	return regex.nfa.addRegexExpression(reText, groupId)
}

func (regex *Regex) AddRegexExpressions(reTexts []string, groupIds []int) error {
	if len(reTexts) != len(groupIds) {
		return fmt.Errorf("Invalid parameter size")
	}

	for i, reText := range reTexts {
		if err := regex.AddRegexExpression(reText, groupIds[i]); err != nil {
			return fmt.Errorf("Can't parse %s, %s", reText, err.Error())
		}
	}
	return nil
}

func (regex *Regex) Compile() {
	regex.dfa = regex.nfa.toDfa()
	regex.nfa = nil
}

func (regex *Regex) Match(reText string) (int, []int) {
	return regex.dfa.match(reText)
}
