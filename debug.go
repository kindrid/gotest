package gotest

import (
	"fmt"
	"runtime"
)

// CallerInfo gives info about the current call stack.
func CallerInfo(depth int) (msg, fileName string, fileLine int) {
	_, fileName, fileLine, ok := runtime.Caller(depth)

	msg = "Caller info not found."
	if ok {
		msg = fmt.Sprintf("%s:%d", fileName, fileLine)
	}
	return
}

// CallerSimple gives a single string with condensed information about the current call stack.
func CallerSimple(depth int) string {
	msg, _, _ := CallerInfo(depth)
	return msg
}
