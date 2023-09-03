package std

import "reflect"

func Type(i any) string {
	return reflect.TypeOf(i).Kind().String()
}