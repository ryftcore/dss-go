// Ported from dss-json-common's eu.europa.esig.json.RFC3339DateUtils (DSS
// 6.5.RC1), collapsed into this package per doc.go's rationale - only the
// two methods CryptographicSuiteJsonCatalogue calls (getDate/getDateTime)
// are ported.
package cryptojson

import (
	"fmt"
	"time"
)

// rfc3339DateLayout is the Go reference-time layout for
// RFC3339DateUtils.DATE_PATTERN ("yyyy-MM-dd"), always interpreted in UTC
// (java.util.TimeZone.getTimeZone("UTC")).
const rfc3339DateLayout = "2006-01-02"

// rfc3339DateTimeLayouts are the Go reference-time layouts for
// RFC3339DateUtils.DATE_TIME_PATTERNS, tried in the same order, each in
// UTC.
var rfc3339DateTimeLayouts = []string{
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000Z07:00",
}

// rfc3339GetDate parses an IETF RFC 3339 date string. Ports
// RFC3339DateUtils#getDate; the Java IllegalArgumentException on an
// unparseable dateString becomes a returned error per PORTING.md, and a
// nil error/zero time.Time is returned for an empty dateString per Java's
// null-string short-circuit (see this package's callers, which only
// invoke this on a non-empty string obtained from asString - see
// json_object.go's doc comment on why "" is otherwise the not-present
// sentinel).
func rfc3339GetDate(dateString string) (time.Time, error) {
	if dateString == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(rfc3339DateLayout, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid 'date format': %s", dateString)
	}
	return t.UTC(), nil
}

// rfc3339GetDateTime parses an IETF RFC 3339 date-time string, trying each
// of rfc3339DateTimeLayouts in turn. Ports RFC3339DateUtils#getDateTime.
func rfc3339GetDateTime(dateTimeString string) (time.Time, error) {
	if dateTimeString == "" {
		return time.Time{}, nil
	}
	for _, layout := range rfc3339DateTimeLayouts {
		if t, err := time.Parse(layout, dateTimeString); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("unparseable 'date-time': %s", dateTimeString)
}
