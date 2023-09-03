package lexer

import (
	"strings"
	"unicode"
)


type Lexer struct {
	source   [][]string // Stores the source array for the lexer
	position *Position  // Stores the position
	Token    []Token    // Stores the array of pack tokens
	capture  bool
}

// NewLexer creates a new lexer structure
func NewLexer(source, split string, capture bool, body ...bool) *Lexer {
	lex := new(Lexer)

	// Ranges through the array of source split
	for _, line := range strings.Split(source, split) {
		lex.source = append(lex.source, strings.Split(line, ""))
	}

	lex.position = new(Position)
	lex.Token    = make([]Token, 0)
	lex.capture  = !capture
	

	// Returns the object
	return lex
}

// Start will run the lexer on the source field
func (lex *Lexer) Start() error {
	for linePosition, lineContains := range lex.source {
		lex.position.Line = linePosition


		if !lex.capture && !strings.Contains(strings.Join(lineContains, ""), "<kbm>") {
			var literal string = strings.Join(lineContains, "")
			if linePosition + 1 != len(lex.source) {
				literal += "\r\n"
			}
			lex.Token = append(lex.Token, *NewToken(AnsiUtil(literal), TEXT, *lex.position))
			continue
		}

		for column := 0; column < len(lineContains); column++ {
			lex.position.Column = column

			if !lex.capture { // Detects there will be an open statement inline
				var capture string = ""

				// Ranges through gets the tags within the network locations
				for column = lex.position.Column; column < len(strings.Split(strings.Join(lineContains, ""), "")); column++ {
					lex.position.Column = column
					if strings.Contains(strings.Join(lineContains, "")[column:], "<kbm>") && lex.isTag(column, strings.Join(lineContains, ""), "<kbm>") {
						lex.capture = !lex.capture; break
					}

					capture += strings.Split(strings.Join(lineContains, ""), "")[column]
				}

				if linePosition + 1 != len(lex.source) && column + 1 == len(lineContains) {
					capture += "\r\n"
				}
			
				column--
				lex.Token = append(lex.Token, *NewToken(AnsiUtil(capture), TEXT, *lex.position))
				continue
			}

			Token, err := lex.tokenize(rune(lineContains[column][0]))
			if err != nil {
				return err
			}

			if Token == nil { // blank token detected
				continue
			}

			column += len(Token.Literal) - 1
			if Token.Sort == Comment{continue}
			lex.Token = append(lex.Token, *Token)
		}
	}

	return nil
}

