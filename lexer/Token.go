package lexer


// TokenType storage container
type TokenType int


type Token struct {
	Literal	 	string		// literal format
	Sort		TokenType	// acts as tokenType
	Position    *Position   // position of the token
}


// Makes the new token to be added onto the array of tokens
func NewToken(Literal string, Type TokenType, Position Position) *Token {
	return &Token{
		Literal: 	Literal,		// Literal format
		Sort: 		Type,			// Token type
		Position:   &Position,		// Token position
	}
}
