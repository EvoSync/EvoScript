package parse

import (
	"EvoScript/lexer"
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	CurrentPosition   int
	Tokens            []lexer.Token
	Lexer             *lexer.Lexer
	GainedExpressions []Instruction
}

func NewParser(Lex *lexer.Lexer) *Parser {
	return &Parser{
		Lexer:             Lex,
		Tokens:            Lex.Token,
		CurrentPosition:   0, // starts at 0
		GainedExpressions: make([]Instruction, 0),
	}
}

// Start will run the parser from the lexers output
func (Parser *Parser) Start() error {

	// Loops through the tokens
	for position := Parser.CurrentPosition; position < len(Parser.Tokens); position++ {
		Parser.CurrentPosition = position // Sets the position

		// Tries to parse the expression
		instruction, err := Parser.ParseExpression()
		if err != nil {
			return err
		}

		// Unknown reason for a nil pointer results in us looping again
		if instruction == nil {
			continue
		}

		position += len(instruction.TokensAxis) - 1 // Allows us to skip through the tokens inside a instruction

		// Appends the instruction gained on to the GainedExpressions array
		Parser.GainedExpressions = append(Parser.GainedExpressions, *instruction)
	}

	return nil
}

func (Parser *Parser) ParseExpression() (*Instruction, error) {

	// Detects for tokens as text types within the system
	if Parser.Tokens[Parser.CurrentPosition].Sort == lexer.TEXT {
		var new *Instruction = new(Instruction)           // Makes new structure
		new.NodeType = 6                                  // Sets nodetype
		new.Text = &Parser.Tokens[Parser.CurrentPosition] // Sets text output
		new.TokensAxis = append(new.TokensAxis, *new.Text)
		return new, nil
	} else if Parser.Tokens[Parser.CurrentPosition].Sort == lexer.OpenBody || Parser.Tokens[Parser.CurrentPosition].Sort == lexer.CloseBody {
		return nil, nil
	}

	// implements the keyword system inside the parser
	switch Parser.Tokens[Parser.CurrentPosition].Literal {

	case FUNCTION:
		return Parser.parseFunction()
	case IF:
		return Parser.parseIF()

	case VAR, CONST: // Keyword for declaration statements
		return Parser.declarationStatement()

	case RETURN: // Keyword for return statements
		return Parser.returnStatement()

	default:
		// Onlys accepts cerain type
		if Parser.Tokens[Parser.CurrentPosition].Sort != lexer.INDENT && Parser.Tokens[Parser.CurrentPosition].Sort != lexer.Dollar { // Checks for only the system
			return nil, errors.New(Parser.Tokens[Parser.CurrentPosition].Sort.String() + " type given can't enter the parser as position: " + strconv.Itoa(Parser.Tokens[Parser.CurrentPosition].Position.Line+1) + ":" + strconv.Itoa(Parser.Tokens[Parser.CurrentPosition].Position.Column))
		}

		// Variable involded statement
		if Parser.Tokens[Parser.CurrentPosition].Sort == lexer.Dollar {
			return Parser.expressionStatement()
		} else {
			// Gets all the inline position
			given := Parser.inlineBodyApplication(Parser.CurrentPosition, make([]lexer.Token, 0), lexer.SemiColon)
			analyse := 0

			// Ranges through the tokens
			for token := range given {

				if given[token].Sort == lexer.Dollar {
					analyse = 1
					break // Variable caught within an object called
				} else if given[token].Sort == lexer.OpenParenthesis || given[token].Sort == lexer.CloseParenthesis {
					analyse = 2
					break // Function caught within an object called
				}
			}

			switch analyse { // Counter acts the analyse

			case 0:
				return &Instruction{NodeType: 5, TokensAxis: given}, nil
			case 1: // Variable caught within the network
				return Parser.expressionStatement()
			case 2: // Function caught within the network
				return Parser.functionStatement()
			}
		}
	}

	return nil, nil
}

