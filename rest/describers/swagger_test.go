package describers

import (
	"reflect"
	"testing"
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
		t.Errorf("Expecting topics %v, got %v.", expected, list)
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

	// // Test Scenarios extraction (filtered and unfiltered)
	// list = spec.Operations("")
	// expected = []string{"/pets/{id}.get", "addPet", "deletePet", "findPets"}
	// if !reflect.DeepEqual(list, expected) {
	// 	t.Errorf("Expecting operations \n    %#v, \ngot %#v.", expected, list)
	// }

}
