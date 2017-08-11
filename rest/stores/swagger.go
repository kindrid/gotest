package stores

import "github.com/kindrid/gotest/rest"

// SwaggerScenarioStore provides restful scenarios from a swagger file.
type SwaggerScenarioStore struct {
	topics []*rest.Topic
}

func (ss *SwaggerScenarioStore) Topics() {
	return ss.topics
}
