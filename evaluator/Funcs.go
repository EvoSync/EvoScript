package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
)

var inCallArgs []lexer.TokenType = []lexer.TokenType{
	lexer.Comma,
}

// functionCall will try to register and identify the function state
func (eval *Evaluator) functionCall(call *formatter.FunctionCall) ([]Pointer, error) {
	var Values [][]lexer.Token = make([][]lexer.Token, 1) // all operators

	// If the token is found within a body we will ignore any operators and just append the value
	var insideBody int = 0

	// Loops through all values
	for _, current := range call.Args {
		if current.Sort == lexer.OpenParenthesis {
			insideBody++
		} else if current.Sort == lexer.CloseParenthesis {
			insideBody--
		}

		if insideBody <= 0 && current.Sort == lexer.Comma {
			Values = append(Values, make([]lexer.Token, 0)) // Makes a new element space
			continue
		}

		// Saves into the map
		Values[len(Values)-1] = append(Values[len(Values)-1], current)
	}

	// Stores all the different arg segmants
	var inbetween [][]lexer.Token = make([][]lexer.Token, 1)
	var inside int = 0
	for _, arg := range call.Args { // Loops through all the arguments

		if arg.Sort == lexer.OpenParenthesis {
			inside++
		} else if arg.Sort == lexer.CloseParenthesis {
			inside--
		}

		if inside <= 0 && unwrapTokenArray(inCallArgs, arg.Sort) {
			inbetween = append(inbetween, make([]lexer.Token, 0))
			continue
		}

		// Appends the current token into the array
		inbetween[len(inbetween)-1] = append(inbetween[len(inbetween)-1], arg)
	}

	var pure []Pointer = make([]Pointer, 0) // Stores the information worked pure args
	for _, pureArg := range inbetween {     // Loops through all the arguments
		if len(pureArg) == 0 {
			continue
		}

		// works the args out
		current, err := eval.args(pureArg, 0)
		if err != nil {
			return nil, err
		}

		// Appends the pure value
		pure = append(pure, *current)
	}

	// Looks up and executes the function
	return eval.lookupMemory(call.Path, pure...)
}
