// Ported from dss-policy-jaxb/.../policy/DateUtils.java (DSS 6.5.RC1).
package policy

import "time"

// DateUtilsDefaultDateFormat is the default date format. Ports
// DateUtils#DEFAULT_DATE_FORMAT.
const DateUtilsDefaultDateFormat = "yyyy-MM-dd"

// javaToGoDateLayout translates the small subset of java.text.SimpleDateFormat
// pattern letters this package's callers use (see rule_utils.go and
// cryptographic_constraint_wrapper.go, whose only format strings are
// DateUtilsDefaultDateFormat itself and whatever AlgoExpirationDate.Format
// carries from a policy document - always "yyyy-MM-dd" in every upstream
// resource this chunk ports, see policy/jaxb/testdata) into the equivalent
// Go reference-time layout.
var javaToGoDateLayout = map[byte]string{
	'y': "2006",
	'M': "01",
	'd': "02",
	'H': "15",
	'm': "04",
	's': "05",
}

// DateUtilsParseDate converts the given string representation of the date
// using the format pattern. Ports DateUtils#parseDate(String, String).
//
// Java's SimpleDateFormat.parse (with setLenient(false)) throws a checked
// ParseException that DateUtils wraps and rethrows as an
// IllegalArgumentException; both become a returned error here per
// PORTING.md.
func DateUtilsParseDate(format, dateString string) (time.Time, error) {
	layout := javaSimpleDateFormatToGoLayout(format)
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return time.Time{}, &dateUtilsParseError{dateString: dateString, format: format, cause: err}
	}
	return t, nil
}

// javaSimpleDateFormatToGoLayout translates a java.text.SimpleDateFormat
// pattern into a Go reference-time layout, one run of identical pattern
// letters at a time - sufficient for the "yyyy-MM-dd"-shaped formats this
// package's callers use (see javaToGoDateLayout's doc comment). Any other
// character (notably '-' and '/') passes through unchanged, matching
// SimpleDateFormat's treatment of non-letter pattern characters as literal
// text.
func javaSimpleDateFormatToGoLayout(format string) string {
	var out []byte
	for i := 0; i < len(format); {
		c := format[i]
		layout, known := javaToGoDateLayout[c]
		if !known {
			out = append(out, c)
			i++
			continue
		}
		j := i
		for j < len(format) && format[j] == c {
			j++
		}
		out = append(out, layout...)
		i = j
	}
	return string(out)
}

// dateUtilsParseError ports the IllegalArgumentException DateUtils#parseDate
// throws, preserving the original message shape ("Unable to parse date %s
// (format:%s)") and wrapping the underlying error, mirroring the Java
// exception's cause chain.
type dateUtilsParseError struct {
	dateString string
	format     string
	cause      error
}

func (e *dateUtilsParseError) Error() string {
	return "Unable to parse date " + e.dateString + " (format:" + e.format + ")"
}

func (e *dateUtilsParseError) Unwrap() error {
	return e.cause
}