// parseFunction will parse the information needed for functions
func (Parser *Parser) parseFunction() (*Instruction, error) {
	var instruction *Instruction = new(Instruction)
	instruction.TokensAxis = make([]lexer.Token, 0)
	instruction.NodeType = 8

	instruction.FunctionCreation = new(FunctionBody)

	if Parser.CurrentPosition+5 >= len(Parser.Tokens) {
		return nil, fmt.Errorf("%d:%d missing core statements in function declaration", Parser.Tokens[Parser.CurrentPosition].Position.Line, Parser.Tokens[Parser.CurrentPosition].Position.Column)
	}

	if Parser.Tokens[Parser.CurrentPosition+1].Sort != lexer.INDENT {
		return nil, fmt.Errorf("%d:%d keyword must only be an indent", Parser.Tokens[Parser.CurrentPosition].Position.Line, Parser.Tokens[Parser.CurrentPosition].Position.Column)
	}

	// Saves into the token axis
	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition:Parser.CurrentPosition+2]...)
	instruction.FunctionCreation.Keyword = Parser.Tokens[Parser.CurrentPosition+1]

	if Parser.Tokens[Parser.CurrentPosition+2].Sort != lexer.OpenParenthesis {
		return nil, fmt.Errorf("%d:%d missing open parathesis where `%s` stands", Parser.Tokens[Parser.CurrentPosition].Position.Line, Parser.Tokens[Parser.CurrentPosition].Position.Column, Parser.Tokens[Parser.CurrentPosition+2].Literal)
	}

	Args, err := Parser.multiLineBodyApplication(Parser.CurrentPosition+2, lexer.OpenParenthesis, lexer.CloseParenthesis, make([]lexer.Token, 0))
	if err != nil {
		return nil, err
	}

	instruction.FunctionCreation.ArgsWants = Args
	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition+2:Parser.CurrentPosition+4+len(Args)]...)
	Parser.CurrentPosition += len(instruction.TokensAxis)

	if Parser.Tokens[Parser.CurrentPosition].Sort == lexer.OpenParenthesis {
		ReturnsArgs, err := Parser.multiLineBodyApplication(Parser.CurrentPosition+1, lexer.OpenParenthesis, lexer.CloseParenthesis, make([]lexer.Token, 0))
		if err != nil {
			return nil, err
		}

		instruction.FunctionCreation.ReturnArgs = ReturnsArgs
		instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition:Parser.CurrentPosition+2+len(ReturnsArgs)]...)
		Parser.CurrentPosition += 2 + len(ReturnsArgs)
	} else if Parser.Tokens[Parser.CurrentPosition].Sort == lexer.INDENT {
		instruction.FunctionCreation.ReturnArgs = []lexer.Token{Parser.Tokens[Parser.CurrentPosition]}
		instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition])
		Parser.CurrentPosition++
	} else if Parser.Tokens[Parser.CurrentPosition].Sort != lexer.OpenBraces {
		return nil, fmt.Errorf("%d:%d return args could not be parsed", Parser.Tokens[Parser.CurrentPosition].Position.Line, Parser.Tokens[Parser.CurrentPosition].Position.Column)
	}

	if Parser.Tokens[Parser.CurrentPosition].Sort != lexer.OpenBraces {
		return nil, fmt.Errorf("%d:%d missing open brace where `%s` stands", Parser.Tokens[Parser.CurrentPosition].Position.Line, Parser.Tokens[Parser.CurrentPosition].Position.Column, Parser.Tokens[Parser.CurrentPosition].Literal)
	}

	ParserBody, err := Parser.multiLineBodyApplication(Parser.CurrentPosition+1, lexer.OpenBraces, lexer.CloseBraces, make([]lexer.Token, 0))
	if err != nil {
		return nil, err
	}

	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition])
	instruction.TokensAxis = append(instruction.TokensAxis, ParserBody...)
	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition+1+len(ParserBody)])
	temp := NewParser(tempLexer(ParserBody))
	if err := temp.Start(); err != nil {
		return nil, err
	}

	instruction.FunctionCreation.Bodys = temp.GainedExpressions
	return instruction, nil
}

