package describers

import (
	"reflect"
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

	// Test Topic (path) extraction.
	list := spec.Topics()
	expected := []string{"/pets", "/pets/{id}"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting topics \n    %#v, \ngot %#v.", expected, list)
	}

	// Test Operations extraction (filtered and unfiltered)
	list = spec.Operations("")
	expected = []string{"/pets/{id}.get", "addPet", "deletePet", "findPets"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting operations \n    %#v, \ngot %#v.", expected, list)
	}
	list = spec.Operations("/pets")
	expected = []string{"addPet", "findPets"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting operations \n    %#v, \ngot %#v.", expected, list)
	}

	// Test Scenarios extraction
	list = spec.Scenarios("addPet")
	expected = []string{"addPet.200", "addPet.default"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting scenarios \n    %#v, \ngot %#v.", expected, list)
	}

	// Test Requests extraction
	list = spec.Requests("addPet.200")
	expected = []string{"addPet.200.0"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting requests \n    %#v, \ngot %#v.", expected, list)
	}

	// Test Sample request formation
	reqID := "addPet.200.0"
	req, rsp, err := spec.GetRequest(reqID, "")
	q.Q("\n=======\n\n", req, rsp, err)
	if err != nil {
		t.Errorf("error making request %s", reqID)
	}
	if req == nil {
		t.Errorf("nil request for %s", reqID)
	}
	if rsp == nil {
		t.Errorf("nil response for %s", reqID)
	}

	// Test Types extraction
	list = spec.Types()
	expected = []string{"Error", "NewPet", "Pet"}
	if !reflect.DeepEqual(list, expected) {
		t.Errorf("Expecting types \n    %#v, \ngot %#v.", expected, list)
	}

}
