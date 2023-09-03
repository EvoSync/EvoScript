package lexer



type Position struct {
	Column int // Column number
	Line   int // Line number
}


// peekNext will peek the next charater
func (lex *Lexer) peekNext() rune {

	// Checks if it can fit inside the line
	if len(lex.source[lex.position.Line]) <= lex.position.Column + 1 {
		return 0
	}

	// Gets the next inline charater of the source entry
	return rune(lex.source[lex.position.Line][lex.position.Column + 1][0])
}