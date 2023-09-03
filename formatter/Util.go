package formatter

import (
	"EvoScript/lexer"
)

// sampleSplitByLine takes an array of tokens and sorts them by line
func (F *Format) sampleSplitByLine(tokens []lexer.Token) [][]lexer.Token {

	var (
		// Stores our future objects inside the system
		segmants  [][]lexer.Token = make([][]lexer.Token, 1)
		startLine int             = 1
	)

	// Loops through all the different tokens
	for scanner := 0; scanner < len(tokens); scanner++ {

		if tokens[scanner].Position.Line > startLine || tokens[scanner].Sort == lexer.SemiColon {
			startLine = tokens[scanner].Position.Line                                      // Updates the new line
			segmants = append(segmants, make([]lexer.Token, 0))                            // Makes the new element
			segmants[len(segmants)-1] = append(segmants[len(segmants)-1], tokens[scanner]) // Appends the current charater
			continue
		}

		//if tokens[scanner].Position.Line == startLine {										// Checks for new line
		//	startLine = tokens[scanner].Position.Line - 1										// Updates the new line
		//	segmants = append(segmants, make([]lexer.Token, 0))								// Makes the new element
		//	segmants[len(segmants)-1] = append(segmants[len(segmants)-1], tokens[scanner])	// Appends the current charater
		//	continue
		//}

		// Appends onto the array ensuring its done
		segmants[len(segmants)-1] = append(segmants[len(segmants)-1], tokens[scanner])
	}

	// Returns the segmants
	return segmants
}
