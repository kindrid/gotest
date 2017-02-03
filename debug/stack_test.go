package debug

import (
	"testing"

	"github.com/kindrid/gotest"
	"github.com/kindrid/gotest/should"
)

func TestStackTraces(t *testing.T) {
	caller := CallerSimple(1)
	gotest.Assert(t, caller, should.NotBeBlank)
	gotest.Assert(t, CallStack(5)[4], should.BeBlank)
	// gotest.Assert(t, CallStack(5)[0], should.NotBeBlank)  // TODO check this
}
