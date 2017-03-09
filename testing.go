package gotest

import (
	"flag"
	"fmt"

	"github.com/kindrid/gotest/debug"
	"github.com/kindrid/gotest/should"
)

// StackDepth sets the maximum stack depth reported with errors. 0 disables. It
// is puporsefully public so tests using this library can manipulate it and
// check it.
var StackDepth int

// Verbosity sets a level of "chattiness" for the tests. It is puporsefully
// public so tests using this library can manipulate it and check it.
var Verbosity int

// // Verbosity Levels: these are conventions only. Assertions and test
// functions can interpret these however they want. Stack traces, however should
// be controlled by the StackDepth variable.
const (
	Silent    = iota - 1 // (hopefully) no output except panics
	Short                // Shows the first line of the failure string (up to \n) and info to re-run the particular tests
	Long                 // Add the next level of failure string (up to \n\n)
	Actuals              // Add the entire failure string and a (possibly shortened) representation of the actual value
	Expecteds            // Add  (possibly shortened) representation(s) of  expected values.
	Debug                // Adds granular information to diagnose the failure and ensures that all representations are unabridged. This level may inject flags into the tested item to make it more verbose.
	Insane               // Adds information to test meta concerns, such as logic within assertions.
)

func init() {
	flag.IntVar(&StackDepth, "gotest-stack", 0, "stack-trace depth on failure")
	flag.IntVar(&Verbosity, "gotest-verbosity", 0, "verbosity level: 0=silent, 1=short, 2=long, 3=show-actuals, \n\t4=show-expecteds, 5=show-debugging-details, 6=show-test-internals")
}

// T describes the interface provided by Go's std.testing.T. If only they had
// made that an interface!
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	Logf(format string, args ...interface{})
}

// Vocal makes an easy way to gate operations by verbosity level. It returns true if Verbosity is < minLevel.
func Vocal(minLevel int) bool {
	return Verbosity >= minLevel
}

// Sprintv formats a string if Verbosity >= minLevel, otherwise returns ""
func Sprintv(minLevel int, format string, args ...interface{}) string {
	if !Vocal(minLevel) {
		return ""
	}
	return fmt.Sprintf(format, args...)
}

// Inspectv returns a detailed introspection of objects if Verbosity >= minLevel.
func Inspectv(minLevel int, label string, inspected ...interface{}) (result string) {
	if !Vocal(minLevel) {
		return
	}
	if label != "" {
		result = fmt.Sprintf("%s: \n", label)
	}
	return // check out spew
}

// Assert wraps any standard Assertion for use with Go's std.testing library.
func Assert(t T, actual interface{}, assertion should.Assertion, expected ...interface{}) {
	fail := assertion(actual, expected...)
	if fail != "" {
		terseMsg, extraMsg, detailsMsg, metaMsg := should.SplitMsg(fail)
		msg := ""
		if StackDepth > 0 {
			msg += fmt.Sprintf("\nTest Failure Stack Trace: %s\n\n", debug.FormattedCallStack(StackDepth))
		}
		msg += Sprintv(Short, "Test failure: %s.\nTest path: %s\n", terseMsg, "testPath")
		msg += Sprintv(Long, "%s\n", extraMsg)
		msg += Inspectv(Actuals, "Actual", actual)
		msg += Inspectv(Expecteds, "Expected", expected)
		msg += Sprintv(Debug, "Failure Details: %s\n", detailsMsg)
		msg += Sprintv(Insane, "Meta Details: %s\n", metaMsg)
		t.Error(msg)
	}
}

// Deny negates any standard Assertion for use with Go's std.testing library.
// You may also want to use should.Not()--it will give more accurate reporting.
func Deny(t T, actual interface{}, assertion should.Assertion, expected ...interface{}) {
	fail := assertion(actual, expected...)
	if fail == "" {
		t.Error("Expected a failure")
	}
}

// Later describes pending tests.
func Later(t T, desc string, ignored ...interface{}) {
	t.Logf("LATER: %s", desc)
}
