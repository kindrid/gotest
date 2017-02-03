package debug

import "testing"

func TestStackTraces(t *testing.T) {
	caller := CallerSimple(1)
	if caller == "" {
		t.Error("Call stack should not be blank.")
	}
	caller = FormattedCallStack(5)
	if caller == "" {
		t.Error("Call stack should not be blank.")
	}
}
