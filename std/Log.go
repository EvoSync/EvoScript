package std

import (
	"fmt"
	"log"
	"strings"
)

func Log(logs []string) {
	log.Println(strings.Join(logs, " "))
}

func Logf(format string, v []any) {
	log.Println(fmt.Sprintf(format, v...))
}