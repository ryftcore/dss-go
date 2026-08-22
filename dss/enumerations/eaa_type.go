// Ported from dss-enumerations/.../EAAType.java (DSS 6.5.RC1).
package enumerations

// EAAType defines a list of EAA types known by the current implementation.
// NOTE: This type relates to a format of an EAA signed token.
type EAAType string

const (
	// EAATypeSDJWTVC is a realization of EAA that implements EAA as a
	// JSON Web Signature as specified in IETF RFC 7515, built on IETF
	// SD-JWT VC, which further profiles a Selective Disclosure JSON Web
	// Token as specified in IETF RFC 9901 "Selective Disclosure for JSON
	// Web Tokens".
	EAATypeSDJWTVC EAAType = "SD_JWT_VC"
	// EAATypeISOIECMDoc is a realization of EAA that implements EAA
	// built on the data structures defined in ISO/IEC 18013-5.
	EAATypeISOIECMDoc EAAType = "ISO_IEC_MDOC"
	// EAATypeW3CVC is a realization of EAA based on the JSON-LD
	// (specified in W3C Recommendation: "JSON-LD 1.1. A JSON-based
	// Serialization for Linked Data") serialization of W3C Recommendation
	// (15 May 2025): "Verifiable Credentials Data Model v2.0".
	EAATypeW3CVC EAAType = "W3C_VC"
	// EAATypeX509AC is a realization of EAA based on X.509 Attribute
	// certificates as specified in IETF RFC 5755.
	EAATypeX509AC EAAType = "X509_AC"
)

// EAATypeValues returns all constants in declaration order.
func EAATypeValues() []EAAType {
	return []EAAType{
		EAATypeSDJWTVC,
		EAATypeISOIECMDoc,
		EAATypeW3CVC,
		EAATypeX509AC,
	}
}
