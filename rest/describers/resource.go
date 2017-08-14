package describers

import (
	"reflect"
	"regexp"
)

// Resources represents a group of registered resources
type Resources struct {
	All    []*Resource
	allMap map[string]*Resource
}

// Resource describes a type of value in a response.
type Resource struct {
	Type        string // string that tags this type of resource
	Description string
	Attributes  []*Attribute // The native attributes of the core resource
	Includes    map[string]string
	Filters     map[string]string
	IDs         map[string]string // example ID's. IDs["default"] should always be set.
	Special     bool              // if true, it's not a normal resource
	Finalized   bool
}

// Permission is a bit field of the field-level actions a role is permitted.
type Permission int

// Base permission constants
const (
	Read   Permission = 1 << iota // Read means the field can be read
	Update                        // Update means the field can be written to in an existing record
	Create                        // Create means the field can be written to in a new record
)

// Combined permission constants
const (
	RU  Permission = Read | Update          // RU permits Read and Update
	RC             = Read | Create          // RC permits Read and Create
	UC             = Update | Create        // UC permits Update and Create
	RUC            = Read | Update | Create // RUC permits Read, Update, and Create
)

// reflection constants used to represent MSON data types
const (
	BoolField        reflect.Kind = reflect.Bool    // BoolField  JSON boolean value
	StringField                   = reflect.String  // JSON string value
	NumberField                   = reflect.Float64 // JSON Number primitive value
	ArrayField                    = reflect.Slice   // JSON Array
	ObjectField                   = reflect.Map     // JSON Object
	CustomStructType              = reflect.Struct  // MSON defined type...won't we use relationships for this?
)

// Attribute describe one field in a response.
// Consider Default and example values might be best handled by an instance of the struct itself.
// or more type-safe (allowing us to excise reflect), a factory function
type Attribute struct {
	Name              string            // name of the field
	Description       string            // description of the field
	Kind              reflect.Kind      // type of field
	ExampleString     string            // string version of example value
	Required          bool              // true if the record must have this attribute filled
	Nullable          bool              // true if a required value can be nil (if required is false, nullable is forced true).
	Fixed             bool              // true if the value is a fixed value (in which case .Example is that value)
	Enums             map[string]string // string form of legal values if an enum field and their descriptions
	Regexp            *regexp.Regexp    // embedded validator regexp
	PermissionsInDocs Permission        // we'll expand these later
	Deprecation       string            // if not empty the field is being deprecated with this explanation/info/timeline
	Finalized         bool              // you shouldn't use unfinalized attributes. Call resources.Register() do it.
}
