package rest

import (
	"net/http"
)

/* Describer holds a description of an API in a hierarchy:

- Topics: usually hold broad types of resources and contain
- Scenarios: series of actions towards a single goal or example
- Requests: pairs of requests and expected responses

Each of these items, Topics, Scenarios, and Requests must have a unique id. The
id may be blank if it is the only entity of its type.

For methods that filter, such as Requests(topicId, scenarioId), "" means all the
topic or scenarios that exist.

*/
type Describer interface {
	// usually paths or resource types (first path segment)
	Topics() (TopicIDs []string)

	// Operations usually HTTPVERB + path or an id indicating it
	Operations(topicId string) (OperationIDs []string)

	// Scenarios varios status responses for an endpoint, but for embedded tests, can also be richer
	Scenarios(operationID string) (ScenarioIDs []string)

	// Requests returns a list of id's for items that points to an actual Request
	// Response pair. They may also
	// Requests incorporate mimetype, otherwise the first example is used.
	Requests(scenarioID string) (RequestIDs []string)

	// Types names of structures defined in the specification
	Types() (typeIDs []string)

	// GetRequest applies any params to path, query, and body template,
	// returning a request and the expected response.
	//
	// requestID: a string that uniquely identifies a request response pair
	// params:  a list of strings, [key1, value1, key2, value2, ...].
	//
	// Keys must have one of these prefixes:
	//
	// 	  ":" - indicates an html header as a string
	//    "&" - indicates a URL param as a string
	//    "=" - treated as a raw string in path and body templating, ADD QUOTES if you want quotes.
	GetRequest(requestID string, body string, params ...string) (
		req *http.Request, expected *http.Response, err error)

	// GetSchema gets a described schema from the specification.  Might actually be better to
	// pass in a structure explorer and compare, since Resource is one level only. Otherwise,
	// what format is the resource schema? Could use the custom Resource, but as mentioned, that's
	// one level only. There's github.com/gima/govalid which is nice and Go-like. There's
	// https://github.com/xeipuuv/gojsonschema which introduces the intricacies of
	// jsonschema. There's the github.com/go-openapi/spec.Schema which is pre-built
	// to handle swagger, but might not handle other things.
	// for now, returns an interface to be interpreted however
	GetSchema(typeID string) (interface{}, error)
}
