package should

import (
	"fmt"
	"strings"
)

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
func SplitMsg(msg string) (short, long, details, meta string) {
	chop := func(s, sep string) (found, rest string) {
		pos := strings.Index(s, sep)
		if pos < 0 {
			rest = s
		} else {
			found = s[:pos]
			rest = s[pos+len(sep):]
		}
		return
	}
	short, long = chop(msg, ShortSeparator)
	// fmt.Printf("msg %s\n short %s\n", msg, short)
	long, details = chop(long, LongSeparator)
	// fmt.Printf("long %s\n details %s\n", long, details)
	details, meta = chop(details, DetailsSeparator)
	// fmt.Printf("details %s\n meta %s\n", details, meta)
	if short == "" {
		return msg, "", "", ""
	} else if long == "" && details == "" {
		return short, meta, "", ""
	} else if details == "" {
		return short, long, meta, ""
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
