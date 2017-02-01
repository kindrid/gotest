package cuke

import (
	"github.com/DATA-DOG/godog"
)

// Cuke encapsulates cucumber test run settings
type Cuke struct {
	FeaturesDir  string       // directory or file to search for features (see GoDog docs https://github.com/DATA-DOG/godog)
	decideFilter string       // which features to run, a Regexp?, comma sep?
	Format       OutputFormat // output format: junit, pretty, progress
	FailFast     bool         // if true, exit on first failure

}

// OutputFormat are output choices for cucumber tests.
type OutputFormat int

const (
	text OutputFormat = iota
	pretty
	terse
	markdown
	html
	junit
)

// Run runs Cucumber-style feature tests using GoDog for step definitions. IWBNI
// we could factor all the Godog dependent stuff into this file.
//
// steps - a typical GoDog FeatureContext() function that sets up step
// definitions, etc.
func Run(name string, featuresLocations []string, steps func(s *godog.Suite)) (failures int) {
	opts := godog.Options{
		Format: "progress",
		Paths:  featuresLocations,
	}
	return godog.RunWithOptions(name, steps, opts)
}