// ParseIF will parse the if statement detected
func (Parser *Parser) parseIF() (*Instruction, error) {
	var instruction *Instruction = new(Instruction)
	instruction.IF = new(IfStatement)
	instruction.NodeType = 7
	instruction.IF.Expression = make([]Expression, 0)

	for system := Parser.CurrentPosition; system < len(Parser.Tokens); system++ {
		instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system])

		switch Parser.Tokens[system].Literal {

		case "if", "elif":
			if len(Parser.Tokens) > system+2 && Parser.Tokens[system+1].Sort != lexer.OpenParenthesis {
				return nil, fmt.Errorf("%d:%d if requires ( and not %s", Parser.Tokens[system].Position.Line, Parser.Tokens[system].Position.Column, Parser.Tokens[system+1].Literal)
			}

			var new *Expression = new(Expression)
			new.Keyword = Parser.Tokens[system]
			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system+1]) // Adds the information
			args, err := Parser.multiLineBodyApplication(system+2, lexer.OpenParenthesis, lexer.CloseParenthesis, make([]lexer.Token, 0))
			if err != nil {
				return nil, err
			}

			new.Args = args
			system += len(args) + 2
			instruction.TokensAxis = append(instruction.TokensAxis, args...)
			new.Args = args
			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system])
			if len(Parser.Tokens) <= system+2 {
				return nil, fmt.Errorf("%d:%d if requires {", Parser.Tokens[system].Position.Line, Parser.Tokens[system].Position.Column)
			}

			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system+1])
			body, err := Parser.multiLineBodyApplication(system+2, lexer.OpenBraces, lexer.CloseBraces, make([]lexer.Token, 0))
			if err != nil {
				return nil, err
			}

			// Runs the parser on the body
			p := NewParser(tempLexer(body))
			if err := p.Start(); err != nil {
				return nil, err
			}

			system += len(body) + 2
			new.Body = p.GainedExpressions
			instruction.TokensAxis = append(instruction.TokensAxis, body...)
			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system])
			instruction.IF.Expression = append(instruction.IF.Expression, *new)
		case "else":
			if len(Parser.Tokens) > system+2 && Parser.Tokens[system+1].Sort != lexer.OpenBraces {
				return nil, fmt.Errorf("%d:%d if requires ( and not %s", Parser.Tokens[system].Position.Line, Parser.Tokens[system].Position.Column, Parser.Tokens[system+1].Literal)
			}

			var expression *Expression = new(Expression)
			expression.Keyword = Parser.Tokens[system]
			body, err := Parser.multiLineBodyApplication(system+2, lexer.OpenBody, lexer.CloseBraces, make([]lexer.Token, 0))
			if err != nil {
				return nil, err
			}

			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system])
			instruction.TokensAxis = append(instruction.TokensAxis, body...)

			// Runs the parser on the body
			p := NewParser(tempLexer(body))
			if err := p.Start(); err != nil {
				return nil, err
			}

			system += len(body) + 2
			expression.Body = p.GainedExpressions
			instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[system])
			instruction.IF.Expression = append(instruction.IF.Expression, *expression)
			return instruction, nil

		}

	}
	return instruction, nil
}

func (Parser *Parser) functionStatement() (*Instruction, error) {
	var instruction *Instruction = new(Instruction)
	instruction.FunctionInstruction = new(FunctionCall)
	instruction.NodeType = 4

	var err error = nil // Holds the error from the multiline function
	instruction.FunctionInstruction.Leadup = Parser.inlineBodyApplication(Parser.CurrentPosition, make([]lexer.Token, 0), lexer.OpenParenthesis)
	instruction.FunctionInstruction.Args, err = Parser.multiLineBodyApplication(Parser.CurrentPosition+len(instruction.FunctionInstruction.Leadup)+1, lexer.OpenParenthesis, lexer.CloseParenthesis, make([]lexer.Token, 0))
	if err != nil {
		return nil, err
	}

	instruction.TokensAxis = append(instruction.TokensAxis, instruction.FunctionInstruction.Leadup...)                                                                                     // Leadup
	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition+len(instruction.FunctionInstruction.Leadup)])                                             // Open arg
	instruction.TokensAxis = append(instruction.TokensAxis, instruction.FunctionInstruction.Args...)                                                                                       // Args for the system
	instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition+len(instruction.FunctionInstruction.Leadup)+len(instruction.FunctionInstruction.Args)+1]) // End part of the token
	return instruction, nil
}

func (Parser *Parser) expressionStatement() (*Instruction, error) {
	var instruction *Instruction = new(Instruction) // Holds the structure information
	instruction.NodeType = 3                        // Holds the node type
	instruction.TokensAxis = Parser.inlineBodyApplication(Parser.CurrentPosition, make([]lexer.Token, 0), lexer.CloseBody)
	//if len(instruction.TokensAxis) > len(Parser.Tokens) { // Err handles the statement
	//	return nil, fmt.Errorf("%d:%d expression presented with no variable", instruction.TokensAxis[0].Position.Line, instruction.TokensAxis[0].Position.Column)
	//}

	return instruction, nil
}

