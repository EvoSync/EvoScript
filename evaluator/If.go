package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
	"fmt"
	"strconv"
)

// if_exec will execute & handle the if statements
func (E *Evaluator) if_exec(point []formatter.ExpressionIF) error {
	// Ranges through the objects given
	for _, object := range point {

		// Checks if args are needed inside the statement
		if object.Keyword.Literal == "if" || object.Keyword.Literal == "elif" {
			point, err := E.objectiveLookup(object.Args)
			if err != nil {
				return err
			}

			// false
			if !point {
				continue
			}

			// Executes the body which has been passed
			eval := NewEvaluator(E.writer, &formatter.Format{Manuals: object.Body}, E.allocatedMemory)
			if _, err := eval.Execute(E.writer); err != nil {
				return err
			}
			break
		}

		// Executes the body which has been passed
		eval := NewEvaluator(E.writer, &formatter.Format{Manuals: object.Body}, E.allocatedMemory)
		if _, err := eval.Execute(E.writer); err != nil {
			return err
		}

		break
	}

	return nil
}

var IF_Sections []lexer.TokenType = []lexer.TokenType{
	lexer.EqualEqual, lexer.NotEqual,
	lexer.GreaterEqual, lexer.GreaterThan,
	lexer.LessEqual, lexer.LessThan,
	lexer.AndAnd, lexer.OrOr,
}

// objectiveLookup will take args sort into sections then run lookup
func (E *Evaluator) objectiveLookup(args []lexer.Token) (bool, error) {
	var sections [][]lexer.Token = make([][]lexer.Token, 1) // Ranges through sections
	var opened int = 0                                      // Stores all opened
	var operators []lexer.Token = make([]lexer.Token, 0)

	// Ranges through the args
	for _, object := range args {
		if object.Sort == lexer.OpenParenthesis {
			opened++
		} else if object.Sort == lexer.CloseParenthesis {
			opened--
		}

		if opened <= 0 && unwrapTokenArray(IF_Sections, object.Sort) {
			operators = append(operators, object)               // Saves into the map
			sections = append(sections, make([]lexer.Token, 0)) // Makes a new element space
			continue
		}
		sections[len(sections)-1] = append(sections[len(sections)-1], object)
	}

	var (
		// All token stages inside the objects
		singleStage []lexer.Token   = make([]lexer.Token, 0)
		locked      lexer.TokenType = 0
	)
	for pos, stage := range sections {
		single, err := E.args(stage, 0)
		if err != nil {
			return false, err
		}

		if pos == 0 {
			locked = single.memory.(Var).value.Sort
		} else if single.memory.(Var).value.Sort != locked {
			return false, fmt.Errorf("%d:%d all types must be the same", single.memory.(Var).value.Position.Line, single.memory.(Var).value.Position.Column)
		}

		singleStage = append(singleStage, *single.memory.(Var).value)
	}

	// Checks for a boolean inside the object
	if len(singleStage) <= 1 && singleStage[0].Sort == lexer.BOOL {
		formatted, err := strconv.ParseBool(singleStage[0].Literal)
		if err != nil {
			return false, err
		}

		return formatted, nil
	}

	var (
		targ  lexer.Token = singleStage[0]
		found bool        = false
	)

	for val := 1; val < len(singleStage); val++ {

		switch operators[val-1].Sort {
		case lexer.EqualEqual: // ==
			found = targ.Literal == singleStage[val].Literal
		case lexer.NotEqual: // !=
			found = targ.Literal != singleStage[val].Literal
		case lexer.GreaterThan:
			one, err := strconv.Atoi(singleStage[val].Literal)
			if err != nil {
				return false, err
			}

			two, err := strconv.Atoi(targ.Literal)
			if err != nil {
				return false, err
			}

			found = two > one
		case lexer.GreaterEqual:
			one, err := strconv.Atoi(singleStage[val].Literal)
			if err != nil {
				return false, err
			}

			two, err := strconv.Atoi(targ.Literal)
			if err != nil {
				return false, err
			}

			found = two >= one
		case lexer.LessThan:
			one, err := strconv.Atoi(singleStage[val].Literal)
			if err != nil {
				return false, err
			}

			two, err := strconv.Atoi(targ.Literal)
			if err != nil {
				return false, err
			}

			found = two < one
		case lexer.LessEqual:
			one, err := strconv.Atoi(singleStage[val].Literal)
			if err != nil {
				return false, err
			}

			two, err := strconv.Atoi(targ.Literal)
			if err != nil {
				return false, err
			}

			found = two <= one
		}
	}

	return found, nil
}
