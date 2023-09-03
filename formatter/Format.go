package formatter

import (
	"EvoScript/lexer"
	"EvoScript/parse"
)

type Format struct {
	Parser *parse.Parser // Holds mainly the instructions
	Lexer  *lexer.Lexer  // Holds mainly the tokens

	Manuals  []FormattedPacket // Holds our entire formatter output
	Position int               // Scanners current position inside the array
}

// Makes the new formatter model
func NewFormat(parser *parse.Parser) *Format {
	return &Format{
		Parser:  parser,                     // Parser
		Lexer:   parser.Lexer,               // Lexer from parser
		Manuals: make([]FormattedPacket, 0), // New array formed for storage
	}
}

// Format will run the formatter on the information
func (f *Format) Format() error {

	// Ranges through the parser instructions
	for instruction := range f.Parser.GainedExpressions {
		inst := f.Parser.GainedExpressions[instruction] // Indexes the item
		f.Position = instruction                        // updates the position

		switch inst.NodeType { // Control seqs through

		case 8:
			function, err := f.FuncCreate(inst.FunctionCreation)
			if err != nil {
				return err
			}

			f.Manuals = append(f.Manuals, *function)
		case 7: // IF
			manual, err := f.IF(inst.IF)
			if err != nil {
				return err
			}

			// Appends into the manual list
			f.Manuals = append(f.Manuals, *manual)
		case 6:
			f.Manuals = append(f.Manuals, FormattedPacket{Node: 6, Axis: inst.TokensAxis})
		case 3:
			f.Manuals = append(f.Manuals, FormattedPacket{Node: 3, Axis: inst.TokensAxis})
		case 2:
			f.Manuals = append(f.Manuals, FormattedPacket{Node: 2, Axis: inst.TokensAxis, Returns: inst.TokensAxis[1:]})
		case 5:
			f.Manuals = append(f.Manuals, FormattedPacket{Node: 5, Axis: inst.TokensAxis})
		case 1: // Declaration statement
			formatted, err := f.var_format()
			if err != nil {
				return err
			}

			// Appends with the past structures
			f.Manuals = append(f.Manuals, *formatted)

		case 4: // Appends into the manuals structure
			f.Manuals = append(f.Manuals, FormattedPacket{Axis: inst.TokensAxis, Node: 4, CallStatement: &FunctionCall{Path: inst.FunctionInstruction.Leadup, Args: inst.FunctionInstruction.Args}})
		}
	}

	return nil
}

// FuncCreate will return the function inside a formatted packet
func (F *Format) FuncCreate(body *parse.FunctionBody) (*FormattedPacket, error) {
	var format *FormattedPacket = new(FormattedPacket)
	format.Node = 8
	format.FunctionBody = new(FunctionBody)

	format.FunctionBody.Keyword = body.Keyword
	format.FunctionBody.ArgsWants = body.ArgsWants
	format.FunctionBody.ReturnArgs = body.ReturnArgs

	formats := NewFormat(&parse.Parser{GainedExpressions: body.Bodys})
	if err := formats.Format(); err != nil {
		return nil, err
	}

	format.FunctionBody.Bodys = formats.Manuals
	return format, nil
}

// IF will work the needed information within the if statement
func (F *Format) IF(body *parse.IfStatement) (*FormattedPacket, error) {
	var format *FormattedPacket = new(FormattedPacket)
	format.Node = 7
	format.IFCall = make([]ExpressionIF, 0)

	// Ranges through the if statement options given
	for _, exp := range body.Expression {
		var newIF *ExpressionIF = new(ExpressionIF)
		formats := NewFormat(&parse.Parser{GainedExpressions: exp.Body})
		if err := formats.Format(); err != nil {
			return nil, err
		}

		newIF.Args = exp.Args
		newIF.Keyword = exp.Keyword
		newIF.Body = formats.Manuals
		format.IFCall = append(format.IFCall, *newIF)
	}

	return format, nil
}
