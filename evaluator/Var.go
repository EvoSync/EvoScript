package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
	"sort"

	"golang.org/x/exp/maps"
)

// declare_var will register the variable inside the memory
func (E *Evaluator) declare_var(indexed formatter.FormattedPacket) error {

	// Ranges through all the methods within the system
	for index := range indexed.DeclareStatement.Methods {
		indexed := indexed.DeclareStatement.Methods[index]

		// Function execution route
		if indexed.WholeSomeFunction {
			return E.MultiValues(indexed)
		}

		// Ranges through the memorys allocated
		for _, memory := range indexed.Models {

			// Works the args for the memory function
			v, err := E.args(memory.Values, memory.Locked)
			if err != nil {
				return err
			}

			var Pure *Pointer = new(Pointer)
			switch v.memory.(type) {
			case Var:
				mem, err := E.allocatedVar(memory.Model.Literal, v.memory.(Var).value)
				if err != nil {
					return err
				}

				Pure = mem
			case Object:
				converted, err := E.go2Evo(memory.Model.Literal, v.memory.(Object).Value)
				if err != nil {
					return err
				}

				Pure = converted
			}

			//// Saves into memory
			E.allocatedMemory = append(E.allocatedMemory, *Pure)
		}

	}
	return nil
}

// Allows function to return multiply values and hold in memory
func (E *Evaluator) MultiValues(values formatter.DeclareLine) error {
	var functionPath []lexer.Token = make([]lexer.Token, 0)
	var assign []lexer.Token = make([]lexer.Token, 0)

	keys := maps.Keys(values.Models)
	sort.Ints(keys)

	for _, value := range keys { // Ranges through the values
		module := values.Models[value]
		if value == 0 { // First position holds function
			functionPath = module.Values // Sets the function path to hold that information
		}

		if module.Locked != 0 {
			module.Model.Sort = module.Locked
		}

		assign = append(assign, module.Model)
	}

	var (
		paths    []lexer.Token   = make([]lexer.Token, 0)   // path to the function
		tokens   [][]lexer.Token = make([][]lexer.Token, 1) // tokens used as args
		switched bool            = false
	)

	// Ranges through acting as a multireader
	for _, path := range functionPath {

		// Checks for a body
		if path.Sort == lexer.OpenParenthesis {
			switched = !switched // flips the boolean
			continue             // continues to loop
		} else if path.Sort == lexer.CloseParenthesis {
			break // breaks from loop
		}

		if !switched { // saves as path
			paths = append(paths, path)
		} else { // Checks for an operator and makes a new allig
			if unwrapTokenArray([]lexer.TokenType{lexer.Comma}, path.Sort) {
				tokens = append(tokens, make([]lexer.Token, 0))
				continue
			}

			// saves into array
			tokens[len(tokens)-1] = append(tokens[len(tokens)-1], path)
		}
	}

	var pure []Pointer = make([]Pointer, 0) // Temp array
	for _, arg := range tokens {            // Ranges through tokens

		// Only accepts length larger than 0
		if len(arg) <= 0 {
			continue
		}

		// args worker for values
		current, err := E.args(arg, 0)
		if err != nil {
			return err
		}

		pure = append(pure, *current)
	}

	// Lookup memory will lookup the function and execute it
	vals, err := E.lookupMemory(paths, pure...)
	if err != nil {
		return err
	}

	// Ranges through the values
	for pos, sysVal := range assign {
		valueFromFunction := vals[pos]

		switch valueFromFunction.memory.(type) {

		case Var:
			mem, err := E.allocatedVar(sysVal.Literal, valueFromFunction.memory.(Var).value)
			if err != nil {
				return err
			}

			E.allocatedMemory = append(E.allocatedMemory, *mem)
		case Object:
			converted, err := E.go2Evo(sysVal.Literal, valueFromFunction.memory.(Object).Value)
			if err != nil {
				return err
			}

			//fmt.Println(converted.memory)
			E.allocatedMemory = append(E.allocatedMemory, *converted)
		}

	}

	return nil
}
