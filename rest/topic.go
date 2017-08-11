package rest

import "net/http"

// Topic holds sets of scenarios--such as all the requests for a resource
type Topic struct {
	ID        string      // often the resource name
	Scenarios []*Scenario //
}

// Scenario holds a logical set of restful exchanges
type Scenario struct {
	ID        string      // often verb plus endpoint name
	Exchanges []*Exchange // the requests to run and expected responses
}

// Exchange holds a Restful Exchange
type Exchange struct {
	ID       string         // if blank, assigned 0-based indices
	Request  *http.Request  // the request to run
	Response *http.Response // the response
}
