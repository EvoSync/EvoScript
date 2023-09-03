package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
	"EvoScript/parse"
	"io"
)

// EvoScriptEngineDriver will execute the tags inside this engine driver
func (eval *Evaluator) EvoScriptEngineDriverEngineDriver(tag string, wr io.Writer) (int, error) {

	lexer := lexer.NewLexer(tag, "\n", false)
	if err := lexer.Start(); err != nil {
		return 0, err
	}

	parse := parse.NewParser(lexer)
	if err := parse.Start(); err != nil {
		return 0, err
	}

	format := formatter.NewFormat(parse)
	if err := format.Format(); err != nil {
		return 0, err
	}

	e := NewEvaluator(eval.writer, format, eval.allocatedMemory)
	_, err := e.Execute(wr)
	eval.allocatedMemory = append(eval.allocatedMemory, e.allocatedMemory...)
	return 0, err
}
