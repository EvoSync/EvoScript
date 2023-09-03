package engine

import (
	"io"
	"strings"
)

// Holds the EvoScript engine version.
const EngineVersion string = "IN_DEV - v1.0"

type NewEngine struct {
	Source string 				// stores the source for the engine to run over
	Tags   []string				// Holds the tags for the engine
	Pos	   int					// Position of the executer
}

// New creates the structure
func New(source string, tags ...string) *NewEngine {

	if len(tags) != 2 {
		tags = make([]string, 0)
		tags = append(tags, "<<")	// OpenTag
		tags = append(tags, ">>")	// CloseTag
	}

	return &NewEngine{
		Source: source, 			// Holds the source for the engine
		Tags: 	tags,				// Sets the tags for the engine
		Pos: 	0,
	}
}

// ExecuteWithContextTags will execute the given source under the source passed
func (engine *NewEngine) ExecuteWithContextTags(exe func(string, io.Writer) (int, error)) (string, error) {
	var writer strings.Builder

	var charaters []string = strings.Split(engine.Source, "") // splits charater by charater
	for pos := 0; pos < len(charaters); pos++ {				  // Loops through the charaters
		engine.Pos = pos

		
		// Checks if the line is a tag
		if charaters[pos] == string(engine.Tags[0][0]) && engine.checkInlineTag(0, charaters[pos:]) {
			var Ctag    string = "" // Holds the tag insider which was called
			var tag  int = 0		// The emulators position on the string
			
			
			// Loops through the tag until it closes
			for tag = pos + len(engine.Tags[0]); tag < len(charaters); tag++ {
				Ctag += charaters[tag]

				// Checks the charater if its a valid tag and begins closing down the forloop
				if charaters[tag] == string(engine.Tags[1][0]) && engine.checkInlineTag(1, charaters[tag:]) {
					Ctag += engine.Tags[1][1:]
					break
				}
			}

			// Detects if the tag was valid or not
			if strings.Join(strings.Split(Ctag, "")[len(Ctag)-len(engine.Tags[1]):], "") != engine.Tags[1] {
				continue // continues to loop through
			}

			// Removes the tag ending
			Ctag = strings.TrimSuffix(Ctag, engine.Tags[1])
			pos  = tag+len(engine.Tags[1])-1

			// Executes the tag route
			if _, err := exe(Ctag, &writer); err != nil {
				return "", err
			}


			continue
		}

		writer.WriteString(charaters[pos])
	}


	return writer.String(), nil
}