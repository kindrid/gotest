package should

import (
	"fmt"
	"strings"
)

// tools for working with failure messages

const (
	// ShortSeparator ends the failure message short portion (and begins the long
	// portion)
	ShortSeparator = "\n"

	// ShortLength is an arbitrary length used to shorten failure messages that
	// don't contain ShortSeparator
	ShortLength = 80

	// LongSeparator ends the failure message explanation section (and begins the
	// details)
	LongSeparator = "\n# DETAILS:"

	// DetailsSeparator ends the failure debugging portion (and begins details for
	// debugging asserts and test runner.)
	DetailsSeparator = "\n# INTERNALS:"

	// SectionSeparator separates the long, details, and internals sections.
	SectionSeparator = "\n~~~~~~~~~~\n"
)

func trim(s string) string {
	return strings.Trim(s, " \n\t\r")
}

func splitShortLong(s string) (short, long string) {
	sl := strings.SplitN(trim(s), ShortSeparator, 2)
	if len(sl) > 1 {
		return trim(sl[0]), trim(sl[1])
	}
	if len(s) > ShortLength {
		return trim(s[:ShortLength]), trim(s[ShortLength:])
	}
	return trim(s), ""
}

// SplitMsg divides a failure message into parts that may be muted depending on verbosity levels
func SplitMsg(msg string) (short, long, details, meta string) {
	if msg == "" {
		return
	}
	secs := strings.Split(msg, SectionSeparator)
	short, long = splitShortLong(secs[0])
	if len(secs) > 1 {
		details = trim(secs[1])
	}
	if len(secs) > 2 {
		meta = trim(secs[2])
	}
	return
}

func oldSplitMsg(msg string) (short, long, details, meta string) {
	msgZ := len(msg)
	findAndClean := func(start int, find string) (string, int) {
		fmt.Printf("Searching from %d for %s.\n", start, find)
		if start >= msgZ {
			return "", msgZ
		}
		pos := strings.Index(msg[start:], find)
		if pos < 0 { // not found
			return "", start
		}
		return msg[start:pos], pos + len(find)
	}
	// shortZ := findOrSetToLength(ShortSeparator)
	// longZ := findOrSetToLength(LongSeparator)
	pos := 0
	short, pos = findAndClean(pos, ShortSeparator)
	if short == "" { // special case: not formatted by our rules or only a short message
		short = msg
		return
	}
	long, _ = findAndClean(pos, LongSeparator)
	// details = msg[detailsZ:]
	// meta = fail[metaZ:]
	return
}

// JoinMsg creates a failure message from its components
func JoinMsg(short, long, details, meta string) (result string) {
	result = short + ShortSeparator + long + SectionSeparator + details + SectionSeparator + meta
	return
}
