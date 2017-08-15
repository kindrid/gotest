package should

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/kindrid/gotest/rest"
)

// RESTHarness provides an engine that can construct requests, run them, and
// prepare the results for testing.
type RESTHarness struct {
	API       rest.Describer
	Requester RequestMaker
	Parser    StructureParser
}

// RESTExchange holds one HTTP request, expected response, and actual response
type RESTExchange struct {
	Request  *http.Request   // The request
	Expected *BodiedResponse // The response we should have got
	Actual   *BodiedResponse // The response we actually got
	Err      error           // any error running the request
}

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

// RunRequest executes an HTTP request and returns the expected and actual response in a
// *RESTExchange. For the format of params, see rest.Describer's documentation, currently:
//
// Params is a list of strings, [name1, value1, name2, value2, ...]. Keys should have one
// of these prefixes:
//
// 	  ":" - indicates an html header as a string
//    "&" - indicates a URL param as a string
//    "=" - treated as a raw string in path and body templating, ADD QUOTES if you want quotes.
func (har *RESTHarness) RunRequest(requestID string, params []string, body string) (result *RESTExchange) {
	var expected, actual *http.Response
	// Grab information from the Describer (API specification)
	result = &RESTExchange{}
	// q.Q("Before get Request", requestID, params, body)
	result.Request, expected, result.Err = har.API.GetRequest(requestID, params, body)
	if result.Err != nil {
		return
	}

	if expected != nil {
		result.Expected, result.Err = ReadResponseBody(expected, har.Parser)
		if result.Err != nil {
			return
		}
	}

	// // Run the request
	if har.Requester == nil {
		result.Err = fmt.Errorf("a RESTHarness needs a request function to run a request")
	}
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
