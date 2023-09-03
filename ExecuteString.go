package EvoScript

import (
	"EvoScript/evaluator"
	"EvoScript/formatter"
	"EvoScript/lexer"
	"EvoScript/parse"
	"io"
)

// ExecuteString is the mainly supported function for executing EvoScript within Applicaitons.
func ExecuteString(source string, wr io.Writer, elements map[string]any) error {
	lex := lexer.NewLexer(source, "\n", true)
	if err := lex.Start(); err != nil {
		return err
	}

	par := parse.NewParser(lex)
	if err := par.Start(); err != nil {
		return err
	}

	f := formatter.NewFormat(par)
	if err := f.Format(); err != nil {
		return err
	}

	evaluator := evaluator.NewEvaluator(wr, f)
	for name, value := range elements {
		evaluator.Go2Evo(name, value)
	}

	_, err := evaluator.Execute(wr)
	return err
}

