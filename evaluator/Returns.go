package evaluator

import (
	"EvoScript/lexer"
)

// ReturnStatement will take what it should return and work the args out
func (E *Evaluator) ReturnStatement(returns []lexer.Token) ([]lexer.Token, error) {
	var Sections [][]lexer.Token = make([][]lexer.Token, 1)
	for _, token := range returns {
		if token.Sort == lexer.Comma {
			Sections = append(Sections, make([]lexer.Token, 0))
			continue
		}

		Sections[len(Sections)-1] = append(Sections[len(Sections)-1], token)
	}

	if len(Sections[0]) == 0 {
		return nil, nil
	}

	returns = make([]lexer.Token, 0)
	for _, section := range Sections {
		answer, err := E.args(section, 0)
		if err != nil {
			return nil, err
		}

		returns = append(returns, *answer.memory.(Var).value)
	}

	return returns, nil
}
