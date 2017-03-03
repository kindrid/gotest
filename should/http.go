package should

import (
	"fmt"
	"net/http"
)

// DescribeResponse give tests a way to, if failing, dump the request and
// response without having every test method do it.
func DescribeResponse(rsp *http.Response) (desc string) {
	return "WIP--see goblue.blueprint.Harness.formatResponse"
}

// MatchHTTPStatusCode asserts that the documented and actual HTTP status codes match
func MatchHTTPStatusCode(actual interface{}, expected ...interface{}) (fail string) {
	if msg := exactly(1, expected); msg != Ok {
		return msg
	}
	eStatus, ok := expected[0].(int)
	if !ok {
		return fmt.Sprintf("Expected value should be an int, not a %T.", expected[0])
	}
	aRsp, ok := actual.(*http.Response)
	if !ok {
		return "Actual value should be a *http.Response"
	}
	aStatus := aRsp.StatusCode
	if Equal(aStatus, eStatus) == "" {
		return ""
	}

	return fmt.Sprintf(
		"HTTP DEBUG Status Expected: %d %s. Got: %d %s.",
		eStatus, http.StatusText(int(eStatus)),
		aStatus, http.StatusText(aStatus),
	)
}
