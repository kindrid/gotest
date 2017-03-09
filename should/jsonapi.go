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
	fail = HaveFields(json, fields...)
	if fail == "" {
		fail += HaveOnlyFields(json, fields...)
	}
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

	// var links, included string
	// These aren't checked yet, so this throws an ineffectual assignment warning
	// if element := json.Search("links"); element != nil {
	// 	links = BeValidLinks(element)
	// }
	// if element := json.Search("included"); element != nil {
	// 	included = BeValidIncluded(element)
	// }

	fail = FailFirst(
		NotJSONAPIError(json),
		HaveFields(json, "meta", reflect.Map, "data", reflect.Map),
		HaveOnlyFields(json, "meta", reflect.Map, "data", reflect.Map, "links", reflect.Map, "included", reflect.Slice),
		BeValidRecord(json.Search("data")),
		HaveOnlyCamelcaseKeys(json),
		// BeValidMeta(json.Search("meta")),
		// links,
		// included,
	)
	return
}

// NotJSONAPIError returns a non-blank string if the JSON appears to be a
// JSONAPI error response
func NotJSONAPIError(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJSONAPIArray expects a single string argument and passes if that argument parses as a JSONAPI multi-object array."
	if actual == nil {
		return usage
	}
	json, err := ParseJSON(actual)
	if err != nil {
		return err.Error()
	}
	if json.PathExists("errors") {
		fail = "You expected a jsonapi return, but you got what looks like a jsonapi error."
	}
	if json.PathExists("error") {
		fail += " You got what looks like a malformed jsonapi error. (It should only have .errors, an array, instead of .error, a singleton.)"
	}
	return
}

// BeJSONAPIArray passes if actual seems to be a complete JSONAPI-format
// response where response.data is a multi-object array.
func BeJSONAPIArray(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJSONAPIArray expects a single string argument and passes if that argument parses as a JSONAPI multi-object array."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}

	FailFirst(
		NotJSONAPIError(json),
		HaveFields(json, "meta", reflect.Map, "data", reflect.Slice),
		HaveOnlyFields(json, "meta", reflect.Map, "data", reflect.Slice, "links", reflect.Map, "included", reflect.Slice),
		BeValidRecordArray(json.Search("data")),
		// check links
		// check includes
	)

	return
}

// BeJSONAPI passes if actual seems to be a complete JSONAPI-format response
// where response.data is a multi-object array or a single JSON object.
func BeJSONAPI(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJSONAPIArray expects a single string argument and passes if that argument parses as a JSONAPI return value."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	if HaveFields(json, "data", reflect.Slice) == "" {
		return BeJSONAPIArray(actual, expected)
	}
	return BeJSONAPIRecord(actual, expected)
}

// func BeValidMeta(json *gabs.Container) (fail string) {
// 	// Kindrid Specific:
// 	// fail += HaveFields(json, "apiVersion", reflect.String, "formatVersion", reflect.String)
// 	return
// }

// BeValidRecord returns a non-blank string if json doesn't comply with
// JSONAPI rules for a resource record
func BeValidRecord(json *gabs.Container) (fail string) {
	fail = HaveFields(json, "id", reflect.String, "type", reflect.String)
	return
}

// BeValidRecordArray returns a non-blank string if json doesn't comply with
// JSONAPI rules for an array of resource records
func BeValidRecordArray(json *gabs.Container) (fail string) {
	children, err := json.Children()
	if err != nil {
		return err.Error()
	}
	for _, record := range children {
		if fail = BeValidRecord(record); fail != "" {
			return
		}
	}
	return
}

//
// func BeValidLinks(json *gabs.Container) (fail string) {
// 	return
// }
//
// func BeValidIncluded(json *gabs.Container) (fail string) {
// 	return
// }
