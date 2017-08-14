package should

import (
	"bytes"
	"net/http"

	"github.com/kindrid/gotest/rest"
)

// Provides a harness around REST API descriptions

// RequesterMaker is a function that the tested code will use to simulate or actually perform a request.
type RequestMaker func(*http.Request) (*http.Response, error)

// BodiedResponse holds a response with the Body already read and parsed
type BodiedResponse struct {
	Response *http.Response    // incorporate response
	Raw      string            // the raw body
	Parsed   StructureExplorer // parsed body
}

func ReadResponseBody(rsp *http.Response, parser StructureParser) (result *BodiedResponse, err error) {
	result = &BodiedResponse{Response: rsp}
	buf := new(bytes.Buffer)
	buf.ReadFrom(rsp.Body)
	result.Raw = buf.String()
	result.Parsed, err = parser(result.Raw)
	return
}

// Exchange holds one HTTP request, expected response, and actual response
type Exchange struct {
	Request  *http.Request   // The request
	Expected *BodiedResponse // The response we should have got
	Actual   *BodiedResponse // The response we actually got
	Err      error           // any error running the request
}

// RESTHarness provides an engine that can construct requests, run them, and
// prepare the results for testing.
type RESTHarness struct {
	API       rest.Describer
	Requester RequestMaker
	Parser    StructureParser
}

// RunRequest executes an HTTP request and returns the expected and actual response in an *Exchange
func (har *RESTHarness) RunRequest(requestID string, params map[string]string, body *string) (result *Exchange) {
	var expected, actual *http.Response
	// Grab information from the Describer (API specification)
	result.Request, expected, result.Err = har.API.GetRequest(requestID, params, body)
	if result.Err != nil {
		return
	}
	result.Expected, result.Err = ReadResponseBody(expected, har.Parser)
	if result.Err != nil {
		return
	}

	// Run the request
	actual, result.Err = har.Requester(result.Request)
	if result.Err != nil {
		return
	}
	result.Actual, result.Err = ReadResponseBody(actual, har.Parser)
	if result.Err != nil {
		return
	}

	return
}
