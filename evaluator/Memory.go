package evaluator

import (
	"EvoScript/formatter"
	"EvoScript/lexer"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Pointer struct {
	address []byte // Stores the memory address
	memory  any    // Stores the value of the memory allocated
}

type Var struct {
	module string       // Name of the memory module
	value  *lexer.Token // Value of the memory to be allocated
}

type Object struct {
	module   string    // Name of the memory object to be allocated
	subheads []Pointer // Stores all the names of the subheaders inside the object
	Value    any
}

type GoFunction struct {
	module string            // Name of the module function to be allocated
	wanted []lexer.TokenType // Stores all the args needed for the function seal
	xseal  any               // Stores the function model
	rets   []lexer.TokenType // Stores all objects it returns
}

type CreatedFunction struct {
	Keyword    lexer.Token                 // Stores the keyword
	ArgsWants  []lexer.Token               // Stores all the arguments needed
	ReturnArgs []lexer.Token               // Stores all the return arguments
	Bodys      []formatter.FormattedPacket // Stores the body of the function
}

// allocatedVar will allocate var the memory inside the memory relays
func (eval *Evaluator) allocatedVar(header string, value *lexer.Token) (*Pointer, error) {
	call := "$" + header       // Automatically adds the var callsign "$"
	address := make([]byte, 7) // Randomly generates a memory address
	rand.Read(address)         // Randomly updates the memory address
	return &Pointer{address: address, memory: Var{module: call, value: value}}, nil
}

// allocatedFunction will allocate var the memory inside the memory relays
func (eval *Evaluator) allocatedFunction(function *formatter.FunctionBody) (*Pointer, error) {
	address := make([]byte, 7) // Randomly generates a memory address
	rand.Read(address)         // Randomly updates the memory address
	return &Pointer{address: address, memory: CreatedFunction{Keyword: function.Keyword, ArgsWants: function.ArgsWants, ReturnArgs: function.ReturnArgs, Bodys: function.Bodys}}, nil
}

// lookupMemory looks up the memory address inside the index
func (eval *Evaluator) lookupMemory(path []lexer.Token, args ...Pointer) ([]Pointer, error) {
	var segmants [][]lexer.Token = make([][]lexer.Token, 1) // Makes an array of array tokens

	// loops through all the tokens in path
	for position := 0; position < len(path); position++ {
		token := path[position] // Grabs token expression

		// Ignores fullstop tokens
		if token.Sort == lexer.Fullstop {
			segmants = append(segmants, make([]lexer.Token, 0))
			continue
		}

		// Appends into the current memory segmant
		segmants[len(segmants)-1] = append(segmants[len(segmants)-1], token)
	}

	var base []Pointer = eval.allocatedMemory // Stores the base allocated memory
	for pos, weight := range segmants {       // Ranges through all the segmants inside the path
		if weight[0].Sort != lexer.Dollar { // Checks for vars inside the weight
			var found bool = false         // Stores if we found the memory base
			for _, setting := range base { // Ranges through all the possible bases

				// Only accepts object structures past this point
				if strings.ToLower(strings.Replace(filepath.Ext(reflect.TypeOf(setting.memory).String()), ".", "", 1)) != "object" {
					continue
				}

				// Only accepts the validated object structure passed here
				if setting.memory.(Object).module != weight[0].Literal {
					continue
				}

				found = true                            // Confirms we found the object
				base = setting.memory.(Object).subheads // Updates the pointers base into the objects base

				if pos+1 == len(segmants) {
					return []Pointer{setting}, nil
				}
			}

			if found { // Object found
				continue // Continue to loop
			}
		}

		// Detects a variable lookup inside the object
		if weight[0].Sort == lexer.Dollar {

			// Indexes the object inside memory
			for position := 0; position < len(base); position++ {
				if strings.ToLower(strings.Replace(filepath.Ext(reflect.TypeOf(base[position].memory).String()), ".", "", 1)) != "var" {
					continue
				}

				// Found the memory object within the statement
				if base[position].memory.(Var).module == "$"+weight[1].Literal {
					return []Pointer{base[position]}, nil
				}
			}

			return nil, errors.New(strconv.Itoa(path[0].Position.Line) + ":" + strconv.Itoa(path[0].Position.Column) + ": var `" + weight[len(weight)-1].Literal + "` is not an exported value inside memory")
		}

		// Function execution validation
		if weight[0].Sort == lexer.INDENT {
			for position := 0; position < len(base); position++ {
				switch strings.ReplaceAll(filepath.Ext(reflect.TypeOf(base[position].memory).String()), ".", "") {
				case "GoFunction": // Function execution
					var function GoFunction = base[position].memory.(GoFunction)
					if weight[0].Literal != function.module {
						continue
					}

					if len(function.wanted) != len(args) && !HasVariadic(function.wanted) {
						return nil, fmt.Errorf("%d:%d function wants %d args but has %d args", weight[0].Position.Line, weight[0].Position.Column, len(function.wanted), len(args))
					}

					var values []reflect.Value = make([]reflect.Value, 0)
					for pos := 0; pos < len(args); pos++ {
						arg := args[pos]
						switch arg.memory.(type) {

						case Object:
							values = append(values, reflect.ValueOf(arg.memory.(Object).Value))
						case Var:
							object := arg.memory.(Var).value

							if object.Sort != function.wanted[pos] && !HasVariadic(function.wanted) && function.wanted[pos] != lexer.Any {
								return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", weight[0].Position.Line, weight[0].Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(function.wanted), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
							}

							if function.wanted[pos] >= 90 {
								if pos != len(function.wanted)-1 {
									return nil, fmt.Errorf("%d:%d variadic args must only appear at the end of args", weight[0].Position.Line, weight[0].Position.Column)
								}

								var index []reflect.Value = make([]reflect.Value, 0)
								for _, argument := range args[pos:] {
									switch argument.memory.(type) {
									case Object:
										index = append(index, reflect.ValueOf(argument.memory.(Object).module))
										continue
									case Var:
										value, err := ToValue(argument.memory.(Var).value, function.wanted, args)
										if err != nil {
											return nil, err
										}

										if function.wanted[pos] == lexer.VariadicString && value.Kind() != reflect.String {
											return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", weight[0].Position.Line, weight[0].Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(function.wanted), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
										} else if function.wanted[pos] == lexer.VariadicInt && value.Kind() != reflect.Int {
											return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", weight[0].Position.Line, weight[0].Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(function.wanted), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
										} else if function.wanted[pos] == lexer.VariadicBool && value.Kind() != reflect.Bool {
											return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", weight[0].Position.Line, weight[0].Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(function.wanted), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
										}

										index = append(index, *value)
										continue
									}
								}

								if function.wanted[pos] == lexer.VariadicString {
									var str []string = make([]string, 0)
									for _, v := range index {
										str = append(str, v.String())
									}
									values = append(values, reflect.ValueOf(str))
								} else if function.wanted[pos] == lexer.VariadicInt {
									var ints []int = make([]int, 0)
									for _, i := range index {
										conversion, _ := strconv.Atoi(fmt.Sprint(i.Interface()))
										ints = append(ints, conversion)
									}
									values = append(values, reflect.ValueOf(ints))
								} else if function.wanted[pos] == lexer.VariadicBool {
									var bools []bool = make([]bool, 0)
									for _, b := range index {
										conversion, _ := strconv.ParseBool(fmt.Sprint(b.Interface()))
										bools = append(bools, conversion)
									}
									values = append(values, reflect.ValueOf(bools))
								} else if function.wanted[pos] == lexer.VariadicAny {
									var indexed []any = make([]any, 0)
									for _, a := range index {
										indexed = append(indexed, a.Interface())
									}
									values = append(values, reflect.ValueOf(indexed))
								}

								pos += len(index)
								break
							}

							value, err := ToValue(object, function.wanted, args)
							if err != nil {
								return nil, err
							}

							values = append(values, *value)
						}
					}

					values_returned := reflect.ValueOf(function.xseal).Call(values)

					var objects []Pointer = make([]Pointer, 0)
					for _, val := range values_returned {
						pointer, err := eval.go2Evo("", val.Interface())
						if err != nil {
							return nil, err
						}

						objects = append(objects, *pointer)
					}

					return objects, nil

				case "CreatedFunction": // Function execution

				}
			}

		}
	}
	return nil, fmt.Errorf("%d:%d memory passed was not found: [%s]", path[0].Position.Line, path[0].Position.Line, strings.Join(convertTokensToTokenLiteralArray(path), ""))
}

func ToValue(token *lexer.Token, wants []lexer.TokenType, args []Pointer) (*reflect.Value, error) {
	switch token.Sort {

	case lexer.STRING:
		token.Literal = lexer.AnsiUtil(token.Literal)
		val := reflect.ValueOf(token.Literal)
		return &val, nil
	case lexer.INT:
		rawINT, err := strconv.Atoi(token.Literal)
		if err != nil {
			return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", token.Position.Line, token.Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(wants), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
		}

		val := reflect.ValueOf(rawINT)
		return &val, nil
	case lexer.BOOL:
		rawBOOL, err := strconv.ParseBool(token.Literal)
		if err != nil {
			return nil, fmt.Errorf("%d:%d args dont match what are wanted. Wants: func(%s) Has: func(%s)", token.Position.Line, token.Position.Column, strings.Join(convertTokenTypesToTokenTypeArray(wants), ", "), strings.Join(convertTokensToTokenTypeArray(args), ", "))
		}

		val := reflect.ValueOf(rawBOOL)
		return &val, nil
	case lexer.Any:
		fmt.Println("HERE")
	default:
		return nil, nil
	}

	return nil, nil
}

func HasVariadic(args []lexer.TokenType) bool {
	for _, arg := range args {
		if arg >= 90 {
			return true
		}
	}

	return false
}