func (Parser *Parser) returnStatement() (*Instruction, error) {
	var instruction *Instruction = new(Instruction)
	instruction.TokensAxis = Parser.inlineBodyApplication(Parser.CurrentPosition, make([]lexer.Token, 0), lexer.SemiColon)
	instruction.NodeType = 2
	return instruction, nil
}

func (Parser *Parser) declarationStatement() (*Instruction, error) {
	var instruction *Instruction = new(Instruction)            // Makes Instruction structure
	instruction.DeclarationInstruction = new(DeclareStatement) // Makes declaration statement structure

	// Checks if the declaration is a constant declaration
	if Parser.Tokens[Parser.CurrentPosition].Literal == CONST {
		instruction.DeclarationInstruction.Constant = true
	}

	// Checks for a multi-line body declaration statement
	if Parser.Tokens[Parser.CurrentPosition+1].Sort == lexer.OpenParenthesis {
		instruction.DeclarationInstruction.BodyTYPED = true

		// Adds the first information from the tokenized information
		instruction.TokensAxis = append(instruction.TokensAxis, []lexer.Token{Parser.Tokens[Parser.CurrentPosition], Parser.Tokens[Parser.CurrentPosition+1]}...)

		// Parses the multiline body application
		TokensAxis, err := Parser.multiLineBodyApplication(Parser.CurrentPosition+2, lexer.OpenParenthesis, lexer.CloseParenthesis, make([]lexer.Token, 0))
		if err != nil {
			return nil, err
		}

		instruction.TokensAxis = append(instruction.TokensAxis, TokensAxis...)                                                     // Sets the tokenaxis
		instruction.TokensAxis = append(instruction.TokensAxis, Parser.Tokens[Parser.CurrentPosition+len(instruction.TokensAxis)]) // Adds the last peice of information
	} else { // Not a multi-line body application
		instruction.DeclarationInstruction.BodyTYPED = false                                                     // Not a multi-line body application
		instruction.TokensAxis = Parser.inlineBodyApplication(Parser.CurrentPosition, make([]lexer.Token, 0), 0) // Gets the inline body tokens for the application
	}

	instruction.NodeType = 1 // Sets the declaration type
	return instruction, nil  // Returns the instruction
}

// multiLineBodyApplication will read the bodys as a multi-line body
func (Parser *Parser) multiLineBodyApplication(position int, open, close lexer.TokenType, collection []lexer.Token) ([]lexer.Token, error) {

	var openBodys int = 0 // number of opened bodies inside of the main body
	for pos := position; pos < len(Parser.Tokens); pos++ {
		// Detects the new open in the body
		if Parser.Tokens[pos].Sort == open {
			openBodys++ // Open body
		}

		// Detects the new close in the body
		if Parser.Tokens[pos].Sort == close {
			if openBodys <= 0 {
				return collection, nil
			}

			openBodys-- // Close body
		}

		// Appends into the array
		collection = append(collection, Parser.Tokens[pos])
	}

	// Returns the error
	return make([]lexer.Token, 0), errors.New("body opened but not closed: " + strconv.Itoa(Parser.Tokens[position].Position.Line+1) + ":" + strconv.Itoa(Parser.Tokens[position].Position.Column))
}

// inlineBodyApplication will read the InLine tokens until a breakpoint is found
func (Parser *Parser) inlineBodyApplication(position int, collection []lexer.Token, close lexer.TokenType) []lexer.Token {

	// Runs the system information without errors happening
	for pos := position; pos < len(Parser.Tokens); pos++ {

		// Checks for the breakpoint on the line
		if Parser.Tokens[pos].Sort == lexer.SemiColon || Parser.Tokens[position].Position.Line < Parser.Tokens[pos].Position.Line || close > 0 && Parser.Tokens[pos].Sort == close {
			return collection
		}

		// Appends onto the collection
		collection = append(collection, Parser.Tokens[pos])
	}

	return collection
}

func tempLexer(body []lexer.Token) *lexer.Lexer {
	return &lexer.Lexer{
		Token: body,
	}
}
