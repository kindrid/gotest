package rest

// ScenarioStore holds a series of requests and responses and can retrieve
// them for testing
type ScenarioStore interface {
	Topics() []*Topic
	// GetTopic(id string) *Topic
	// GetScenario(id string) *Scenario
}
