package std

import (
	"fmt"
	"strings"
)

// Prints the object
func Print(p []string) {
	fmt.Println(strings.Join(p, " "))
}

// Prints the object but formatted
func Printf(format string, args []interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Returns the object
func Sprint(p []interface{}) string {
	return fmt.Sprint(p...)
}

// Returns the object but formatted
func Sprintf(format string, args []interface{}) string {
	return fmt.Sprintf(format, args...)
}