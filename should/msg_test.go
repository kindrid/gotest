package should

import "testing"

const (
	shortMsg       = "Brief message"
	longMsg        = "Extra message with a long line.\n"
	detailsMsg     = "Here are some \ndetails"
	metaMsg        = "Hereare a bunch of technical details about test workings \nfor when you doubt the assertion or the runner."
	failShort      = shortMsg
	failExtra      = ShortSeparator + longMsg
	failDetails    = LongSeparator + detailsMsg
	failMeta       = DetailsSeparator + metaMsg
	failShortExtra = failShort + failExtra
	failAll        = failShort + failExtra + failDetails + failMeta

	failShortMeta = failShort + failMeta
)

func testStringEqual(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Errorf("Expected %#v\nTo Equal %#v", actual, expected)
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
	testMessageParse(t, "", "", "", "", "")
	testMessageParse(t, failShort, shortMsg, "", "", "")
	// testMessageParse(t, failAll, shortMsg, longMsg, detailsMsg, metaMsg)

}
