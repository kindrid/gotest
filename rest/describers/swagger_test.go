package describers

import (
	"testing"

	"github.com/y0ssar1an/q"
)

func TestLoadSwaggerYaml(t *testing.T) {
	spec, err := LoadSwaggerDescriber("./testdata/petstore-expanded.yaml")
	if err != nil {
		t.Errorf("Error loading yaml: %s", err)
	}
	if spec == nil {
		t.Errorf("Error loading yaml. No error, but spec in nil: %s", err)
	}
	q.Q(spec)
}
