// Ported from dss-enumerations/.../CertificateType.java (DSS 6.5.RC1).
package enumerations

// CertificateType contains a list of possible certificate types.
type CertificateType string

const (
	// CertificateTypeESign is for electronic signature.
	CertificateTypeESign CertificateType = "ESIGN"
	// CertificateTypeESeal is for electronic seal.
	CertificateTypeESeal CertificateType = "ESEAL"
	// CertificateTypeWSA is for Web authentication.
	CertificateTypeWSA CertificateType = "WSA"
	// CertificateTypeUnknown is unknown.
	CertificateTypeUnknown CertificateType = "UNKNOWN"
)

var certificateTypeLabel = map[CertificateType]string{
	CertificateTypeESign:   "eSig",
	CertificateTypeESeal:   "eSeal",
	CertificateTypeWSA:     "WSA",
	CertificateTypeUnknown: "unknown",
}

// CertificateTypeValues returns all constants in declaration order.
func CertificateTypeValues() []CertificateType {
	return []CertificateType{
		CertificateTypeESign,
		CertificateTypeESeal,
		CertificateTypeWSA,
		CertificateTypeUnknown,
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
