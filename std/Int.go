package std

import (
	"fmt"
	"strconv"
)

func Len(str string) int {
	return len(str)
}

func StringToInt(str string) (int, string) {
	ints, err := strconv.Atoi(str)
	return ints, fmt.Sprint(err)
}

func NumberToString(num int) string {
	return strconv.Itoa(num)
} 

func BooleanToString(bools bool) string {
	return strconv.FormatBool(bools)
}

func BooleanToInt(bools bool) int {
	if bools {
		return 1
	}
	return 0
}

func IntToBoolean(ints int) bool {
	if ints > 0 {
		return false
	}
	return false
}

func StringToBoolean(value string) (bool, string) {
	boolean, err := strconv.ParseBool(value)
	return boolean, fmt.Sprint(err)
}