package evaluator

import (
	"EvoScript/formatter"

	"EvoScript/engine"
	"EvoScript/lexer"
	"EvoScript/parse"
	"fmt"
	"io"
	"strings"
)

type Evaluator struct {
	Lexer           *lexer.Lexer      // Past tense lexer
	Parser          *parse.Parser     // Past tenser parser
	Formatter       *formatter.Format // Past tense formatter
	allocatedMemory []Pointer         // Stores all allocated memory
	writer          io.Writer         // Stores the writers input for the interface
	packages        map[string]any
}

// Makes the new evaluator
func NewEvaluator(wr io.Writer, Format *formatter.Format, memory ...[]Pointer) *Evaluator {
	var allocated []Pointer = make([]Pointer, 0)
	if len(memory) >= 1 {
		for _, memo := range memory {
			allocated = append(allocated, memo...)
		}
	}
	return &Evaluator{
		Lexer:           Format.Lexer,  // Lexer
		Parser:          Format.Parser, // Parser
		Formatter:       Format,        // Formatter
		allocatedMemory: allocated,     // Allocate a non specification limit of memory
		packages:        make(map[string]any),
		writer:          wr,
	}
}

// ExecuteString will execute EvoScript and log into a string
func (E *Evaluator) ExecuteString(tags ...string) (string, error) {
	var capture *strings.Builder = new(strings.Builder)

	if _, err := E.Execute(capture, tags...); err != nil {
		return "", err
	}

	return capture.String(), nil
}

func (E *Evaluator) Writer() io.Writer {
	return E.writer
}

// Run will execute the evaluator on the formatted output
func (E *Evaluator) Execute(wr io.Writer, tags ...string) ([]lexer.Token, error) {

	// Registers the include feature as default
	if err := E.Go2Evo("include", E.Include); err != nil {
		return nil, err
	}

	// Registers the echo feature inside EvoScript
	if err := E.Go2Evo("echo", E.echo); err != nil {
		return nil, err
	}

	// By default STD is imported on render
	if err := E.AccessSTD(); err != nil {
		return nil, err
	}

	// Ranges through the system
	for point := range E.Formatter.Manuals {
		indexed := E.Formatter.Manuals[point]

		//fmt.Println(indexed.Node, indexed.Axis)
		switch indexed.Node { // Controls what executes

		case 8: // Func creation
			function, err := E.allocatedFunction(indexed.FunctionBody)
			if err != nil {
				return nil, err
			}

			E.allocatedMemory = append(E.allocatedMemory, *function)
		case 7: // IF
			if err := E.if_exec(indexed.IFCall); err != nil {
				return nil, err
			}

		case 6: // Text
			raw, err := engine.New(indexed.Axis[0].Literal, tags...).ExecuteWithContextTags(E.EvoScriptEngineDriverEngineDriver)
			if err != nil {
				return nil, err
			}

			var message = raw
			if _, err := wr.Write([]byte(message)); err != nil {
				return nil, err
			}

		case 2: // Return statement is active here
			return E.ReturnStatement(indexed.Returns)

		case 3, 5: // Var/Func executed
			vars, err := E.lookupMemory(indexed.Axis) // Lookups the memory address
			if err != nil || len(vars) <= 0 {         // err handles and checks the length
				return nil, fmt.Errorf("%d:%d memory address was not found", indexed.Axis[0].Position.Line, indexed.Axis[0].Position.Column)
			}

			// Ranges through the objects provided
			for _, pointer := range vars {
				switch pointer.memory.(type) {

				case Var: // Variable detected
					if _, err := wr.Write([]byte(pointer.memory.(Var).value.Literal)); err != nil {
						return nil, err
					}

				case Object:
					if _, err := wr.Write([]byte(fmt.Sprint(pointer.memory.(Object).Value))); err != nil {
						return nil, err
					}
				}
			}

		case 1: // DeclareStatement
			if err := E.declare_var(indexed); err != nil {
				return nil, err
			}

		case 4: // FunctionCall statement
			if values, err := E.functionCall(indexed.CallStatement); err != nil {
				return nil, err
			} else {
				// Ranges through the objects provided
				for _, pointer := range values {
					switch pointer.memory.(type) {

					case Var: // Variable detected
						if _, err := wr.Write([]byte(pointer.memory.(Var).value.Literal)); err != nil {
							return nil, err
						}

					case Object:
						if _, err := wr.Write([]byte(fmt.Sprint(pointer.memory.(Object).Value))); err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	return nil, nil
}
