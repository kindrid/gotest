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
	Document *loads.Document
	aswagger *analysis.Spec             // analyzed swagger doc
	swagger  *spec.Swagger              // swagger doc (flattened?)
	requests map[string]*swaggerRequest // mapped by operationId
}

type swaggerRequest struct {
	path      string
	method    string
	operation *spec.Operation
}

// LoadSwaggerDescriber creates a SwaggerDescriber from a yaml or json OpenApi file.
func LoadSwaggerDescriber(filename string) (result *SwaggerDescriber, err error) {
	ss := SwaggerDescriber{requests: make(map[string]*swaggerRequest)}
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

	for method, paths := range ss.aswagger.Operations() {
		for path, op := range paths {
			if op.ID == "" {
				op.ID = path + "." + strings.ToLower(method)
			}
			if _, alreadyHas := ss.requests[op.ID]; alreadyHas {
				msg := "more than one operationId `%s`"
				return nil, fmt.Errorf(msg, op.ID)
			}
			ss.requests[op.ID] = &swaggerRequest{path, method, op}
		}
	}

	result = &ss
	return
}

// Topics lists the paths in the swagger document
func (ss *SwaggerDescriber) Topics() (result []string) {
	for path, _ := range ss.Document.Analyzer.AllPaths() {
		result = append(result, path)
	}
	return
}

// Operations lists the swagger operation IDs under a path, all if topicID is blank
func (ss *SwaggerDescriber) Operations(topicID string) (result []string) {
	for opID, req := range ss.requests {
		if topicID == "" || topicID == req.path {
			result = append(result, opID)
		}
	}
	sort.Strings(result)
	return
}

// Scenarios lists the different status code examples of an operation (swagger verb + path)
func (ss *SwaggerDescriber) Scenarios(operationID string) (result []string) {
	return
}

// Requests lists the request(s) that make up a scenario.
// Swagger 2.0 only supports one request per scenario
func (ss *SwaggerDescriber) Requests(topicID, scenarioID string) (result []string) {
	// return ss.Document.Analyzer.OperationMethodPaths()
	return

}

// Types lists the embedded type definitions in the swagger spec.
func (ss *SwaggerDescriber) Types() (result []string) {
	return ss.Document.Analyzer.AllDefinitionReferences()
}

// GetRequest implements the Describer interface
func (ss *SwaggerDescriber) GetRequest(requestID string, params []string, body string) (req *http.Request, expected *http.Response, err error) {
	return
}

// GetSchema implements the Describer interface
func (ss *SwaggerDescriber) GetSchema(typeID string) (result *Resource) {
	return
}
