package debug

import (
	"fmt"
	"runtime"
	"strings"
)

// CallerInfo gives info about the current call stack.
func CallerInfo(depth int) (msg, fileName string, fileLine int) {
	_, fileName, fileLine, ok := runtime.Caller(depth)
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

// CallStack gives a short array of the call stack descriptions.
func CallStack(maxDepth int) (result []string) {
	const offset = 1
	result = make([]string, maxDepth)
	for i := offset; i < maxDepth+offset; i++ {
		msg, _, _ := CallerInfo(i)
		if msg == "" {
			break
		}
		result = append(result, msg)
	}
	return
}

// FormattedCallStack returns the call stack printout as lines.
func FormattedCallStack(maxDepth int) string {
	return strings.Trim(strings.Join(CallStack(maxDepth), "\n"), " \n\r\t")
}
