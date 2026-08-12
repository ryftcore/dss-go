package spi

import (
	"testing"
	"time"
)

// TestDSSUtilsFormatDateWithCustomFormatKAT pins the java.text.SimpleDateFormat pattern
// translation, in particular the millisecond field 'S', which has no direct Go layout
// counterpart. The expected values come from upstream DSSUtils#formatDateWithCustomFormat on
// OpenJDK 21 with the UTC time zone.
func TestDSSUtilsFormatDateWithCustomFormatKAT(t *testing.T) {
	testCases := []struct {
		millis   int64
		pattern  string
		expected string
	}{
		{0, "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", "1970-01-01T00:00:00.000Z"},
		{1786565956000, "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", "2026-08-12T20:19:16.000Z"},
		{1786565956123, "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", "2026-08-12T20:19:16.123Z"},
		{-1000, "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", "1969-12-31T23:59:59.000Z"},
		{253402300799000, "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", "9999-12-31T23:59:59.000Z"},
		{1786565956123, "yyyy-MM-dd'T'HH:mm:ss'Z'", "2026-08-12T20:19:16Z"},
		{1786565956123, "yyyy-MM-dd", "2026-08-12"},
	}
	for _, testCase := range testCases {
		date := time.UnixMilli(testCase.millis).UTC()
		if got := DSSUtilsFormatDateWithCustomFormat(date, testCase.pattern); got != testCase.expected {
			t.Errorf("DSSUtilsFormatDateWithCustomFormat(%d, %q) = %q, want %q",
				testCase.millis, testCase.pattern, got, testCase.expected)
		}
	}
}

// TestDSSUtilsParseDateKAT pins the two SimpleDateFormat behaviours the Go time package does
// not share: parse() stops as soon as the pattern is satisfied and ignores the rest of the
// string, while it never accepts a fractional second the pattern does not declare. The
// expected values come from upstream DSSUtils#parseISO8601Date / #parseRFCDate on OpenJDK 21.
func TestDSSUtilsParseDateKAT(t *testing.T) {
	const iso8601Midnight = 1786492800000 // 2026-08-12T00:00:00Z
	const rfcInstant = 1786565956000      // 2026-08-12T20:19:16Z

	testCases := []struct {
		value       string
		iso8601     int64 // 0 means "unparseable"
		rfc         int64
		isISO8601   bool
		isRFCFormat bool
	}{
		{"2026-08-12", iso8601Midnight, 0, true, false},
		// "yyyy-MM-dd" is satisfied by the first ten characters; the time is dropped.
		{"2026-08-12T20:19:16Z", iso8601Midnight, rfcInstant, true, true},
		{"2026-08-12T20:19:16+02:00", iso8601Midnight, 0, true, false},
		// A fractional second the RFC 3339 pattern does not declare is rejected.
		{"2026-08-12T20:19:16.123Z", iso8601Midnight, 0, true, false},
		{"not a date", 0, 0, false, false},
		{"", 0, 0, false, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.value, func(t *testing.T) {
			parsedISO := DSSUtilsParseISO8601Date(testCase.value)
			if testCase.iso8601 == 0 {
				if !parsedISO.IsZero() {
					t.Errorf("DSSUtilsParseISO8601Date(%q) = %v, want no date", testCase.value, parsedISO)
				}
			} else if got := parsedISO.UnixMilli(); got != testCase.iso8601 {
				t.Errorf("DSSUtilsParseISO8601Date(%q) = %d, want %d", testCase.value, got, testCase.iso8601)
			}

			parsedRFC := DSSUtilsParseRFCDate(testCase.value)
			if testCase.rfc == 0 {
				if !parsedRFC.IsZero() {
					t.Errorf("DSSUtilsParseRFCDate(%q) = %v, want no date", testCase.value, parsedRFC)
				}
			} else if got := parsedRFC.UnixMilli(); got != testCase.rfc {
				t.Errorf("DSSUtilsParseRFCDate(%q) = %d, want %d", testCase.value, got, testCase.rfc)
			}

			if got := DSSUtilsIsISO8601Date(testCase.value); got != testCase.isISO8601 {
				t.Errorf("DSSUtilsIsISO8601Date(%q) = %t, want %t", testCase.value, got, testCase.isISO8601)
			}
			if got := DSSUtilsIsRFCDate(testCase.value); got != testCase.isRFCFormat {
				t.Errorf("DSSUtilsIsRFCDate(%q) = %t, want %t", testCase.value, got, testCase.isRFCFormat)
			}
		})
	}
}
