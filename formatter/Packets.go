package formatter

import (
	"EvoScript/lexer"
)

// Stores the formatted node type
type FormatNode int

// Stores the formatted node output
type FormattedPacket struct {
	Node             FormatNode     // Stores the formatted nodeType
	DeclareStatement *Declare       // Holds information if the statement is a declaration
	CallStatement    *FunctionCall  // Holds information if the statement is a call statement
	IFCall           []ExpressionIF // Holds information if the call statement is if statement
	Returns          []lexer.Token  // Holds information if the statement is a return statement
	Axis             []lexer.Token  // Holds all the token axis attributes
	FunctionBody     *FunctionBody
}
