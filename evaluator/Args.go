package evaluator

import (
	"errors"
	"fmt"
	"EvoScript/formatter"
	"EvoScript/lexer"
	"strconv"
	"strings"
)

// Stores all basic operators which could exist
var Operators []lexer.TokenType = []lexer.TokenType{
	lexer.Addition, lexer.Subtraction, lexer.Modulus,
	lexer.Division, lexer.Multiplication,
}

// args will work the predicted outcome from the arguments passed
func (eval *Evaluator) args(v []lexer.Token, t lexer.TokenType) (*Pointer, error) {
	var indexedOperators []lexer.Token = make([]lexer.Token, 0)      // all operators
	var indexedValues    [][]lexer.Token = make([][]lexer.Token, 1)	// all values

	// If the token is found within a body we will ignore any operators and just append the value
	var insideBody int = 0

	// Loops through all values
	for _, current := range v {

		if current.Sort == lexer.OpenParenthesis {
			insideBody++
		} else if current.Sort == lexer.CloseParenthesis {
			insideBody--
		}


		if insideBody <= 0 && unwrapTokenArray(Operators, current.Sort) {
			indexedOperators = append(indexedOperators, current)			// Saves into the map
			indexedValues = append(indexedValues, make([]lexer.Token, 0))	// Makes a new element space
			continue
		}


		// Saves into the map
		indexedValues[len(indexedValues)-1] = append(indexedValues[len(indexedValues)-1], current)
	}

	// Ranges through indexedValues
	for index := 0; index < len(indexedValues); index++ {

		switch indexedValues[index][0].Sort {

		case lexer.STRING: // String
			if len(indexedValues[index]) > 1 { // Checks the length of the object array
				return nil, fmt.Errorf("%d:%d STRING type can only contain 1 element per operator", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column)
			}


			if t == 0 && index == 0 {
				t = lexer.STRING
			} else if t != lexer.STRING {
				return nil, fmt.Errorf("%d:%d STRING type cant appear inside a %s type", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column, t.String())
			}

			if indexedValues[index][0].Literal[0] == '"' && indexedValues[index][0].Literal[len(indexedValues[index][0].Literal)-1] == '"' {
				var capture *lexer.Token = &indexedValues[index][0]	// Captures the current input
				capture.Literal = removeSpeach(capture.Literal)		// Removes the speach indicator
				indexedValues[index] = make([]lexer.Token, 0)		// Clears array
				indexedValues[index] = append(indexedValues[index], *capture) // Sets array
			}
		case lexer.INT: // Int
			if len(indexedValues[index]) > 1 { // Checks the length of the object array
				return nil, fmt.Errorf("%d:%d INT type can only contain 1 element per operator", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column)
			}

			if t == 0 && index == 0 {
				t = lexer.INT
			} else if t != lexer.INT {
				return nil, fmt.Errorf("%d:%d INT type cant appear inside a %s type", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column, t.String())
			}
		case lexer.BOOL: // Bool
			if len(indexedValues[index]) > 1 { // Checks the length of the object array
				return nil, fmt.Errorf("%d:%d BOOL type can only contain 1 element per operator", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column)
			}

			if t == 0 && index == 0 {
				t = lexer.BOOL
			} else if t != lexer.BOOL {
				return nil, fmt.Errorf("%d:%d BOOL type cant appear inside a %s type", indexedValues[index][0].Position.Line, indexedValues[index][0].Position.Column, t.String())
			}
		default: // Indentation
			var (
				memory []Pointer = make([]Pointer, 0)
				err    error         = nil
			)

			if indexedValues[index][len(indexedValues[index])-1].Sort == lexer.CloseParenthesis {
			
				var path []lexer.Token = make([]lexer.Token, 0)	// Path to the memory
				var args []lexer.Token = make([]lexer.Token, 0)	// Args used within the function
				var flip bool = false							// Flips the object onto different save

				// Looks through and places into segmants
				for sys := 0; sys < len(indexedValues[index]); sys++ {
					if indexedValues[index][sys].Sort == lexer.OpenParenthesis {
						flip = !flip; continue
					} else if indexedValues[index][sys].Sort == lexer.CloseParenthesis {
						break
					}

					if flip { // Stores the tokens as args
						args = append(args, indexedValues[index][sys])
					} else { // Stores the tokens as leadup
						path = append(path, indexedValues[index][sys])
					}
				}

				// Acts as the main function call working the information
				memory, err = eval.functionCall(&formatter.FunctionCall{Path: path, Args: args})
			} else {
				memory, err = eval.lookupMemory(indexedValues[index])
			}

			if err != nil || len(memory) == 0 {
				if err != nil {
					return nil, errors.New(strconv.Itoa(v[0].Position.Line)+":"+strconv.Itoa(v[0].Position.Column)+" "+err.Error())
				}
				return nil, errors.New("resolver["+strconv.Itoa(v[0].Position.Line)+":"+strconv.Itoa(v[0].Position.Column)+"]: object was not resolved and had been passed into wanted type")
			}

			for _, point := range memory {
				switch point.memory.(type) {

				case Object:
					return &point, nil

				case Var:
					selected := point.memory.(Var)
					if t == 0 && index == 0 {
						t = selected.value.Sort
					}

					indexedValues[index] = make([]lexer.Token, 0)				// Clears the array
					indexedValues[index] = append(indexedValues[index], *selected.value)// Saves the current object into memor
				}
			}
			continue
		}
	}

	switch t {
		case lexer.STRING: // String worker
			return eval.string(indexedValues, indexedOperators)
		case lexer.INT: // Int worker works like maths eval
			return eval.int(indexedValues, indexedOperators)
		case lexer.BOOL: // Boolean indent worker
			if len(indexedValues) > 1 && len(indexedOperators) > 0 {
				return nil, fmt.Errorf("%d:%d boolean must not contain operators", indexedValues[0][0].Position.Line, indexedValues[0][0].Position.Column)
			}

			return eval.allocatedVar("", &indexedValues[0][0])
	}

	return nil, nil
}



// UnwrapTokenArray will unwrap into an array
func unwrapTokenArray(tokens []lexer.TokenType, value lexer.TokenType) bool {

	// Ranges through all the token ids
	for _, element := range tokens {
		if value == element {
			return true
		}
	}

	return false
}

func removeSpeach(t string) string {
	new := strings.Split(t, "")
	new =  new[:len(new)-1]						// Removes last arg
	new =  new[1:]								// Removes first arg
	return strings.Join(new, "")				// Returns complete string
}


// convertTokensToTokenTypeArray will convert the entire array to string
func convertTokensToTokenTypeArray(tokens []Pointer) []string {
	var source []string = make([]string, len(tokens))

	for pos, current := range tokens {
		switch current.memory.(type) {

		case Object:
			source[pos] = current.memory.(Object).module
		case Var:
			source[pos] = current.memory.(Var).value.Sort.String()
		}
	
	}

	return source
}

// convertTokensToTokenTypeArray will convert the entire array to string
func convertTokenTypesToTokenTypeArray(tokens []lexer.TokenType) []string {
	var source []string = make([]string, len(tokens))

	for pos, current := range tokens {
		source[pos] = current.String()
	}

	return source
}


// convertTokensToTokenLiteralArray will convert the entire array to string
func convertTokensToTokenLiteralArray(tokens []lexer.Token) []string {
	var source []string = make([]string, len(tokens))

	for pos, current := range tokens {
		source[pos] = current.Literal
	}

	return source
}
