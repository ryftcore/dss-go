// Ported from dss-enumerations/.../TimestampQualification.java (DSS 6.5.RC1).
package enumerations

// TimestampQualification defines possible timestamp qualification types.
type TimestampQualification string

const (
	// TimestampQualification_QTSA is a Qualified timestamp token.
	TimestampQualification_QTSA TimestampQualification = "QTSA"
	// TimestampQualification_TSA is a Not-qualified timestamp token.
	TimestampQualification_TSA TimestampQualification = "TSA"
	// TimestampQualification_NA is Not-applicable (not determined).
	TimestampQualification_NA TimestampQualification = "NA"
)

type timestampQualificationFields struct {
	readable string
	label    string
	uri      string
}

// timestampQualificationData holds the (readable, label, uri) tuple for
// each constant.
var timestampQualificationData = map[TimestampQualification]timestampQualificationFields{
	TimestampQualification_QTSA: {"QTSA", "Qualified timestamp", "urn:cef:dss:timestampQualification:QTSA"},
	TimestampQualification_TSA:  {"TSA", "Not qualified timestamp", "urn:cef:dss:timestampQualification:TSA"},
	TimestampQualification_NA:   {"N/A", "Not applicable", "urn:cef:dss:timestampQualification:notApplicable"},
}

// TimestampQualificationValues returns all constants in declaration order.
func TimestampQualificationValues() []TimestampQualification {
	return []TimestampQualification{
		TimestampQualification_QTSA,
		TimestampQualification_TSA,
		TimestampQualification_NA,
	}
}

// Readable returns a short name of the qualification status.
func (t TimestampQualification) Readable() string {
	return timestampQualificationData[t].readable
}

// Label returns a complete name of the qualification status.
func (t TimestampQualification) Label() string {
	return timestampQualificationData[t].label
}

// URI returns a URI of the qualification status.
func (t TimestampQualification) URI() string {
	return timestampQualificationData[t].uri
}
