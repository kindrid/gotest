package gotest

import (
	"fmt"
	"testing"
)

func TestDebug(t *testing.T) {
	fmt.Println(CallerSimple(1))
}
