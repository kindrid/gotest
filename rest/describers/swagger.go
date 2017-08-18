package describers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// SwaggerDescriber provides restful scenarios from a swagger file.
type SwaggerDescriber struct {
	aswagger   *analysis.Spec             // analyzed swagger doc
	swagger    *spec.Swagger              // swagger doc (flattened?)
	operations map[string]*swaggerRequest // operationID to partial swagger request
	scenarios  map[string]*swaggerRequest // scenarioID to partial swagger request
	requests   map[string]*swaggerRequest // requestID to swagger request
	urlBase    string                     // e.g. http://localhost:8081
}

type swaggerRequest struct {
	method     string
	path       string
	operation  *spec.Operation
	scenarioID string
	order      int // ordinal location within the scenario
	code       int
	response   *spec.Response
}

// LoadSwaggerDescriber creates a SwaggerDescriber from a yaml or json OpenApi file.
func LoadSwaggerDescriber(filename string) (result *SwaggerDescriber, err error) {
	var doc, doc2 *loads.Document
	// Read the document
	if doc, err = loads.Spec(filename); err != nil {
		return
	}

	// Check that we actually got well-formed swagger
	if ver := doc.Spec().Swagger; ver != "2.0" {
		return nil, fmt.Errorf("expected version 2.0 swagger, got '%s'", ver)
	}

	// Turn $ref's into literals.
	// Note analysis.Flatten looks like it does a more complete job
	opts := &spec.ExpandOptions{}
	if doc2, err = doc.Expanded(opts, opts); err != nil { //opts repeated because of a bug in https://github.com/go-openapi/loads/blob/master/spec.go#L189 easy fix, but in the middle of POCing this
		return
	}

	// Construct the object
	ss := SwaggerDescriber{
		aswagger:   doc2.Analyzer,
		swagger:    doc2.Spec(),
		operations: make(map[string]*swaggerRequest),
		scenarios:  make(map[string]*swaggerRequest),
		// scenarios:  make(map[string]string),
		// requests:  make(map[string]string),
		requests: make(map[string]*swaggerRequest),
	}
	if err = ss.populate(); err == nil {
		result = &ss
	}
	return
}

func (ss *SwaggerDescriber) populate() (err error) {
	for method, paths := range ss.aswagger.Operations() {
		for path, op := range paths {
			// Ensure a unique operation ID
			if op.ID == "" {
				op.ID = path + "." + strings.ToLower(method)
			}
			if _, alreadyHas := ss.operations[op.ID]; alreadyHas {
				msg := "more than one operationId `%s`"
				return fmt.Errorf(msg, op.ID)
			}

			// Add Operation
			operation := swaggerRequest{path: path, method: method, operation: op}
			ss.operations[op.ID] = &operation

			if err = ss.populateScenarios(operation); err != nil {
				return
			}
		}
	}
	return nil
}

func (ss *SwaggerDescriber) populateScenarios(template swaggerRequest) (err error) {
	op := template.operation
	if op.Responses.Default != nil {
		rsp := op.Responses.Default
		sc := template
		sc.scenarioID = fmt.Sprintf("%s.%s", sc.operation.ID, "default")
		sc.code = 0
		ss.scenarios[sc.scenarioID] = &sc
		ss.populateRequests(sc, []*spec.Response{rsp})
	}
	for status, rsp := range op.Responses.StatusCodeResponses {
		sc := template
		sc.scenarioID = fmt.Sprintf("%s.%d", sc.operation.ID, status)
		sc.code = status
		ss.scenarios[sc.scenarioID] = &sc
		ss.populateRequests(sc, []*spec.Response{&rsp})
	}
	return
}

func (ss *SwaggerDescriber) populateRequests(template swaggerRequest, responses []*spec.Response) (err error) {
	// Populate Requests (for now, one per scenario)
	for i, rsp := range responses {
		rq := template
		rq.response = rsp
		rq.order = i
		requestID := fmt.Sprintf("%s.%d", rq.scenarioID, i)
		ss.requests[requestID] = &rq
	}
	return
}

// Topics lists the paths in the swagger document
func (ss *SwaggerDescriber) Topics() (result []string) {
	for path := range ss.aswagger.AllPaths() {
		result = append(result, path)
	}
	sort.Strings(result)
	return
}

// Operations lists the swagger operation IDs under a path, all if topicID is blank
func (ss *SwaggerDescriber) Operations(topicID string) (result []string) {
	for opID, req := range ss.operations {
		if topicID == "" || topicID == req.path {
			result = append(result, opID)
		}
	}
	sort.Strings(result)
	return
}

// Scenarios lists the different status code responses of an operation (swagger verb + path)
func (ss *SwaggerDescriber) Scenarios(operationID string) (result []string) {
	for scID, req := range ss.scenarios {
		if operationID == "" || operationID == req.operation.ID {
			result = append(result, scID)
		}
	}
	sort.Strings(result)
	return
}

// Requests lists the request(s) that make up a scenario.
// Swagger 2.0 only supports one request per return code, but some systems may
// allow more to represent tests
func (ss *SwaggerDescriber) Requests(scenarioID string) (result []string) {
	for id, req := range ss.requests {
		if scenarioID == "" || scenarioID == req.scenarioID {
			result = append(result, id)
		}
	}
	sort.Strings(result)
	return
}

// Types lists the embedded type definitions in the swagger spec.
func (ss *SwaggerDescriber) Types() (result []string) {
	for name := range ss.swagger.Definitions {
		result = append(result, name)
	}
	sort.Strings(result)
	return
}

// GetRequest implements the Describer interface
func (ss *SwaggerDescriber) GetRequest(requestID string, body string, params ...string) (
	req *http.Request, expected *http.Response, err error) {
	sreq, ok := ss.requests[requestID]
	if !ok {
		err = fmt.Errorf("cannot find requestID==%s", requestID)
		return
	}
	if req, err = ss.makeRequest(sreq, body, params...); err != nil {
		return
	}
	expected, err = ss.makeResponse(sreq)
	return
}

func (ss *SwaggerDescriber) makeRequest(sr *swaggerRequest, body string, params ...string) (result *http.Request, err error) {
	url := ss.urlBase + sr.path // TODO apply path param and query params
	result, err = http.NewRequest(sr.method, url, strings.NewReader(body))
	if err != nil {
		return
	}
	// TODO apply all the layers of headers
	// any other param locations for swagger?
	return
}

func (ss *SwaggerDescriber) makeResponse(sr *swaggerRequest) (result *http.Response, err error) {
	result = &http.Response{
		Status:     http.StatusText(sr.code),
		StatusCode: sr.code,
	}
	return
}

// GetSchema implements the Describer interface
func (ss *SwaggerDescriber) GetSchema(typeID string) (result *Resource) {
	// schema := ss.swagger.Definitions[typeID]
	// q.Q(schema)
	return
}
