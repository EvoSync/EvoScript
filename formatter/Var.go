package formatter

import (
	"EvoScript/lexer"
	"errors"
	"sort"
	"strconv"

	"golang.org/x/exp/maps"
)

// var_format will run the formatter on the current param
func (F *Format) var_format() (*FormattedPacket, error) {
	indexed := F.Parser.GainedExpressions[F.Position]

	if indexed.DeclarationInstruction == nil { // Checks for a nil pointer statement
		return nil, errors.New("unformattable declaration statement passed at: " + strconv.Itoa(indexed.TokensAxis[0].Position.Line+1) + ":" + strconv.Itoa(indexed.TokensAxis[0].Position.Column))
	}

	// Non multi-line declarations statement
	if !indexed.DeclarationInstruction.BodyTYPED {
		line, err := F.var_inline_statement(indexed.TokensAxis[1:])
		if err != nil {
			return nil, err
		}

		// Returns the formatted packets
		return &FormattedPacket{Node: 1, DeclareStatement: &Declare{Methods: []DeclareLine{*line}}, Axis: indexed.TokensAxis}, nil
	}

	var lines []DeclareLine = make([]DeclareLine, 0)
	// Ranges through each sample line within the settings
	for _, lineTokens := range F.sampleSplitByLine(indexed.TokensAxis[1:]) {
		if len(lineTokens) == 0 {
			continue
		}

		// Checks for the parathesis on the line ensures its done
		if lineTokens[0].Sort == lexer.OpenParenthesis || lineTokens[0].Sort == lexer.CloseParenthesis {
			lineTokens = lineTokens[1:] // Removes the first charater
		}

		if len(lineTokens) <= 0 { // Checks the length
			continue // Blank line
		}

		if len(lineTokens) < 3 { // Checks the length on the system passed
			return nil, errors.New("flagged declaration line [possible syntax]: " + strconv.Itoa(lineTokens[0].Position.Line+1) + ":" + strconv.Itoa(lineTokens[0].Position.Column))
		}

		// var_inline_statement will work the information
		current, err := F.var_inline_statement(lineTokens)
		if err != nil {
			return nil, err
		}

		// Appends onto the array
		lines = append(lines, *current)
	}

	// Returns the values within the system
	return &FormattedPacket{Node: 1, DeclareStatement: &Declare{Methods: lines}, Axis: indexed.TokensAxis}, nil
}

// var_inline_statement will parse the inline statement and return a model
func (F *Format) var_inline_statement(tokens []lexer.Token) (*DeclareLine, error) {
	line := new(DeclareLine)                  // Makes the new structure pointer
	line.Models = make(map[int]*declareModel) // Makes the map needed to store the information
	var rendered_exit_signal int = 0          // Stores our rendered_exit_signal leave which allows us to start later

	// Ranges through the tokens until the exit code
	for position := 0; position < len(tokens); position++ {

		// Break point has been detected
		if tokens[position].Sort == lexer.Equal {
			rendered_exit_signal = position
			break
		}

		// Comma has been detected and ignored
		if tokens[position].Sort == lexer.Comma || tokens[position].Sort == lexer.SemiColon {
			continue
		}

		line.Models[position] = new(declareModel)      // New model
		line.Models[position].Model = tokens[position] // Sets the new model

		// Possible locked down type given at this vendor
		if position+1 < len(tokens) && tokens[position+1].Sort == lexer.INDENT {
			switch tokens[position+1].Literal { // Allocates the type
			case "string", "str":
				line.Models[position].Locked = lexer.STRING // string
			case "number", "int":
				line.Models[position].Locked = lexer.INT // integer
			case "boolean", "bool":
				line.Models[position].Locked = lexer.BOOL // Boolean
			}

			position++ // additional movement
		}
	}

	key := maps.Keys(line.Models) // Access all the keys
	sort.Ints(key)                // Sorts by size order
	shifter := 0                  // Stores the current key
	var bodys int = 0

	for new := rendered_exit_signal + 1; new < len(tokens); new++ {

		if bodys == 0 && tokens[new].Sort == lexer.Comma {
			shifter++
			continue
		}

		if tokens[new].Sort == lexer.OpenParenthesis {
			bodys++
		} else if tokens[new].Sort == lexer.CloseParenthesis {
			bodys--
		}

		line.Models[key[shifter]].Values = append(line.Models[key[shifter]].Values, tokens[new])
	}

	if shifter != len(line.Models)-1 {
		if line.Models[0].Values[len(line.Models[0].Values)-1:][0].Sort == lexer.CloseParenthesis {
			line.WholeSomeFunction = true
		} else {
			return nil, errors.New("missing values associated with variables at: " + strconv.Itoa(line.Models[0].Model.Position.Line) + strconv.Itoa(line.Models[0].Model.Position.Column))
		}
	}

	return line, nil
}
