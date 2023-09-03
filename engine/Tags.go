package engine

// checkInlineTag will checks if the tag in the line is valid or not
func (engine *NewEngine) checkInlineTag(number int, current []string) bool {
	var position int = 0 // Holds the position for inside the for loop

	// Loops throughout the text ensuring that its found
	for position = 0; position < len(current) && position < len(engine.Tags[number]); position++ {

		// Compares the values inside the system
		if current[position] != string(engine.Tags[number][position]) {
			return false
		}
	}

	return position == len(engine.Tags[number])
}
