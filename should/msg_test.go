package should

import "testing"

const (
	shortMsg      = "Brief message"
	longMsg       = "Extra message with a long line.\n"
	detailsMsg    = "Here are some \ndetails"
	metaMsg       = "Here are a bunch of technical details about test workings \nfor when you doubt the assertion or the runner."
	failShort     = shortMsg
	failLong      = ShortSeparator + longMsg
	failDetails   = LongSeparator + detailsMsg
	failMeta      = DetailsSeparator + metaMsg
	failShortLong = failShort + failLong
	failAll       = failShort + failLong + failDetails + failMeta

	failShortMeta = failShort + failMeta
)

func testStringEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Actual %#v\nExpected %#v", actual, expected)
	}
}

func testMessageParse(t *testing.T, fullMsg, short, long, details, meta string) {
	s, l, d, m := SplitMsg(fullMsg)
	testStringEqual(t, s, short)
	testStringEqual(t, l, long)
	testStringEqual(t, d, details)
	testStringEqual(t, m, meta)
}

func TestFailureMessageCreation(t *testing.T) {
	// Test everything and nothing
	testStringEqual(t, JoinMsg(shortMsg, longMsg, detailsMsg, metaMsg), failAll)
	testStringEqual(t, JoinMsg("", "", "", ""), "")

	// Test combinations
	testStringEqual(t, JoinMsg(shortMsg, "", "", ""), failShort)
	testStringEqual(t, JoinMsg(shortMsg, "", detailsMsg, ""), failShort+failDetails)
	testStringEqual(t, JoinMsg(shortMsg, "", "", metaMsg), failShortMeta)
}

func TestFailureMessageParsing(t *testing.T) {
	// test normal order
	testMessageParse(t, "", "", "", "", "")
	testMessageParse(t, failShort, shortMsg, "", "", "")
	testMessageParse(t, failShort+failLong, shortMsg, longMsg, "", "")
	testMessageParse(t, failAll, shortMsg, longMsg, detailsMsg, metaMsg)

	// short grabs everything if the msg seems non-compliant
	testMessageParse(t, failShort, failShort, "", "", "")
	testMessageParse(t, failLong, failLong, "", "", "")
	testMessageParse(t, failDetails, failDetails, "", "", "")
	testMessageParse(t, failMeta, failMeta, "", "", "")

	// test some malformed messages (missing sections)
	testMessageParse(t, failShort+failMeta, shortMsg, "long", "", metaMsg)

}
