package rest

import (
	"net/http"

	"github.com/kindrid/gotest/rest/describers"
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
	Topics() (TopicIDs []string)
	Scenarios(topicID string) (ScenarioIDs []string)
	Requests(topicID, scenarioID string) (RequestIDs []string)
	Types() (typeIDs []string)

	// Get request applies any params to path and query, returning a request and the expected response.
	GetRequest(requestID string, params map[string]string, body *string) (req *http.Request, expected *http.Response, err error)
	GetSchema(typeID string) *describers.Resource
}