// tokenize will tokenize the current source column on line
// this will return out an output of the token and an error if one is present
func (lex *Lexer) tokenize(runCharater rune) (*Token, error) {

	switch runCharater {

	case '?': // Question mark
		return NewToken("?", Question, *lex.position), nil

	case ',': // Comma
		return NewToken(",", Comma, *lex.position), nil

	case '.': // Fullstop
		return NewToken(".", Fullstop, *lex.position), nil

	case ';': // SemiColon
		return NewToken(";", SemiColon, *lex.position), nil

	case ':': // Colon
		return NewToken(":", Colon, *lex.position), nil

	case '>': // Greater symbolic
		if lex.peekNext() == '=' {
			return NewToken(">=", GreaterEqual, *lex.position), nil
		}

		return NewToken("<", GreaterThan, *lex.position), nil

	case '<': // Lesser symbolic
		if lex.peekNext() == '=' {
			return NewToken("<=", LessEqual, *lex.position), nil
		} else if lex.isTag(lex.position.Column, strings.Join(lex.source[lex.position.Line], ""), "<kbm>") {
			lex.capture = true; return NewToken("<kbm>", OpenBody, *lex.position), nil
		} else if lex.isTag(lex.position.Column, strings.Join(lex.source[lex.position.Line], ""), "</kbm>") {
			lex.capture = false; return NewToken("</kbm>", CloseBody, *lex.position), nil
		}

		return NewToken("<", LessThan, *lex.position), nil

	case '!': // Exclamation mark
		if lex.peekNext() == '=' {
			return NewToken("!=", NotEqual, *lex.position), nil
		}

		return NewToken("!", Equal, *lex.position), nil

	case '=': // Equal
		if lex.peekNext() == '=' {
			return NewToken("==", EqualEqual, *lex.position), nil
		}

		return NewToken("=", Equal, *lex.position), nil

	case '+': // Addition
		if lex.peekNext() == '=' {
			return NewToken("+=", AdditionOn, *lex.position), nil
		} else if lex.peekNext() == '+' {
			return NewToken("++", Addition1, *lex.position), nil
		}

		return NewToken("+", Addition, *lex.position), nil

	case '-': // Subtraction
		if lex.peekNext() == '=' {
			return NewToken("-=", SubtractionFrom, *lex.position), nil
		} else if lex.peekNext() == '-' {
			return NewToken("--", Subtraction1, *lex.position), nil
		}

		return NewToken("-", Subtraction, *lex.position), nil

	case '#':
		return NewToken(strings.Join(lex.source[lex.position.Line], "")[lex.position.Column:], Comment, *lex.position), nil

	case '/': // Division
		if lex.peekNext() == '/' {
			return NewToken(strings.Join(lex.source[lex.position.Line], "")[lex.position.Column:], Comment, *lex.position), nil
		} else if lex.peekNext() == '*' {
			var capture string = ""

			for position := lex.position.Column; position < len(lex.source[lex.position.Line]); position++ {
				if strings.Join(lex.source[lex.position.Line], "")[position] == '*' && position + 1 < len(lex.source[lex.position.Line]) && strings.Join(lex.source[lex.position.Line], "")[position+1] == '/' {
					break
				}

				capture += strings.Join(lex.source[lex.position.Line], "")
			}

			return NewToken(capture, Comment, *lex.position), nil
		}
		return NewToken("/", Division, *lex.position), nil

	case '*': // Multiplication
		return NewToken("*", Multiplication, *lex.position), nil

	case '%': // Modulus
		return NewToken("%", Modulus, *lex.position), nil

		case '(': // OpenParenthesis 
		return NewToken("(", OpenParenthesis, *lex.position), nil

	case '[': // OpenSquareBrackets
		return NewToken("[", OpenSquareBrackets, *lex.position), nil

	case '{': // OpenBraces
		return NewToken("{", OpenBraces, *lex.position), nil

	case ')': // CloseParenthesis
		return NewToken(")", CloseParenthesis, *lex.position), nil

	case ']': // CloseSquareBrackets
		return NewToken("]", CloseSquareBrackets, *lex.position), nil

	case '}': // CloseBraces
		return NewToken("}", CloseBraces, *lex.position), nil

	case '$': // Dollar
		return NewToken("$", Dollar, *lex.position), nil

	case '&':
		if lex.peekNext() == '&' {
			return NewToken("&&", AndAnd, *lex.position), nil
		}

	case '|':
		if lex.peekNext() == '|' {
			return NewToken("||", OrOr, *lex.position), nil
		}

	case '"': // String
		return lex.WorkSTRING() // string worker

	default: // Indent, int & space
		if IsINDENT(runCharater) {
			indentaion, err := lex.WorkINDENT() // 	indent worker
			if err != nil {
				return nil, err
			}

			// Checks for boolean type
			if indentaion.Literal == "true" || indentaion.Literal == "false" {
				indentaion.Sort = BOOL // Sets type to boolean
			}

			return indentaion, nil
		} else if IsINT(runCharater) {
			return lex.WorkINT() // int worker
		} else if unicode.IsSpace(runCharater) {
			return nil, nil // blank space
		}
	}

	return nil, ErrNotImplemented
}

func (l *Lexer) isTag(pos int, line, tag string) bool {
	var ticks int = 0
	for position, charater := range strings.Split(tag, "") {

		if position + pos > len(line) {
			return false
		} else if charater == strings.Split(line, "")[position+pos] {
			ticks++
		}
	}
	
	return ticks == len(tag)

}