package parse

import "EvoScript/lexer"

const (
	// Stores all the keywords for parsing
	VAR   = "var"
	CONST = "const"

	// Stores all the keywords for parsing
	RETURN = "return"

	// Stores the keyword for IF
	IF = "if"

	// Stores the keyword for functions
	FUNCTION = "func"
)

type DeclareStatement struct {
	BodyTYPED bool // Stores if the declaration statement is a body typed statement (involes a different type of parsing)
	Constant  bool // Stores if the declaration statement is a constant
}

// Stores the functionCall
type FunctionCall struct {
	Args   []lexer.Token // Stores all the arguments
	Leadup []lexer.Token // Stores all the leadup tokens to the first function
}

// Holds information for if statements
type IfStatement struct {
	Expression []Expression // Stores all the information within the system
}

type Expression struct {
	Args    []lexer.Token // Stores the arguments
	Body    []Instruction // Stores all the parsed body
	Keyword lexer.Token   // Stores the keyword like IF, ELSE, ELSEIF
}

type FunctionBody struct {
	Keyword    lexer.Token   // Stores the keyword
	ArgsWants  []lexer.Token // Stores all the arguments needed
	ReturnArgs []lexer.Token // Stores all the return arguments
	Bodys      []Instruction // Stores the body of the function
}
