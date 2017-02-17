package should

import (
	"strings"
	"testing"

	"github.com/kindrid/gotest/debug"
	"github.com/y0ssar1an/q"
)

// HELPERS for self testing

func Passes(t *testing.T, topic string, a Assertion, actual interface{}, expected ...interface{}) {
	fail := a(actual, expected...)
	if fail == "" {
		return
	}
	q.Q("test debug code check")
	t.Errorf("%s Expected %#v to pass but got '%s'", topic, a, fail)
	t.Errorf(strings.Join(debug.CallStack(5), ", "))
	// if fail != "" {
	// 	t.Errorf("Expected %v to pass. Instead got '%s'.", a, fail)
	// } else if pass == "" {
	// 	t.Errorf("Expected %v to pass, but got a blank pass messasge.", a)
	// }
}

func Fails(t *testing.T, topic string, a Assertion, actual interface{}, expected ...interface{}) {
	fail := a(actual, expected...)
	if fail == "" {
		t.Errorf("Expected %#v (%s) to fail. Instead it passed.", a, topic)
	}
}

func TestAssertion(t *testing.T) {
	t.Run("Assertion fundamentals", testFundamentals)
}

func testFundamentals(t *testing.T) {
	// Forced passes and failures act correctly
	Passes(t, "Forced Pass", AlwaysPass, "PASS!", nil, nil)
	Fails(t, "Forced Fail", AlwaysFail, "FAIL!", nil, nil)

	// Negations correctly invert forced passes and failures
	Passes(t, "Negated Forced Fail", Not(AlwaysFail), "FAIL!", nil, nil)
	Fails(t, "Negated Forced Pass", Not(AlwaysPass), "PASS!", nil, nil)
}
