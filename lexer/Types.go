package lexer

import (
	"strings"
)

const (
	// Tokens below can be represented as staticly typed types models
	INT    TokenType = 1 // Integer can be represented as an example like this: 1
	BOOL   TokenType = 2 // Boolean can be represented as an example like this: true, false
	STRING TokenType = 3 // String can be represented as an example like this: "Hello, world!"
	INDENT TokenType = 4 // Indentation can be represented as an example like this: hello
	TEXT   TokenType = 5 // Text can be represented as an example like this: Hello there!

	// Tokens below can be represented as operators
	Question        TokenType = 10 // Question can be represented as an example like this: ?
	Comma           TokenType = 11 // Comma can be represented as an example like this: ,
	Fullstop        TokenType = 12 // Fullstop can be represented as an example like this: .
	SemiColon       TokenType = 13 // Semi-colon can be represented as an example like this: ;
	Colon           TokenType = 14 // Colon can be represented as an example like this: :
	GreaterThan     TokenType = 15 // Greaterthan can be represented as an example like this: >
	GreaterEqual    TokenType = 16 // Greaterequal can be represented as an example like this: >=
	LessThan        TokenType = 17 // Lessthan can be represented as an example like this: <
	LessEqual       TokenType = 18 // Lessequal can be represented as an example like this: <=
	NotEqual        TokenType = 19 // NotEqual can be represented as an example like this: !=
	Bang            TokenType = 20 // Bang can be represented as an example like this: !
	EqualEqual      TokenType = 21 // EqualEqual can be represented as an example like this: ==
	Equal           TokenType = 22 // Equal can be represented as an example like this: =
	Addition        TokenType = 23 // Addition can be represented as an example like this: +
	AdditionOn      TokenType = 24 // AdditionOn can be represented as an example like this: +=
	Addition1       TokenType = 25 // AdditionOn can be represented as an example like this: ++
	Subtraction     TokenType = 26 // Subtraction can be represented as an example like this: -
	Subtraction1    TokenType = 27 // Subtraction can be represented as an example like this: --
	SubtractionFrom TokenType = 28 // SubtractionFrom can be represented as an example like this: -=
	Division        TokenType = 29 // Division can be represented as an example like this: /
	Multiplication  TokenType = 30 // Multiplication can be represented as an example like this: *
	Modulus         TokenType = 31 // Modulus can be represented as an example like this: %
	Dollar			TokenType = 32 // Dollar can be represented as an example like this: $
	AndAnd	        TokenType = 33 // AndAnd can be represented as an example like this: &&
	OrOr			TokenType = 34 // OrOr can be represented as an example like this: ||

	// Tokens below this point can be classified as bodys
	OpenParenthesis    TokenType = 40 // OpenParenthesis can be represented as an example like this: (
	OpenSquareBrackets TokenType = 41 // OpenSquareBrackets can be represented as an example like this: [
	OpenBraces         TokenType = 42 // OpenBraces can be represented as an example like this: {

	CloseParenthesis    TokenType = 43 // CloseParenthesis can be represented as an example like this: )
	CloseSquareBrackets TokenType = 44 // CloseSquareBrackets can be represented as an example like this: ]
	CloseBraces         TokenType = 45 // CloseBraces can be represented as an example like this: }

	OpenBody            TokenType = 46 // OpenBody can be represented as an example like this: <?
	CloseBody           TokenType = 47 // CloseBody can be represented as an example like this: ?>
	Any	        		TokenType = 48 // Any can be represented as an example like this: <any>

	VariadicString		TokenType = 90
	VariadicBool		TokenType = 91
	VariadicInt			TokenType = 92
	VariadicAny			TokenType = 93 // support all types above but randomized and mixed
	Comment	        	TokenType = 48 // support for comments inside the lexer

	INDENTCOG string = "qwertyuiopasdfghjklzxcvbnm_"
	INTCOG    string = "1234567890"

)

var (
	Escapes map[string]string = map[string]string{
		"\\x1b":"\x1b", "\\u001b":"\u001b", "\\033":"\033", 
		"\\r":"\r", "\\n":"\n", "\\a":"\a", 
		"\\b":"\b", "\\t":"\t", "\\v":"\v",
		"\\f":"\f", "\\007":"\007",
	}
)

func (T TokenType) String() string {
	switch T {

	case TEXT:
		return "TEXT"
	case INT:
		return "INT"
	case BOOL:
		return "BOOL"
	case STRING:
		return "STRING"
	case INDENT:
		return "INDENT"
	case VariadicString:
		return "...STRING"
	case VariadicInt:
		return "...INT"
	case VariadicBool:
		return "...BOOL"
	case VariadicAny:
		return "...ANY"
	case OpenBody:
		return "OB"
	case CloseBody:
		return "CB"
	default:
		return "EOF"
	}
}

// isINT checks for an int cog
func IsINT(s rune) bool {
	return strings.Contains(INTCOG, strings.ToLower(string(s)))
}

// isINDENT checks for an indent cog
func IsINDENT(s rune) bool {
	return strings.Contains(INDENTCOG, strings.ToLower(string(s)))
}

// WorkINDENT will workout the initial length of the token
// this will format the token into the offical segment length until its voided
func (lex *Lexer) WorkINDENT() (*Token, error) {

	// Makes the new token module
	var token = NewToken("", INDENT, *lex.position)

	for position := lex.position.Column; position < len(lex.source[lex.position.Line]); position++ {

		// Checks the type for a exit code
		if position == 0 && !IsINDENT(rune(lex.source[lex.position.Line][position][0])) || !IsINDENT(rune(lex.source[lex.position.Line][position][0])) && !IsINT(rune(lex.source[lex.position.Line][position][0])) && lex.source[lex.position.Line][position][0] != '_' {
			return token, nil // Returns the token
		} else {
			token.Literal += lex.source[lex.position.Line][position] // Addons on length
		}
	}

	// Returns the token
	return token, nil
}

// WorkINT will workout the initial length of the token
// this will format the token into the offical segment length until its voided
func (lex *Lexer) WorkINT() (*Token, error) {

	// Makes the new token module
	var token = NewToken("", INT, *lex.position)

	for position := lex.position.Column; position < len(lex.source[lex.position.Line]); position++ {

		// Checks the type for a exit code
		if !IsINT(rune(lex.source[lex.position.Line][position][0])) {
			return token, nil // Returns the token
		} else {
			token.Literal += lex.source[lex.position.Line][position] // Addons on length
		}
	}

	// Returns the token
	return token, nil
}

// WorkSTRING will workout the initial length of the token
// this will format the token into the offical segment length until its voided
func (lex *Lexer) WorkSTRING() (*Token, error) {

	// Makes the new token module
	var token = NewToken("\"", STRING, *lex.position)

	for position := lex.position.Column + 1; position < len(lex.source[lex.position.Line]); position++ {

		// Checks the type for a exit code
		if rune(lex.source[lex.position.Line][position][0]) == '"' {
			token.Literal += "\"" // Adds the close escape character
			return token, nil // Returns the token
		} else {
			token.Literal += lex.source[lex.position.Line][position] // Addons on length
		}
	}

	for escape, pure := range Escapes {
		token.Literal = strings.ReplaceAll(token.Literal, escape, pure)
	}

	// Returns the token
	return token, nil
}

func AnsiUtil(src string) string {
	for escape := range Escapes {
		src = strings.ReplaceAll(src, escape, Escapes[escape])
	}
	return src
}