package should

import "strings"

// tools for working with failure messages

const (
	// ShortSeparator ends the failure message brief portion (and begins the explanation)
	ShortSeparator = "\n"
	// LongSeparator ends the failure message explanation section (and begins the details)
	LongSeparator = "\n# DETAILS:"
	// DetailsSeparator ends the failure debugging portion (and begins details for debugging asserts and test runner.)
	DetailsSeparator = "\n# INTERNALS:"
)

// SplitMsg divides a failure message into parts that may be muted depending on verbosity levels
func SplitMsg(fail string) (short, long, details, meta string) {
	findOrSetToLength := func(s string) int {
		r := strings.Index(fail, s)
		if r < 0 {
			r = len(fail)
		}
		return r
	}
	shortZ := findOrSetToLength(ShortSeparator)
	// if shortZ == -1 { // this fail string isn't playing by our rules
	// 	shortZ = failZ
	// }
	// longZ = strings.Index(fail, LongSeparator)
	// detailsZ = strings.Index(fail, DetailsSeparator)
	//

	short = fail[:shortZ]
	// long = fail[longZ:]
	// details = fail[detailsZ:]
	// meta = fail[metaZ:]
	return
}

// JoinMsg creates a failure message from its components
func JoinMsg(short, long, details, meta string) (result string) {
	result = short
	if long != "" {
		result += ShortSeparator + long
	}
	if details != "" {
		result += LongSeparator + details
	}
	if meta != "" {
		result += DetailsSeparator + meta
	}
	return
}
