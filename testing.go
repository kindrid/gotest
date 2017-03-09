package gotest

import (
	"flag"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/kindrid/gotest/debug"
	"github.com/kindrid/gotest/should"
)

// T describes the interface provided by Go's std.testing.T. If only they had
// made that an interface!
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow() // exit the current test immediately
	Logf(format string, args ...interface{})
	Name() string // need Go 1.8 for this.
}

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

// StackDepth sets the maximum stack depth reported with errors. 0 disables. It
// is puporsefully public so tests using this library can manipulate it and
// check it.
var StackDepth int

// Verbosity sets a level of "chattiness" for the tests. It is puporsefully
// public so tests using this library can manipulate it and check it.
var Verbosity int

var failFast bool

func init() {
	flag.IntVar(&StackDepth, "gotest-stack", 0, "stack-trace depth on failure")
	flag.IntVar(&Verbosity, "gotest-verbosity", 0, "verbosity level: -1=silent, 0=short, 1=long, 2=show-actuals, \n\t3=show-expecteds, 4=show-debugging-details, 5=show-test-internals")
	flag.BoolVar(&failFast, "gotest-failfast", false, "cause tests to exit with errorcode=1 after the first assertion failure")
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
// We're considering using spew or a similar library to give verbosity.
func Inspectv(minLevel int, label string, inspected ...interface{}) (result string) {
	if !Vocal(minLevel) {
		return
	}
	if label != "" {
		result = fmt.Sprintf("%s: \n", label)
	}
	// for _, x := range inspected {
	// result += fmt.Sprintf("%#v\n", x)
	// }
	result += spew.Sdump(inspected...)
	return
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
		name := t.Name()
		msg += Sprintv(Short, "Failed %s: %s", name, terseMsg)
		msg += Sprintv(Long, "\n%s\n", extraMsg)
		msg += Inspectv(Actuals, "\nDUMP OF ACTUAL VALUE", actual)
		msg += Inspectv(Expecteds, "\nDUMP OF EXPECTED VALUE", expected)
		if detailsMsg != "" {
			msg += Sprintv(Debug, "\nFAILURE DETAILS: %s\n", detailsMsg)
		}
		if metaMsg != "" {
			msg += Sprintv(Insane, "\nINTERNALS (FOR DEBUGGING ASSERTIONS): %s\n", metaMsg)
		}
		if failFast {
			msg += "\nNOTE: skipping remaining assertions for this test because of --gotest-failfast."
			t.Error(msg)
			t.FailNow()
		}
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
