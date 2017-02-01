package should

import (
	"fmt"
	"net/http"
)

// MatchHTTPStatusCode asserts that the documented and actual HTTP status codes match
func MatchHTTPStatusCode(actual interface{}, expected ...interface{}) (fail string) {
	if msg := exactly(1, expected); msg != Ok {
		return msg
	}
	eStatus, ok := expected[0].(int)
	if !ok {
		return "Actual value should be an integer"
	}
	aRsp, ok := actual.(*http.Response)
	if !ok {
		return "Actual value should be a *http.Response"
	}
	aStatus := aRsp.StatusCode
	if aStatus == eStatus {
		return ""
	}

	return fmt.Sprintf(
		"HTTP Status Expected: %d %s. Got: %d %s. \nRequest: %#v.\nResponse: %#v.",
		eStatus, http.StatusText(eStatus),
		aStatus, http.StatusText(aStatus),
		aRsp.Request, aRsp,
	)
}
