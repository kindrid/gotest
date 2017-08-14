package rest

import (
	"testing"

	"github.com/kindrid/gotest/rest/describers"
)

func TestDescriber(t *testing.T) {
	swag := &describers.SwaggerDescriber{}
	func(_ Describer) {
		t.Log("yep, it's a describer")
	}(swag)

}
