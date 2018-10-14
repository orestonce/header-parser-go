package ymdCppHeaderParser

import "fmt"

func assert(b bool, a ... interface{}) {
	if b == false {
		panic(fmt.Sprintln(a...))
	}
}
