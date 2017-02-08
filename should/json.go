package should

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/Jeffail/gabs"
)

/* About the JSON parser: https://github.com/tidwall/gjson and
/* https://github.com/tidwall/gjson both fit most of our needs. Gjson is faster
/* most of the time, but uses unsafe and doesn't give distinct parse errors. */

func parseJSON(actual interface{}) (*gabs.Container, error) {
	switch v := actual.(type) {
	case string:
		container, err := gabs.ParseJSON([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("Error parsing JSON: %s\nBody: %s", err, v)
		}
		return container, err
	case *string:
		container, err := gabs.ParseJSON([]byte(*v))
		if err != nil {
			return nil, fmt.Errorf("Error parsing JSON: %s\nBody: %s", err, *v)
		}
		return container, err
	case []byte:
		return gabs.ParseJSON(v)
	case *gabs.Container:
		return v, nil
	default:
		return nil, fmt.Errorf("Expecting a JSON string or a structure representing one, not a %T.", actual)
	}
}

// HaveFields passes if the JSON container or string has fields with certain values or types of values:
//
//   HaveFields(json, "id", reflect.String)  // assert that there is a field `id` with a  string value.
//   HaveFields(json, "count", reflect.Float64)  // assert that there is a field `count` with a numeric value.
//   HaveFields(json, "count", 100)  // assert that there is a field `count` with a numeric value equal to 100.
//   HaveFields(json, "default", reflect.Interface)  // assert that there is a field `default` with any type of value.
//
func HaveFields(actual interface{}, expected ...interface{}) (fail string) {
	usage := "HaveFields expects parseable JSON to be compared to fieldPath string, fieldKind reflect.Kind pairs."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	return haveFields(json, true, expected...)
}

// AllowFields passes if fields in the JSON container or string either don't exist or match expected types.
//
//   AllowFields(json, "id", reflect.String)  // assert that there is a field `id` with a  string value.
//   AllowFields(json, "count", reflect.Float64)  // assert that there is a field `count` with an numeric value.
//   AllowFields(json, "default", reflect.Interface)  // assert that there is a field `default` with any type of value.
//
func AllowFields(actual interface{}, expected ...interface{}) (fail string) {
	usage := "HaveFields expects parseable JSON to be compared to fieldPath string, fieldKind reflect.Kind pairs."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	return haveFields(json, false, expected...)
}

// haveFields checks to see if json contains fields and types matching expected.
// if required=false, it tolerates fields not appearing in the object.
// expected is [fieldPath string, fieldKind reflect.Kind, ...]  pairs
func haveFields(json *gabs.Container, required bool, expected ...interface{}) (fail string) {
	for i := 0; i < len(expected); i += 2 {
		// check existence of key
		fieldPath := expected[i].(string)
		container := json.Path(fieldPath)
		if container == nil || container.Data() == nil { // field not found
			if required {
				fail += fmt.Sprintf("Field '%s' is missing.\n%s", fieldPath, json)
			}
			continue
		}

		expectedInterface := expected[i+1]
		expectedKind, ok := expectedInterface.(reflect.Kind)
		if !ok { // We have a value, not a Kind
			return Equal(container.Data(), expectedInterface)
		}
		// check type of value
		if expectedKind == reflect.Interface { // allow any type
			continue
		}
		actualKind := reflect.ValueOf(container.Data()).Kind()
		if actualKind != expectedKind {
			fail += fmt.Sprintf("Expecting a '%s' value of type %s, got %s.\nJSON: %s", fieldPath, expectedKind, actualKind, json)
		}
	}
	return
}

// HaveOnlyFields passes if the JSON container or string has fields with certain types of values:
//
//   HaveOnlyFields(json, "id", reflect.String)  // assert that there may a field `id` with a string value.
//   HaveOnlyFields(json, "count", reflect.Float64)  // assert that there may a field `count` with an numeric value.
//   HaveOnlyFields(json, "default", reflect.Interface)  // assert that there may a field `default` with any type of value.
//
func HaveOnlyFields(actual interface{}, allowed ...interface{}) (fail string) {
	usage := "HaveOnlyFields expects parseable JSON to be compared to an fieldPath string, fieldKind reflect.Kind pairs."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}

	fail += haveFields(json, false, allowed...)
	fail += haveOnlyKeys(json, allowed...)
	return
}

func haveOnlyKeys(json *gabs.Container, allowed ...interface{}) (fail string) {
	children, err := json.ChildrenMap()
	if err != nil {
		return err.Error()
	}
	for key, child := range children {
		found := false
		for _, v := range allowed {
			name, ok := v.(string)
			if !ok { // onlhy pay attention to strings
				continue
			}
			if key == name {
				found = true
				break
			}
		}
		if !found {
			fail += fmt.Sprintf("fields['%s'] (%s) is not allowed. ", key, child)
		}
	}
	return
}

// BeJSON asserts that the first argument can be parsed as JSON.
func BeJSON(actual interface{}, expected ...interface{}) (fail string) {
	usage := "BeJson expects a single string argument and passes if that argument parses as JSON."
	if actual == nil {
		return usage
	}
	_, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}
	return ""
}

// HaveOnlyCamelcaseKeys passes if all the attributes within a JSON container or
// string contain only upper and lower case ASCII letters.
//
func HaveOnlyCamelcaseKeys(actual interface{}, ignored ...interface{}) (fail string) {
	usage := "HaveOnlyCamelcaseKeys expects parseable JSON in actual. It's keys will recursively checked to make sure they have no snake_case keys. Any right-side arguments should be snake_case key names to be ignored."
	if actual == nil {
		return usage
	}
	json, err := parseJSON(actual)
	if err != nil {
		return err.Error()
	}

	ignoreMap := make(map[string]bool)
	for _, igI := range ignored {
		igS, ok := igI.(string)
		if !ok {
			return fmt.Sprintf("%s. One of the ignored values (%#v) is a %T, not a string.\n%#v", usage, igI, igI, ignored)
		}
		ignoreMap[igS] = true
	}

	return checkCamelcaseKeys(json, ignoreMap)
}

var camelCaseRegexp = regexp.MustCompile(`^[a-zA-Z]+$`)

func checkCamelcaseKeys(j *gabs.Container, ignores map[string]bool) (fail string) {
	// if j is an Object with keys, check each of keys and children
	children, notObjectErr := j.ChildrenMap()
	if notObjectErr == nil {
		for k, v := range children {
			if ignores[k] {
				return
			}
			if camelCaseRegexp.MatchString(k) {
				fail = checkCamelcaseKeys(v, ignores)
			} else {
				fail = fmt.Sprintf("Expecting only camelCase keys: found '%s'.\nWithin: %s ",
					k, j)
			}
			if fail != "" {
				return // fail fast to limit recursion
			}
		}
		return
	}

	// If j is an array, check each of it's elements
	// j.Children() will succeed for objects but lose the keys, so order is
	// it's important to check this AFTER the j.ChildMap() check
	elements, notArrayErr := j.Children()
	if notArrayErr == nil {
		for _, element := range elements {
			fail = checkCamelcaseKeys(element, ignores)
			if len(fail) > 0 {
				return // fail fast to limit recursion
			}
		}
		return
	}

	// otherwise this is an atomic object. No check necessary.
	return ""
}
