// Ported from dss-enumerations/.../EAAQualification.java (DSS 6.5.RC1).
package enumerations

// EAAQualification defines possible supported qualification types for an
// Electronic Attestation of Attributes.
type EAAQualification string

const (
	// EAAQualificationQEAA is a Qualified electronic attestation of
	// attributes as defined in Regulation EU 2024/1183, Article 45d.
	EAAQualificationQEAA EAAQualification = "QEAA"
	// EAAQualificationEAA is an Electronic attestation of attributes as
	// defined in Regulation EU 2024/1183, without a qualified status.
	EAAQualificationEAA EAAQualification = "EAA"
	// EAAQualificationPubEAA is an Electronic attestation of attributes
	// issued by or on behalf of a public sector body responsible for an
	// authentic source as defined in Regulation EU 2024/1183, Article 45f.
	EAAQualificationPubEAA EAAQualification = "PUBEAA"
	// EAAQualificationPID is Personal Identification Data (PID).
	EAAQualificationPID EAAQualification = "PID"
	// EAAQualificationUnknown is an Electronic attestation of attributes
	// of unknown or conflicting status.
	EAAQualificationUnknown EAAQualification = "UNKNOWN"
	// EAAQualificationIndeterminateQEAA is an Indeterminate qualified
	// electronic attestation of attributes as defined in Regulation EU
	// 2024/1183, Article 45d.
	EAAQualificationIndeterminateQEAA EAAQualification = "INDETERMINATE_QEAA"
	// EAAQualificationIndeterminateEAA is an Indeterminate electronic
	// attestation of attributes as defined in Regulation EU 2024/1183,
	// without a qualified status.
	EAAQualificationIndeterminateEAA EAAQualification = "INDETERMINATE_EAA"
	// EAAQualificationIndeterminatePubEAA is an Indeterminate electronic
	// attestation of attributes issued by or on behalf of a public sector
	// body responsible for an authentic source as defined in Regulation EU
	// 2024/1183, Article 45f.
	EAAQualificationIndeterminatePubEAA EAAQualification = "INDETERMINATE_PUBEAA"
	// EAAQualificationIndeterminatePID is Indeterminate Personal
	// Identification Data (PID).
	EAAQualificationIndeterminatePID EAAQualification = "INDETERMINATE_PID"
	// EAAQualificationIndeterminateUnknown is an Indeterminate electronic
	// attestation of attributes of unknown or conflicting status.
	EAAQualificationIndeterminateUnknown EAAQualification = "INDETERMINATE_UNKNOWN"
	// EAAQualificationNotEAA is Not electronic attestation of attributes.
	EAAQualificationNotEAA EAAQualification = "NOT_EAA"
	// EAAQualificationNA is Not Applicable.
	EAAQualificationNA EAAQualification = "NA"
)

type eaaQualificationFields struct {
	readable string
	label    string
	uri      string
}

// eaaQualificationData holds the (readable, label, uri) tuple for each
// constant.
var eaaQualificationData = map[EAAQualification]eaaQualificationFields{
	EAAQualificationQEAA:                 {"QEAA", "Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:QEAA"},
	EAAQualificationEAA:                  {"EAA", "Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:EAA"},
	EAAQualificationPubEAA:               {"PuB-EAA", "Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:PUBEAA"},
	EAAQualificationPID:                  {"PID", "Personal Identification Data", "urn:cef:dss:eaaQualification:PID"},
	EAAQualificationUnknown:              {"Unknown", "Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:Unknown"},
	EAAQualificationIndeterminateQEAA:    {"Indeterminate QEAA", "Indeterminate Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateQEAA"},
	EAAQualificationIndeterminateEAA:     {"Indeterminate EAA", "Indeterminate Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateEAA"},
	EAAQualificationIndeterminatePubEAA:  {"Indeterminate Pub-EAA", "Indeterminate Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:indeterminatePUBEAA"},
	EAAQualificationIndeterminatePID:     {"Indeterminate PID", "Indeterminate Personal Identification Data", "urn:cef:dss:eaaQualification:indeterminatePID"},
	EAAQualificationIndeterminateUnknown: {"Indeterminate Unknown", "Indeterminate Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:indeterminateUnknown"},
	EAAQualificationNotEAA:               {"Not EAA", "Not Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:NOTEAA"},
	EAAQualificationNA:                   {"N/A", "Not applicable", "urn:cef:dss:eaaQualification:NA"},
}

// EAAQualificationValues returns all constants in declaration order.
func EAAQualificationValues() []EAAQualification {
	return []EAAQualification{
		EAAQualificationQEAA,
		EAAQualificationEAA,
		EAAQualificationPubEAA,
		EAAQualificationPID,
		EAAQualificationUnknown,
		EAAQualificationIndeterminateQEAA,
		EAAQualificationIndeterminateEAA,
		EAAQualificationIndeterminatePubEAA,
		EAAQualificationIndeterminatePID,
		EAAQualificationIndeterminateUnknown,
		EAAQualificationNotEAA,
		EAAQualificationNA,
	}
}

// eaaQualificationByReadable is a lookup map from the readable short-string
// to the EAAQualification constant.
var eaaQualificationByReadable = func() map[string]EAAQualification {
	m := make(map[string]EAAQualification, len(eaaQualificationData))
	for _, v := range EAAQualificationValues() {
		m[eaaQualificationData[v].readable] = v
	}
	return m
}()

// Readable gets user-friendly name of the enumeration.
func (e EAAQualification) Readable() string {
	return eaaQualificationData[e].readable
}

// Label gets description of the enumeration.
func (e EAAQualification) Label() string {
	return eaaQualificationData[e].label
}

// URI gets a unique URI.
func (e EAAQualification) URI() string {
	return eaaQualificationData[e].uri
}

// EAAQualificationForName gets EAAQualification from an enumeration name.
// Note: the result can be "" (zero value): empty input returns "" with no
// error, mirroring Java's null return; an unknown name is an error,
// mirroring Java's valueOf IllegalArgumentException.
func EAAQualificationForName(value string) (EAAQualification, error) {
	if value == "" {
		return "", nil
	}
	return EAAQualificationValueOf(value)
}

// EAAQualificationValueOf returns the constant matching the given Java enum
// name.
func EAAQualificationValueOf(name string) (EAAQualification, error) {
	for _, v := range EAAQualificationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &eaaQualificationInvalidValueError{name}
}

type eaaQualificationInvalidValueError struct {
	name string
}

func (e *eaaQualificationInvalidValueError) Error() string {
	return "no enum constant EAAQualification." + e.name
}

// EAAQualificationFromReadable gets EAAQualification from a readable
// user-friendly label. Note: EAAQualification can be "" (zero value) if
// readable is empty or unknown, mirroring Java's null return.
func EAAQualificationFromReadable(readable string) EAAQualification {
	if readable == "" {
		return ""
	}
	return eaaQualificationByReadable[readable]
}

// compile-time interface assertion.
var _ UriBasedEnum = EAAQualification("")
