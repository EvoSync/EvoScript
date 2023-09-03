package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
	"EvoScript/parse"
	"EvoScript/std"
	"bytes"
	"io/ioutil"
	"strings"
)

// NewPackage will register the package inside packages map
func (e *Evaluator) NewPackage(name string, packages any) {
	e.packages[name] = packages
}

// Include will try to include packages and register them inside memory
func (e *Evaluator) Include(path []string) {
	for _, paths := range path {

		// Checks for a filepath import
		if strings.Contains(paths, ".") {
			Contents, err := ioutil.ReadFile(paths)
			if err != nil {
				panic(err.Error())
			}

			lexer := lexer.NewLexer(string(Contents), "\n", true)
			if err := lexer.Start(); err != nil {
				continue
			}

			parse := parse.NewParser(lexer)
			if err := parse.Start(); err != nil {
				continue
			}

			format := formatter.NewFormat(parse)
			if err := format.Format(); err != nil {
				continue
			}

			eval := NewEvaluator(e.writer, format, e.allocatedMemory)
			if _, err := eval.Execute(e.writer); err != nil {
				continue
			}

			e.allocatedMemory = append(e.allocatedMemory, eval.allocatedMemory...)

			continue
		}

		// Has the packages all registered
		object, ok := e.packages[paths]
		if !ok {
			return
		}

		// Creates the package object
		e.Go2Evo(paths, object)
	}
}

// AccessSTD will access all STD objects
func (e *Evaluator) AccessSTD() error {
	for name, value := range std.STD() {
		if err := e.Go2Evo(name, value); err != nil {
			return err
		}
	}

	return nil
}

// echo is a builtin function for EvoScript
func (eval *Evaluator) echo(args []string) {
	buf := bytes.NewBuffer(nil)
	for p := range args[:len(args)-1] {
		buf.Write([]byte(args[p] + " "))
	}

	buf.Write([]byte(args[len(args)-1] + "\r\n"))

	eval.writer.Write(buf.Bytes())
}
