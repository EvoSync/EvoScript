package evaluator

import (
	"EvoScript/lexer"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"unicode"
)

func (eval *Evaluator) Go2Evo(name string, module any) error {
	point, err := eval.Go2Evo(name, module)
	if err != nil {
		return err
	}

	eval.allocatedMemory = append(eval.allocatedMemory, *point)
	return nil
}

// Go2Evo acts as an interface between golang and EvoScript values
func (eval *Evaluator) Go2Evo(name string, module any) (*Pointer, error) {
	switch module.(type) { // Controls what gets registered as what kind of module

	case string:
		return eval.allocatedVar(name, &lexer.Token{Sort: lexer.STRING, Literal: module.(string)})
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return eval.allocatedVar(name, &lexer.Token{Sort: lexer.INT, Literal: fmt.Sprint(module)})
	case float64, float32:
		return eval.allocatedVar(name, &lexer.Token{Sort: lexer.INT, Literal: fmt.Sprint(module)})
	case bool:
		return eval.allocatedVar(name, &lexer.Token{Sort: lexer.BOOL, Literal: fmt.Sprint(module)})

	default:
		switch reflect.TypeOf(module).Kind() {

		case reflect.Map:
			var new *Object = new(Object)
			new.subheads = make([]Pointer, 0)
			new.module = name
			new.Value = module

			iter := reflect.ValueOf(module).MapRange()
			for iter.Next() {
				key := iter.Key()
				val := iter.Value()

				converted, err := eval.Go2Evo(key.String(), val.Interface())
				if err != nil {
					return nil, err
				}

				new.subheads = append(new.subheads, *converted)
			}

			address := make([]byte, 7) // Randomly generates a memory address
			rand.Read(address)         // Randomly updates the memory address
			return &Pointer{address: address, memory: *new}, nil

		case reflect.Struct:
			var new *Object = new(Object)     // Makes the object
			new.subheads = make([]Pointer, 0) // Creates the array pointer
			new.module = name                 // Assigns the module name
			new.Value = module
			// Loops through all methods inside the object
			for pos := 0; pos < reflect.TypeOf(module).NumField(); pos++ {
				field := reflect.ValueOf(module).Field(pos) // Gets the field

				if !unicode.IsUpper(rune(reflect.TypeOf(module).Field(pos).Name[0])) {
					continue
				}

				// Runs the non_reg function for Go2Evo
				addr, err := eval.Go2Evo(strings.ToLower(reflect.TypeOf(module).Field(pos).Name), field.Interface())
				if err != nil {
					return nil, err
				}

				// Appends into the structure
				new.subheads = append(new.subheads, *addr)
			}

			address := make([]byte, 7) // Randomly generates a memory address
			rand.Read(address)         // Randomly updates the memory address
			return &Pointer{address: address, memory: *new}, nil

		case reflect.Func:
			// Converts all values into the pure go routines
			takes, err := toTypesArray(strings.Split(strings.Split(strings.Join(strings.Split(reflect.TypeOf(module).String(), "(")[1:], ""), ")")[0], ","))
			if err != nil {
				return nil, err
			}

			var returns []lexer.TokenType = make([]lexer.TokenType, 0)                                             // Stores all the returns value inside the token type
			if strings.Contains(strings.Join(strings.Split(reflect.ValueOf(module).String(), ")")[1:], ""), "(") { // Checks for multiply return args
				returns, err = toTypesArray(strings.Split(strings.Split(strings.Join(strings.Split(reflect.TypeOf(module).String(), "(")[2:], ""), ")")[0], ","))
				if err != nil {
					return nil, err
				}
			} else {
				// Switches through the secondary value types
				switch strings.Split(strings.Join(strings.Split(reflect.ValueOf(module).String(), ")")[1:], ""), " ")[1] {
				case "string": // String
					returns = append(returns, lexer.STRING)
				case "int": // Int
					returns = append(returns, lexer.INT)
				case "bool": // Bool
					returns = append(returns, lexer.BOOL)
				case "any":
					returns = append(returns, lexer.Any)
				}
			}
			address := make([]byte, 7) // Randomly generates a memory address
			rand.Read(address)         // Randomly updates the memory address

			// Returns the pointer of the function including the pointer with gofunction inside
			return &Pointer{address: address, memory: GoFunction{module: name, wanted: takes, rets: returns, xseal: module}}, nil
		}
	}

	return nil, errors.New("type given is not implemented in Go2Evo: " + reflect.TypeOf(module).Kind().String())
}

// toTypesArray converts the strings into lexer acceptable tokenTypes
func toTypesArray(str []string) ([]lexer.TokenType, error) {
	var assign []lexer.TokenType = make([]lexer.TokenType, 0)
	for _, arg := range str {
		arg = strings.ReplaceAll(arg, " ", "") // Removes any space inside the string

		switch arg {

		case "[]string":
			assign = append(assign, lexer.VariadicString)
		case "[]bool":
			assign = append(assign, lexer.VariadicBool)
		case "[]int":
			assign = append(assign, lexer.VariadicInt)
		case "[]any", "[]interface{}":
			assign = append(assign, lexer.VariadicAny)
		case " ", "":
			continue
		case "string": // String
			assign = append(assign, lexer.STRING)
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64": // Multiply INT support
			assign = append(assign, lexer.INT)
		case "bool": // Boolean
			assign = append(assign, lexer.BOOL)
		case "any", "interface{}":
			assign = append(assign, lexer.Any)
		default: // Non implemented
			return nil, errors.New("type not supported")
		}
	}

	return assign, nil
}
