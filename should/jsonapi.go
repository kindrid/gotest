package should

import (
	"reflect"

	"github.com/Jeffail/gabs"
)

// This file relies on parseJSON() in json.go

// BeJSONAPIResourceIdentifier passes if actual seems to be a JSONAPI resource
// typeentifier. Unlike a full JSONAPI resource object, it has no attributes but
// MUST refer to a resource in the included list.
func BeJSONAPIResourceIdentifier(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJSONAPIResourceIdentifier expects a single string argument and passes if that argument parses as a JSONAPI resource identifier."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	return beJSONAPIResourceIdentifier(json)
}

func beJSONAPIResourceIdentifier(json *gabs.Container) (fail string) {
	var fields = []interface{}{"id", reflect.String, "type", reflect.String}
	// fmt.Print("DEBUG", json, fields)
	fail += HaveFields(json, fields...)
	fail += HaveOnlyFields(json, fields...)
	return
}

// BeJSONAPIRecord passes if actual seems to be a complete JSONAPI-format
// response where response.data is a single JSON object.
func BeJSONAPIRecord(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJSONAPIRecord expects a single string argument and passes if that argument parses as a JSONAPI single-object record."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	fail += HaveFields(json, "meta", reflect.Map, "data", reflect.Map)
	fail += HaveOnlyFields(json, "meta", reflect.Map, "data", reflect.Map, "links", reflect.Map, "included", reflect.Slice)
	fail += beValidMeta(json.Search("meta"))
	fail += beValidRecord(json.Search("data"))
	if links := json.Search("links"); links != nil {
		fail += beValidLinks(links)
	}
	if included := json.Search("included"); included != nil {
		fail += beValidIncluded(included)
	}
	return
}

func beValidMeta(json *gabs.Container) (fail string) {
	// Kindrid Specific:
	// fail += HaveFields(json, "apiVersion", reflect.String, "formatVersion", reflect.String)
	return
}

func beValidRecord(json *gabs.Container) (fail string) {
	fail += HaveFields(json, "id", reflect.String, "type", reflect.String)
	return
}

func beValidLinks(json *gabs.Container) (fail string) {
	return
}

func beValidIncluded(json *gabs.Container) (fail string) {
	return
}
