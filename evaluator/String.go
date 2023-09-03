package evaluator

import (
	"EvoScript/lexer"
	"errors"
	"strconv"
)

// String will join all the tokens and operators into one continous breakpoint token
func (eval *Evaluator) string(values [][]lexer.Token, operators []lexer.Token) (*Pointer, error) {

	if len(values) == 0 { // No values specification inside the string literal
		return nil, errors.New("object which has been given contains no subvalues [NOPOS]")
	}

	if len(operators) == 0 { // Non operators given
		return eval.allocatedVar("", &values[0][0])
	}

	if len(values[0]) > 1 || values[0][0].Sort != lexer.STRING { // Checks for no values or invalid type
		return nil, errors.New("resolver[" + strconv.Itoa(values[0][0].Position.Line) + ":" + strconv.Itoa(values[0][0].Position.Column) + "]: object was not resolved and had been passed into STRING()")
	}

	var object *lexer.Token = &values[0][0] // Base object stored inside here

	// Loops through the values
	for sys := 1; sys < len(values); sys++ {

		if len(values[sys]) > 1 { // variable/object/function was not resolved inside the args
			return nil, errors.New("resolver[" + strconv.Itoa(values[sys][0].Position.Line) + ":" + strconv.Itoa(values[sys][0].Position.Column) + "]: object was not resolved and had been passed into STRING()")
		}

		if operators[sys-1].Sort != lexer.Addition { // Only allows one operator to be completed
			return nil, errors.New("string[" + strconv.Itoa(values[sys][0].Position.Line) + ":" + strconv.Itoa(values[sys][0].Position.Column) + "]: (" + operators[sys-1].Literal + ") unknown operator type inside string statement")
		}

		object.Literal += values[sys][0].Literal
	}

	// Returns the pure token value
	return eval.allocatedVar("", object)
}
