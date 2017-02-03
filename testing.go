package gotest

import (
	"github.com/kindrid/gotest/debug"
	"github.com/kindrid/gotest/should"
)

// StackDepth sets the maximum stack depth reported with errors. 0 disables.
var StackDepth = 5

// T describes the interface provided by Go's std.testing.T. If only they had
// made that an interface!
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
}

// Assert wraps any standard Assertion for use with Go's std.testing library.
func Assert(t T, actual interface{}, assertion should.Assertion, expected ...interface{}) {
	fail := assertion(actual, expected...)
	if fail != "" {
		t.Errorf("%s\nStack=%s", fail, debug.FormattedCallStack(StackDepth))
	}
}

// Deny negates any standard Assertion for use with Go's std.testing library.
// You may also want to use should.Not()
func Deny(t T, actual interface{}, assertion should.Assertion, expected ...interface{}) {
	fail := assertion(actual, expected...)
	if fail == "" {
		t.Error("Expected a failure")
	}
}
