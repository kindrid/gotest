package gotest

import "github.com/kindrid/gotest/should"

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
		// q.Q("Value L?", CallerSimple(3), fail)
		// t.Error(CallerSimple())
		t.Error(fail)
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
