package formatter

import (
	"EvoScript/lexer"
)

type Declare struct {
	Methods []DeclareLine // Stores an array of all the line
}

type DeclareLine struct {
	Models            map[int]*declareModel // Stores all the declare models
	WholeSomeFunction bool                  // WholeSomeFunction would be that a function is present inside the arguments when 2 of more declarations are found but not other elements
}

// Extra model used within *Declare
type declareModel struct {
	Model  lexer.Token     // Stores the declaration name
	Locked lexer.TokenType // Stores the declaration locked specification
	Values []lexer.Token   // Stores all the variables given inside the system
}

// Model when functions are used
type FunctionCall struct {
	Path []lexer.Token // Stores the path for the function path
	Args []lexer.Token // Stores the complete leadup for the system
}

// Model used when if statements are found
type ExpressionIF struct {
	Args    []lexer.Token
	Body    []FormattedPacket
	Keyword lexer.Token
}

type FunctionBody struct {
	Keyword    lexer.Token       // Stores the keyword
	ArgsWants  []lexer.Token     // Stores all the arguments needed
	ReturnArgs []lexer.Token     // Stores all the return arguments
	Bodys      []FormattedPacket // Stores the body of the function
}
