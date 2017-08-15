package describers

import (
	"net/http"
)

// SwaggerDescriber provides restful scenarios from a swagger file.
type SwaggerDescriber struct {
}

func (ss *SwaggerDescriber) Topics() (result []string) {
	return
}

func (ss *SwaggerDescriber) Scenarios(topicID string) (result []string) {
	return
}

func (ss *SwaggerDescriber) Requests(topicID, scenarioID string) (result []string) {
	return
}

func (ss *SwaggerDescriber) Types() (result []string) {
	return
}

func (ss *SwaggerDescriber) GetRequest(requestID string, params []string, body *string) (req *http.Request, expected *http.Response, err error) {
	return
}

func (ss *SwaggerDescriber) GetSchema(typeID string) (result *Resource) {
	return
}
