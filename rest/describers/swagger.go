package describers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// SwaggerDescriber provides restful scenarios from a swagger file.
type SwaggerDescriber struct {
	Document *loads.Document
}

// LoadSwaggerDescriber creates a SwaggerDescriber from a yaml or json OpenApi file.
func LoadSwaggerDescriber(filename string) (result *SwaggerDescriber, err error) {
	ss := SwaggerDescriber{}
	var doc *loads.Document
	// Read the document
	if doc, err = loads.Spec(filename); err != nil {
		return
	}
	// Check that we actually got well-formed swagger
	if ver := doc.Spec().Swagger; ver != "2.0" {
		return nil, fmt.Errorf("expected version 2.0 swagger, got '%s'", ver)
	}
	// Turn $ref's into literals  Panics on petstore sample
	opts := &spec.ExpandOptions{}
	if ss.Document, err = doc.Expanded(opts, opts); err != nil { //opts repeated because of a bug in https://github.com/go-openapi/loads/blob/master/spec.go#L189 easy fix, but in the middle of POCing this
		return
	}
	ss.Document = doc
	result = &ss
	return
}

// Topics implements the Describer interface
func (ss *SwaggerDescriber) Topics() (result []string) {
	result = append(result, "")
	return
}

// Scenarios implements the Describer interface
func (ss *SwaggerDescriber) Scenarios(topicID string) (result []string) {
	for path, _ := range ss.Document.Analyzer.AllPaths() {
		result = append(result, path)
	}
	return
}

// Requests implements the Describer interface
func (ss *SwaggerDescriber) Requests(topicID, scenarioID string) (result []string) {
	return ss.Document.Analyzer.OperationIDs()
	// return ss.Document.Analyzer.OperationMethodPaths()
}

// Types implements the Describer interface
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
