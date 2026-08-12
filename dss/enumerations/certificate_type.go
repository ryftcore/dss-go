// Ported from dss-enumerations/.../CertificateType.java (DSS 6.5.RC1).
package enumerations

// CertificateType contains a list of possible certificate types.
type CertificateType string

const (
	// CertificateType_ESIGN is for electronic signature.
	CertificateType_ESIGN CertificateType = "ESIGN"
	// CertificateType_ESEAL is for electronic seal.
	CertificateType_ESEAL CertificateType = "ESEAL"
	// CertificateType_WSA is for Web authentication.
	CertificateType_WSA CertificateType = "WSA"
	// CertificateType_UNKNOWN is unknown.
	CertificateType_UNKNOWN CertificateType = "UNKNOWN"
)

var certificateTypeLabel = map[CertificateType]string{
	CertificateType_ESIGN:   "eSig",
	CertificateType_ESEAL:   "eSeal",
	CertificateType_WSA:     "WSA",
	CertificateType_UNKNOWN: "unknown",
}

// CertificateTypeValues returns all constants in declaration order.
func CertificateTypeValues() []CertificateType {
	return []CertificateType{
		CertificateType_ESIGN,
		CertificateType_ESEAL,
		CertificateType_WSA,
		CertificateType_UNKNOWN,
	}
}

// Label returns the user-friendly label.
func (c CertificateType) Label() string {
	return certificateTypeLabel[c]
}

// CertificateTypeValueOf returns the constant matching the given Java enum
// name.
func CertificateTypeValueOf(name string) (CertificateType, error) {
	for _, v := range CertificateTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &certificateTypeInvalidValueError{name}
}

type certificateTypeInvalidValueError struct {
	name string
}

func (e *certificateTypeInvalidValueError) Error() string {
	return "no enum constant CertificateType." + e.name
}
