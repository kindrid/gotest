package should

import "testing"

func TestRestfulHarness(t *testing.T) {
	har := RESTHarness{}
	if har.Parser != nil {
		t.Error("We're just making things compile.")
	}
}
