// Parses the XML Schema xs:date/xs:dateTime lexical forms this package's
// documents use, standing in for javax.xml.datatype.XMLGregorianCalendar
// (which the JAXB RI binds xs:date/xs:dateTime to) - see
// CryptographicSuiteXmlCatalogue#toDate(XMLGregorianCalendar), collapsed
// into xsdDateToTime/xsdDateTimeToTime below since this port has no
// XMLGregorianCalendar equivalent to convert from.
package cryptoxml

import "time"

// xsdDateLayouts are the xs:date lexical forms this package's documents
// use: with a timezone offset (every dssc:End value in dss-crypto-suite.xml
// - see xml_types.go's ValidityType doc comment) and, defensively, without
// one (permitted by the XSD, even though no in-scope fixture exercises it).
var xsdDateLayouts = []string{
	"2006-01-02Z07:00",
	"2006-01-02",
}

// xsdDateTimeLayouts are the xs:dateTime lexical forms this package's
// documents use (PolicyIssueDate/NextUpdate in dss-crypto-suite.xml, e.g.
// "2024-10-13T00:00:00.000+00:00"), plus the "Z"-suffixed and
// fractional-second-less variants XML Schema also permits.
var xsdDateTimeLayouts = []string{
	"2006-01-02T15:04:05.000Z07:00",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05Z",
}

// xsdDateToTime parses an xs:date lexical value. Ports the date-only half
// of toDate(XMLGregorianCalendar); returns (zero, false) for an empty or
// unparseable value, mirroring toDate's null-in/null-out short-circuit
// (see json_object.go's doc comment on why "" is otherwise this port's
// not-present sentinel for optional string fields).
func xsdDateToTime(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range xsdDateLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

// xsdDateTimeToTime parses an xs:dateTime lexical value. Ports the
// dateTime half of toDate(XMLGregorianCalendar); see xsdDateToTime's doc
// comment for the empty/unparseable handling.
func xsdDateTimeToTime(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range xsdDateTimeLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}
