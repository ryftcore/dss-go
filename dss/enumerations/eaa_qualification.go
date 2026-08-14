// Ported from dss-enumerations/.../EAAQualification.java (DSS 6.5.RC1).
package enumerations

// EAAQualification defines possible supported qualification types for an
// Electronic Attestation of Attributes.
type EAAQualification string

const (
	// EAAQualification_QEAA is a Qualified electronic attestation of
	// attributes as defined in Regulation EU 2024/1183, Article 45d.
	EAAQualification_QEAA EAAQualification = "QEAA"
	// EAAQualification_EAA is an Electronic attestation of attributes as
	// defined in Regulation EU 2024/1183, without a qualified status.
	EAAQualification_EAA EAAQualification = "EAA"
	// EAAQualification_PUBEAA is an Electronic attestation of attributes
	// issued by or on behalf of a public sector body responsible for an
	// authentic source as defined in Regulation EU 2024/1183, Article 45f.
	EAAQualification_PUBEAA EAAQualification = "PUBEAA"
	// EAAQualification_PID is Personal Identification Data (PID).
	EAAQualification_PID EAAQualification = "PID"
	// EAAQualification_UNKNOWN is an Electronic attestation of attributes
	// of unknown or conflicting status.
	EAAQualification_UNKNOWN EAAQualification = "UNKNOWN"
	// EAAQualification_INDETERMINATE_QEAA is an Indeterminate qualified
	// electronic attestation of attributes as defined in Regulation EU
	// 2024/1183, Article 45d.
	EAAQualification_INDETERMINATE_QEAA EAAQualification = "INDETERMINATE_QEAA"
	// EAAQualification_INDETERMINATE_EAA is an Indeterminate electronic
	// attestation of attributes as defined in Regulation EU 2024/1183,
	// without a qualified status.
	EAAQualification_INDETERMINATE_EAA EAAQualification = "INDETERMINATE_EAA"
	// EAAQualification_INDETERMINATE_PUBEAA is an Indeterminate electronic
	// attestation of attributes issued by or on behalf of a public sector
	// body responsible for an authentic source as defined in Regulation EU
	// 2024/1183, Article 45f.
	EAAQualification_INDETERMINATE_PUBEAA EAAQualification = "INDETERMINATE_PUBEAA"
	// EAAQualification_INDETERMINATE_PID is Indeterminate Personal
	// Identification Data (PID).
	EAAQualification_INDETERMINATE_PID EAAQualification = "INDETERMINATE_PID"
	// EAAQualification_INDETERMINATE_UNKNOWN is an Indeterminate electronic
	// attestation of attributes of unknown or conflicting status.
	EAAQualification_INDETERMINATE_UNKNOWN EAAQualification = "INDETERMINATE_UNKNOWN"
	// EAAQualification_NOT_EAA is Not electronic attestation of attributes.
	EAAQualification_NOT_EAA EAAQualification = "NOT_EAA"
	// EAAQualification_NA is Not Applicable.
	EAAQualification_NA EAAQualification = "NA"
)

type eaaQualificationFields struct {
	readable string
	label    string
	uri      string
}

// eaaQualificationData holds the (readable, label, uri) tuple for each
// constant.
var eaaQualificationData = map[EAAQualification]eaaQualificationFields{
	EAAQualification_QEAA:                  {"QEAA", "Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:QEAA"},
	EAAQualification_EAA:                   {"EAA", "Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:EAA"},
	EAAQualification_PUBEAA:                {"PuB-EAA", "Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:PUBEAA"},
	EAAQualification_PID:                   {"PID", "Personal Identification Data", "urn:cef:dss:eaaQualification:PID"},
	EAAQualification_UNKNOWN:               {"Unknown", "Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:Unknown"},
	EAAQualification_INDETERMINATE_QEAA:    {"Indeterminate QEAA", "Indeterminate Qualified Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateQEAA"},
	EAAQualification_INDETERMINATE_EAA:     {"Indeterminate EAA", "Indeterminate Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:indeterminateEAA"},
	EAAQualification_INDETERMINATE_PUBEAA:  {"Indeterminate Pub-EAA", "Indeterminate Electronic Attestation of Attributes issued by or on behalf of a public sector body", "urn:cef:dss:eaaQualification:indeterminatePUBEAA"},
	EAAQualification_INDETERMINATE_PID:     {"Indeterminate PID", "Indeterminate Personal Identification Data", "urn:cef:dss:eaaQualification:indeterminatePID"},
	EAAQualification_INDETERMINATE_UNKNOWN: {"Indeterminate Unknown", "Indeterminate Electronic Attestation of Attributes of unknown type", "urn:cef:dss:eaaQualification:indeterminateUnknown"},
	EAAQualification_NOT_EAA:               {"Not EAA", "Not Electronic Attestation of Attributes", "urn:cef:dss:eaaQualification:NOTEAA"},
	EAAQualification_NA:                    {"N/A", "Not applicable", "urn:cef:dss:eaaQualification:NA"},
}

// EAAQualificationValues returns all constants in declaration order.
func EAAQualificationValues() []EAAQualification {
	return []EAAQualification{
		EAAQualification_QEAA,
		EAAQualification_EAA,
		EAAQualification_PUBEAA,
		EAAQualification_PID,
		EAAQualification_UNKNOWN,
		EAAQualification_INDETERMINATE_QEAA,
		EAAQualification_INDETERMINATE_EAA,
		EAAQualification_INDETERMINATE_PUBEAA,
		EAAQualification_INDETERMINATE_PID,
		EAAQualification_INDETERMINATE_UNKNOWN,
		EAAQualification_NOT_EAA,
		EAAQualification_NA,
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
