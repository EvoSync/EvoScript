package evaluator

import (
	"EvoScript/lexer"
	"errors"
	"strconv"
)

// int will act as the main maths interpreter inside EvoScripts interpreter
func (eval *Evaluator) int(values [][]lexer.Token, operators []lexer.Token) (*Pointer, error) {

	if len(values) == 0 { // No values specification inside the int literal
		return nil, errors.New("object which has been given contains no subvalues [NOPOS]")
	}

	if len(operators) == 0 { // Non operators given
		return eval.allocatedVar("", &values[0][0])
	}

	if len(values[0]) > 1 || values[0][0].Sort != lexer.INT { // Checks for no values or invalid type
		return nil, errors.New("resolver[" + strconv.Itoa(values[0][0].Position.Line) + ":" + strconv.Itoa(values[0][0].Position.Column) + "]: object was not resolved and had been passed into INT()")
	}

	var object *lexer.Token = &values[0][0]  // Base object stored inside here
	for sys := 1; sys < len(values); sys++ { // Loops through all the insider objects

		if len(values[sys]) > 1 { // variable/object/function was not resolved inside the args
			return nil, errors.New("resolver[" + strconv.Itoa(values[sys][0].Position.Line) + ":" + strconv.Itoa(values[sys][0].Position.Column) + "]: object was not resolved and had been passed into INT()")
		}

		// Converts the base object into int format
		convertedBase, err := strconv.Atoi(object.Literal)
		if err != nil {
			return nil, err
		}

		// Converts the current object into an int format
		convertedItem, err := strconv.Atoi(values[sys][0].Literal)
		if err != nil {
			return nil, err
		}

		switch operators[sys-1].Sort {
		case lexer.Addition: // Addition
			convertedBase += convertedItem
		case lexer.Subtraction: // Subtraction
			convertedBase -= convertedItem
		case lexer.Multiplication: // Multiplication
			convertedBase *= convertedItem
		case lexer.Division: // Division
			convertedBase /= convertedItem
		case lexer.Modulus: // Modulus
			convertedBase %= convertedItem
		}

		// Sets the new bounds of the item
		object.Literal = strconv.Itoa(convertedBase)
	}

	object.Sort = lexer.INT
	return eval.allocatedVar("", object)
}
