package describers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// SwaggerDescriber provides restful scenarios from a swagger file.
type SwaggerDescriber struct {
	Document   *loads.Document
	aswagger   *analysis.Spec             // analyzed swagger doc
	swagger    *spec.Swagger              // swagger doc (flattened?)
	operations map[string]*swaggerRequest // mapped by operationID
	scenarios  map[string]string          // scenarioID to operationID
	requests   map[string]string          // requestID to scenarioID

}

type swaggerRequest struct {
	path      string
	method    string
	operation *spec.Operation
	code      int
	response  *spec.Response
}

// LoadSwaggerDescriber creates a SwaggerDescriber from a yaml or json OpenApi file.
func LoadSwaggerDescriber(filename string) (result *SwaggerDescriber, err error) {
	ss := SwaggerDescriber{
		operations: make(map[string]*swaggerRequest),
		scenarios:  make(map[string]string),
		requests:   make(map[string]string),
	}
	var doc *loads.Document
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
	if ss.Document, err = doc.Expanded(opts, opts); err != nil { //opts repeated because of a bug in https://github.com/go-openapi/loads/blob/master/spec.go#L189 easy fix, but in the middle of POCing this
		return
	}

	ss.Document = doc
	ss.aswagger = doc.Analyzer
	ss.swagger = doc.Spec()

	if err = ss.populate(); err == nil {
		result = &ss
	}
	return
}

func (ss *SwaggerDescriber) populate() error {
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
			ss.operations[op.ID] = &swaggerRequest{path: path, method: method, operation: op}

			// Populate Scenarios
			if op.Responses != nil {
				if op.Responses.Default != nil {
					ss.scenarios[op.ID+".default"] = op.ID
				}
				for status := range op.Responses.StatusCodeResponses {
					ss.scenarios[op.ID+"."+strconv.Itoa(status)] = op.ID
				}
			}

			// Populate Requests (for now, one per scenario)
			for sc := range ss.scenarios {
				ss.requests[sc+".0"] = sc
			}
		}
	}
	return nil
}

// Topics lists the paths in the swagger document
func (ss *SwaggerDescriber) Topics() (result []string) {
	for path := range ss.Document.Analyzer.AllPaths() {
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
	for sc, op := range ss.scenarios {
		if operationID == "" || operationID == op {
			result = append(result, sc)
		}
	}
	sort.Strings(result)
	return
}

// Requests lists the request(s) that make up a scenario.
// Swagger 2.0 only supports one request per return code, but some systems may
// allow more to represent tests
func (ss *SwaggerDescriber) Requests(scenarioID string) (result []string) {
	for req, sc := range ss.requests {
		if scenarioID == "" || scenarioID == sc {
			result = append(result, req)
		}
	}
	sort.Strings(result)
	return
}

// Types lists the embedded type definitions in the swagger spec.
func (ss *SwaggerDescriber) Types() (result []string) {

	for name, _ := range ss.swagger.Definitions {
		result = append(result, name)
	}
	sort.Strings(result)
	return
}

// GetRequest implements the Describer interface
func (ss *SwaggerDescriber) GetRequest(requestID string, body string, params ...string) (
	req *http.Request, expected *http.Response, err error) {
	// method, path, op, code, rsp, err := ss.getRequestParents(requestID)
	return
}

func (ss *SwaggerDescriber) getRequestParents(requestID string) (
	method, path string, op *spec.Operation, code int, rsp *spec.Response, err error,
) {
	scenarioID := ss.requests[requestID]
	operation := ss.operations[scenarioID]
	method = operation.method
	path = operation.path
	op = operation.operation

	requestPortion := requestID[len(scenarioID)+1:]
	code, err = strconv.Atoi(requestPortion)
	if err == nil {
		if code == 0 { // default response
			rsp = op.Responses.Default
		} else {
			rspz := (op.Responses.StatusCodeResponses[code])
			rsp = &rspz
		}
	}
	return
}

// GetSchema implements the Describer interface
func (ss *SwaggerDescriber) GetSchema(typeID string) (result *Resource) {
	// schema := ss.swagger.Definitions[typeID]
	// q.Q(schema)
	return
}
